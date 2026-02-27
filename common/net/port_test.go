package net_test

import (
	"testing"

	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
)

func TestPortRangeContains(t *testing.T) {
	portRange := &PortRange{
		From: 53,
		To:   53,
	}

	if !portRange.Contains(Port(53)) {
		t.Error("expected port range containing 53, but actually not")
	}
}
