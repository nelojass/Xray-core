package conf_test

import (
	"testing"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/infra/conf"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/dns"
)

func TestDnsProxyConfig(t *testing.T) {
	creator := func() Buildable {
		return new(DNSOutboundConfig)
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"address": "8.8.8.8",
				"port": 53,
				"network": "tcp"
			}`,
			Parser: loadJSON(creator),
			Output: &dns.Config{
				Server: &net.Endpoint{
					Network: net.Network_TCP,
					Address: net.NewIPOrDomain(net.IPAddress([]byte{8, 8, 8, 8})),
					Port:    53,
				},
			},
		},
	})
}
