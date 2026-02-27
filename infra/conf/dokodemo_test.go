package conf_test

import (
	"testing"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/infra/conf"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/dokodemo"
)

func TestDokodemoConfig(t *testing.T) {
	creator := func() Buildable {
		return new(DokodemoConfig)
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"address": "8.8.8.8",
				"port": 53,
				"network": "tcp",
				"followRedirect": true,
				"userLevel": 1
			}`,
			Parser: loadJSON(creator),
			Output: &dokodemo.Config{
				Address: &net.IPOrDomain{
					Address: &net.IPOrDomain_Ip{
						Ip: []byte{8, 8, 8, 8},
					},
				},
				Port:           53,
				Networks:       []net.Network{net.Network_TCP},
				FollowRedirect: true,
				UserLevel:      1,
			},
		},
	})
}
