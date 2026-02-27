package outbound

import (
	"bytes"
	"context"
	gotls "crypto/tls"
	"encoding/base64"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/dispatcher"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/features/stats"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	utls "github.com/refraction-networking/utls"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/buf"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/errors"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/protocol"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/retry"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/session"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/signal"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/task"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/xudp"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/core"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/features/policy"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/vless"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/vless/encoding"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/vless/encryption"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport/internet"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport/internet/reality"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport/internet/stat"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport/internet/tls"
)

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}

// Handler is an outbound connection handler for VLess protocol.
type Handler struct {
	serverList    *protocol.ServerList
	serverPicker  protocol.ServerPicker
	policyManager policy.Manager
	statManager   stats.Manager
	cone          bool
	encryption    *encryption.ClientInstance
}

// New creates a new VLess outbound handler.
func New(ctx context.Context, config *Config) (*Handler, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Vnext {
		s, err := protocol.NewServerSpecFromPB(rec)
		if err != nil {
			return nil, errors.New("failed to parse server spec").Base(err).AtError()
		}
		serverList.AddServer(s)
	}

	v := core.MustFromContext(ctx)
	handler := &Handler{
		serverList:    serverList,
		serverPicker:  protocol.NewRoundRobinServerPicker(serverList),
		policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
		statManager:   v.GetFeature(stats.ManagerType()).(stats.Manager),
		cone:          ctx.Value("cone").(bool),
	}

	a := handler.serverPicker.PickServer().PickUser().Account.(*vless.MemoryAccount)
	if a.Encryption != "" && a.Encryption != "none" {
		s := strings.Split(a.Encryption, ".")
		var nfsPKeysBytes [][]byte
		for _, r := range s {
			b, _ := base64.RawURLEncoding.DecodeString(r)
			nfsPKeysBytes = append(nfsPKeysBytes, b)
		}
		handler.encryption = &encryption.ClientInstance{}
		if err := handler.encryption.Init(nfsPKeysBytes, a.XorMode, a.Seconds, a.Padding); err != nil {
			return nil, errors.New("failed to use encryption").Base(err).AtError()
		}
	}

	return handler, nil
}

// Process implements proxy.Outbound.Process().
func (h *Handler) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	outbounds := session.OutboundsFromContext(ctx)
	ob := outbounds[len(outbounds)-1]
	if !ob.Target.IsValid() {
		return errors.New("target not specified").AtError()
	}
	ob.Name = "vless"

	var rec *protocol.ServerSpec
	var conn stat.Connection
	if err := retry.ExponentialBackoff(5, 200).On(func() error {
		rec = h.serverPicker.PickServer()
		var err error
		conn, err = dialer.Dial(ctx, rec.Destination())
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return errors.New("failed to find an available destination").Base(err).AtWarning()
	}
	defer conn.Close()

	iConn := conn
	if statConn, ok := iConn.(*stat.CounterConnection); ok {
		iConn = statConn.Connection
	}
	target := ob.Target

	clientId := session.ClientIdFromContext(ctx)
	errors.LogInfo(ctx, "tunneling request to ", target, " via ", rec.Destination().NetAddr(), " clientId:", clientId)
	if h.encryption != nil {
		var err error
		if conn, err = h.encryption.Handshake(conn); err != nil {
			return errors.New("ML-KEM-768 handshake failed").Base(err).AtInfo()
		}
	}

	command := protocol.RequestCommandTCP
	if target.Network == net.Network_UDP {
		command = protocol.RequestCommandUDP
	}
	if target.Address.Family().IsDomain() && target.Address.Domain() == "v1.mux.cool" {
		command = protocol.RequestCommandMux
	}

	//todo: add vless traffic counter
	if sizeStatWriter, ok := link.Writer.(*dispatcher.SizeStatWriter); ok {
		name := "user>>>" + strconv.FormatUint(clientId, 10) + ">>>traffic>>>downlink"
		if c, _ := stats.GetOrRegisterCounter(h.statManager, name); c != nil {
			sizeStatWriter.Counter = c
		}
	} else if buWriter, ok := link.Writer.(*buf.EndpointOverrideWriter); ok {
		name := "user>>>" + strconv.FormatUint(clientId, 10) + ">>>traffic>>>downlink"
		if c, _ := stats.GetOrRegisterCounter(h.statManager, name); c != nil {
			if statWriter, ok := buWriter.Writer.(*dispatcher.SizeStatWriter); ok {
				statWriter.Counter = c
			}
		}
	}

	request := &protocol.RequestHeader{
		Version:  encoding.Version,
		User:     rec.PickUser(),
		Command:  command,
		Address:  target.Address,
		Port:     target.Port,
		ClientId: clientId,
	}

	account := request.User.Account.(*vless.MemoryAccount)

	requestAddons := &encoding.Addons{
		Flow: account.Flow,
	}

	var input *bytes.Reader
	var rawInput *bytes.Buffer
	allowUDP443 := false
	switch requestAddons.Flow {
	case vless.XRV + "-udp443":
		allowUDP443 = true
		requestAddons.Flow = requestAddons.Flow[:16]
		fallthrough
	case vless.XRV:
		ob.CanSpliceCopy = 2
		switch request.Command {
		case protocol.RequestCommandUDP:
			if !allowUDP443 && request.Port == 443 {
				return errors.New("XTLS rejected UDP/443 traffic").AtInfo()
			}
		case protocol.RequestCommandMux:
			fallthrough // let server break Mux connections that contain TCP requests
		case protocol.RequestCommandTCP:
			var t reflect.Type
			var p uintptr
			if commonConn, ok := conn.(*encryption.CommonConn); ok {
				if _, ok := commonConn.Conn.(*encryption.XorConn); ok || !proxy.IsRAWTransportWithoutSecurity(iConn) {
					ob.CanSpliceCopy = 3 // full-random xorConn / non-RAW transport / another securityConn should not be penetrated
				}
				t = reflect.TypeOf(commonConn).Elem()
				p = uintptr(unsafe.Pointer(commonConn))
			} else if tlsConn, ok := iConn.(*tls.Conn); ok {
				t = reflect.TypeOf(tlsConn.Conn).Elem()
				p = uintptr(unsafe.Pointer(tlsConn.Conn))
			} else if utlsConn, ok := iConn.(*tls.UConn); ok {
				t = reflect.TypeOf(utlsConn.Conn).Elem()
				p = uintptr(unsafe.Pointer(utlsConn.Conn))
			} else if realityConn, ok := iConn.(*reality.UConn); ok {
				t = reflect.TypeOf(realityConn.Conn).Elem()
				p = uintptr(unsafe.Pointer(realityConn.Conn))
			} else {
				return errors.New("XTLS only supports TLS and REALITY directly for now.").AtWarning()
			}
			i, _ := t.FieldByName("input")
			r, _ := t.FieldByName("rawInput")
			input = (*bytes.Reader)(unsafe.Pointer(p + i.Offset))
			rawInput = (*bytes.Buffer)(unsafe.Pointer(p + r.Offset))
		}
	default:
		ob.CanSpliceCopy = 3
	}

	var newCtx context.Context
	var newCancel context.CancelFunc
	if session.TimeoutOnlyFromContext(ctx) {
		newCtx, newCancel = context.WithCancel(context.Background())
	}

	sessionPolicy := h.policyManager.ForLevel(request.User.Level)
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, func() {
		cancel()
		if newCancel != nil {
			newCancel()
		}
	}, sessionPolicy.Timeouts.ConnectionIdle)

	clientReader := link.Reader // .(*pipe.Reader)
	clientWriter := link.Writer // .(*pipe.Writer)
	trafficState := proxy.NewTrafficState(account.ID.Bytes())
	if request.Command == protocol.RequestCommandUDP && (requestAddons.Flow == vless.XRV || (h.cone && request.Port != 53 && request.Port != 443)) {
		request.Command = protocol.RequestCommandMux
		request.Address = net.DomainAddress("v1.mux.cool")
		request.Port = net.Port(666)
	}

	postRequest := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

		bufferWriter := buf.NewBufferedWriter(buf.NewWriter(conn))
		if err := encoding.EncodeRequestHeader(bufferWriter, request, requestAddons); err != nil {
			return errors.New("failed to encode request header").Base(err).AtWarning()
		}

		// default: serverWriter := bufferWriter
		serverWriter := encoding.EncodeBodyAddons(bufferWriter, request, requestAddons, trafficState, true, ctx, conn, ob)
		if request.Command == protocol.RequestCommandMux && request.Port == 666 {
			serverWriter = xudp.NewPacketWriter(serverWriter, target, xudp.GetGlobalID(ctx))
		}
		timeoutReader, ok := clientReader.(buf.TimeoutReader)
		if ok {
			multiBuffer, err1 := timeoutReader.ReadMultiBufferTimeout(time.Millisecond * 500)
			if err1 == nil {
				if err := serverWriter.WriteMultiBuffer(multiBuffer); err != nil {
					return err // ...
				}
			} else if err1 != buf.ErrReadTimeout {
				return err1
			} else if requestAddons.Flow == vless.XRV {
				mb := make(buf.MultiBuffer, 1)
				errors.LogInfo(ctx, "Insert padding with empty content to camouflage VLESS header ", mb.Len())
				if err := serverWriter.WriteMultiBuffer(mb); err != nil {
					return err // ...
				}
			}
		} else {
			errors.LogDebug(ctx, "Reader is not timeout reader, will send out vless header separately from first payload")
		}
		// Flush; bufferWriter.WriteMultiBuffer now is bufferWriter.writer.WriteMultiBuffer
		if err := bufferWriter.SetBuffered(false); err != nil {
			return errors.New("failed to write A request payload").Base(err).AtWarning()
		}

		if requestAddons.Flow == vless.XRV {
			if tlsConn, ok := iConn.(*tls.Conn); ok {
				if tlsConn.ConnectionState().Version != gotls.VersionTLS13 {
					return errors.New(`failed to use `+requestAddons.Flow+`, found outer tls version `, tlsConn.ConnectionState().Version).AtWarning()
				}
			} else if utlsConn, ok := iConn.(*tls.UConn); ok {
				if utlsConn.ConnectionState().Version != utls.VersionTLS13 {
					return errors.New(`failed to use `+requestAddons.Flow+`, found outer tls version `, utlsConn.ConnectionState().Version).AtWarning()
				}
			}
		}
		err := buf.Copy(clientReader, serverWriter, buf.UpdateActivity(timer))
		if err != nil {
			return errors.New("failed to transfer request payload").Base(err).AtInfo()
		}

		// Indicates the end of request payload.
		switch requestAddons.Flow {
		default:
		}
		return nil
	}

	getResponse := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

		responseAddons, err := encoding.DecodeResponseHeader(conn, request)
		if err != nil {
			return errors.New("failed to decode response header").Base(err).AtInfo()
		}

		// default: serverReader := buf.NewReader(conn)
		serverReader := encoding.DecodeBodyAddons(conn, request, responseAddons)
		if requestAddons.Flow == vless.XRV {
			serverReader = proxy.NewVisionReader(serverReader, trafficState, false, ctx, conn, input, rawInput, ob)
		}
		if request.Command == protocol.RequestCommandMux && request.Port == 666 {
			if requestAddons.Flow == vless.XRV {
				serverReader = xudp.NewPacketReader(&buf.BufferedReader{Reader: serverReader})
			} else {
				serverReader = xudp.NewPacketReader(conn)
			}
		}

		if requestAddons.Flow == vless.XRV {
			err = encoding.XtlsRead(serverReader, clientWriter, timer, conn, trafficState, false, ctx)
		} else {
			// from serverReader.ReadMultiBuffer to clientWriter.WriteMultiBuffer
			err = buf.Copy(serverReader, clientWriter, buf.UpdateActivity(timer))
		}

		if err != nil {
			return errors.New("failed to transfer response payload").Base(err).AtInfo()
		}

		return nil
	}

	if newCtx != nil {
		ctx = newCtx
	}

	if err := task.Run(ctx, postRequest, task.OnSuccess(getResponse, task.Close(clientWriter))); err != nil {
		return errors.New("connection ends").Base(err).AtInfo()
	}

	return nil
}
