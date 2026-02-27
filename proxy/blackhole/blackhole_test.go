package blackhole_test

import (
	"context"
	"testing"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/buf"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/serial"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/session"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/blackhole"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport/pipe"
)

func TestBlackholeHTTPResponse(t *testing.T) {
	ctx := session.ContextWithOutbounds(context.Background(), []*session.Outbound{{}})
	handler, err := blackhole.New(ctx, &blackhole.Config{
		Response: serial.ToTypedMessage(&blackhole.HTTPResponse{}),
	})
	common.Must(err)

	reader, writer := pipe.New(pipe.WithoutSizeLimit())

	var mb buf.MultiBuffer
	var rerr error
	go func() {
		b, e := reader.ReadMultiBuffer()
		mb = b
		rerr = e
	}()

	link := transport.Link{
		Reader: reader,
		Writer: writer,
	}
	common.Must(handler.Process(ctx, &link, nil))
	common.Must(rerr)
	if mb.IsEmpty() {
		t.Error("expect http response, but nothing")
	}
}
