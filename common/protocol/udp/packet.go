package udp

import (
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/buf"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
)

// Packet is a UDP packet together with its source and destination address.
type Packet struct {
	Payload *buf.Buffer
	Source  net.Destination
	Target  net.Destination
}
