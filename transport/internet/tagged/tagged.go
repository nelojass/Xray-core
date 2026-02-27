package tagged

import (
	"context"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/features/routing"
)

type DialFunc func(ctx context.Context, dispatcher routing.Dispatcher, dest net.Destination, tag string) (net.Conn, error)

var Dialer DialFunc
