package conf_test

import (
	"testing"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/protocol"
	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/infra/conf"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/proxy/freedom"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport/internet"
)

func TestFreedomConfig(t *testing.T) {
	creator := func() Buildable {
		return new(FreedomConfig)
	}

	runMultiTestCase(t, []TestCase{
		{
			Input: `{
				"domainStrategy": "AsIs",
				"redirect": "127.0.0.1:3366",
				"userLevel": 1
			}`,
			Parser: loadJSON(creator),
			Output: &freedom.Config{
				DomainStrategy: internet.DomainStrategy_AS_IS,
				DestinationOverride: &freedom.DestinationOverride{
					Server: &protocol.ServerEndpoint{
						Address: &net.IPOrDomain{
							Address: &net.IPOrDomain_Ip{
								Ip: []byte{127, 0, 0, 1},
							},
						},
						Port: 3366,
					},
				},
				UserLevel: 1,
			},
		},
	})
}
