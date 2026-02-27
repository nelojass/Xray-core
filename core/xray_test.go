package core_test

import (
	"testing"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/dispatcher"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/proxyman"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/protocol"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/serial"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/uuid"
	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/core"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/features/dns"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/features/dns/localdns"
	_ "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/main/distro/all"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/dokodemo"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/vmess"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/vmess/outbound"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/testing/servers/tcp"
	"google.golang.org/protobuf/proto"
)

func TestXrayDependency(t *testing.T) {
	instance := new(Instance)

	wait := make(chan bool, 1)
	instance.RequireFeatures(func(d dns.Client) {
		if d == nil {
			t.Error("expected dns client fulfilled, but actually nil")
		}
		wait <- true
	}, false)
	instance.AddFeature(localdns.New())
	<-wait
}

func TestXrayClose(t *testing.T) {
	port := tcp.PickPort()

	userID := uuid.New()
	config := &Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
		Inbound: []*InboundHandlerConfig{
			{
				ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
					PortList: &net.PortList{
						Range: []*net.PortRange{net.SinglePortRange(port)},
					},
					Listen: net.NewIPOrDomain(net.LocalHostIP),
				}),
				ProxySettings: serial.ToTypedMessage(&dokodemo.Config{
					Address:  net.NewIPOrDomain(net.LocalHostIP),
					Port:     uint32(0),
					Networks: []net.Network{net.Network_TCP},
				}),
			},
		},
		Outbound: []*OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&outbound.Config{
					Receiver: []*protocol.ServerEndpoint{
						{
							Address: net.NewIPOrDomain(net.LocalHostIP),
							Port:    uint32(0),
							User: []*protocol.User{
								{
									Account: serial.ToTypedMessage(&vmess.Account{
										Id: userID.String(),
									}),
								},
							},
						},
					},
				}),
			},
		},
	}

	cfgBytes, err := proto.Marshal(config)
	common.Must(err)

	server, err := StartInstance("protobuf", cfgBytes)
	common.Must(err)
	server.Close()
}
