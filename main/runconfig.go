package main

import (
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/core"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/features/routing"
	"log"
)

func RunJsonConfig(json string) error {
	instance, err := core.StartInstance("json", []byte(json))
	if err != nil {
		log.Println("start instance error:", err.Error())
		return err
	}
	dispatcher := instance.GetFeature(routing.DispatcherType())

	log.Println("start instance success, status:", instance.IsRunning(), " dispatcher type:", dispatcher.Type())
	return nil
}
