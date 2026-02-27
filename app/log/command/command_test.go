package command_test

import (
	"context"
	"testing"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/dispatcher"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/log"
	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/log/command"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/proxyman"
	_ "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/proxyman/inbound"
	_ "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/proxyman/outbound"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/serial"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/core"
)

func TestLoggerRestart(t *testing.T) {
	v, err := core.New(&core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	})
	common.Must(err)
	common.Must(v.Start())

	server := &LoggerServer{
		V: v,
	}
	common.Must2(server.RestartLogger(context.Background(), &RestartLoggerRequest{}))
}
