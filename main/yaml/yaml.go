package yaml

import (
	"context"
	"io"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/cmdarg"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/errors"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/core"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/infra/conf"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/infra/conf/serial"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/main/confloader"
)

func init() {
	common.Must(core.RegisterConfigLoader(&core.ConfigFormat{
		Name:      "YAML",
		Extension: []string{"yaml", "yml"},
		Loader: func(input interface{}) (*core.Config, error) {
			switch v := input.(type) {
			case cmdarg.Arg:
				cf := &conf.Config{}
				for i, arg := range v {
					errors.LogInfo(context.Background(), "Reading config: ", arg)
					r, err := confloader.LoadConfig(arg)
					if err != nil {
						return nil, errors.New("failed to read config: ", arg).Base(err)
					}
					c, err := serial.DecodeYAMLConfig(r)
					if err != nil {
						return nil, errors.New("failed to decode config: ", arg).Base(err)
					}
					if i == 0 {
						// This ensure even if the muti-json parser do not support a setting,
						// It is still respected automatically for the first configure file
						*cf = *c
						continue
					}
					cf.Override(c, arg)
				}
				return cf.Build()
			case io.Reader:
				return serial.LoadYAMLConfig(v)
			default:
				return nil, errors.New("unknown type")
			}
		},
	}))
}
