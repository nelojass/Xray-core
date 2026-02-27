package dokodemo

import (
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/net"
)

// GetPredefinedAddress returns the defined address from proto config. Null if address is not valid.
func (v *Config) GetPredefinedAddress() net.Address {
	addr := v.Address.AsAddress()
	if addr == nil {
		return nil
	}
	return addr
}
