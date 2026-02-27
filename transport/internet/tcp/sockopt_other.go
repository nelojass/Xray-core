//go:build !linux && !freebsd && !darwin
// +build !linux,!freebsd,!darwin

package tcp

import (
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport/internet/stat"
)

func GetOriginalDestination(conn stat.Connection) (net.Destination, error) {
	return net.Destination{}, nil
}
