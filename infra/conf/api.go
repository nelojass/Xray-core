package conf

import (
	"strings"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/commander"
	loggerservice "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/log/command"
	observatoryservice "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/observatory/command"
	handlerservice "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/proxyman/command"
	routerservice "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/router/command"
	statsservice "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/stats/command"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/errors"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/serial"
)

type APIConfig struct {
	Tag      string   `json:"tag"`
	Listen   string   `json:"listen"`
	Services []string `json:"services"`
}

func (c *APIConfig) Build() (*commander.Config, error) {
	if c.Tag == "" {
		return nil, errors.New("API tag can't be empty.")
	}

	services := make([]*serial.TypedMessage, 0, 16)
	for _, s := range c.Services {
		switch strings.ToLower(s) {
		case "reflectionservice":
			services = append(services, serial.ToTypedMessage(&commander.ReflectionConfig{}))
		case "handlerservice":
			services = append(services, serial.ToTypedMessage(&handlerservice.Config{}))
		case "loggerservice":
			services = append(services, serial.ToTypedMessage(&loggerservice.Config{}))
		case "statsservice":
			services = append(services, serial.ToTypedMessage(&statsservice.Config{}))
		case "observatoryservice":
			services = append(services, serial.ToTypedMessage(&observatoryservice.Config{}))
		case "routingservice":
			services = append(services, serial.ToTypedMessage(&routerservice.Config{}))
		}
	}

	return &commander.Config{
		Tag:     c.Tag,
		Listen:  c.Listen,
		Service: services,
	}, nil
}
