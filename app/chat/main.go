package main

import (
	"log"
	"os"

	"github.com/zk101/nixie/app/chat/config"
	"github.com/zk101/nixie/app/chat/lib"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("Config Load failed: %s", err.Error())
	}

	core, err := lib.NewCore(conf)
	if err != nil {
		log.Fatalf("NewCore failed: %s\n", err.Error())
	}

	errorLevel := 0

	<-conf.Signal.Channel

	if err := core.Stop(); err != nil {
		core.Clients.Log.Sugar().Errorf("Core Stop failed: %s", err.Error())
		errorLevel = 1
	}

	config.ShutdownClients(core.Clients)

	os.Exit(errorLevel)
}

// EOF
