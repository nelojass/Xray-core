package all

import (
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/main/commands/all/api"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/main/commands/all/convert"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/main/commands/all/tls"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/main/commands/base"
)

func init() {
	base.RootCommand.Commands = append(
		base.RootCommand.Commands,
		api.CmdAPI,
		convert.CmdConvert,
		tls.CmdTLS,
		cmdUUID,
		cmdX25519,
		cmdWG,
		cmdMLDSA65,
		cmdMLKEM768,
		cmdVLESSEnc,
	)
}
