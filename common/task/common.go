package task

import "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"

// Close returns a func() that closes v.
func Close(v interface{}) func() error {
	return func() error {
		return common.Close(v)
	}
}
