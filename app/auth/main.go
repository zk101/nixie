package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/zk101/nixie/app/auth/config"
	"github.com/zk101/nixie/app/auth/lib"
)

func run(core *lib.Core) error {
	http.HandleFunc("/", core.HandlerBase)

	s := &http.Server{
		Addr:         ":" + strconv.Itoa(core.Config.HTTPD.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.SetKeepAlivesEnabled(core.Config.HTTPD.Keepalive)

	go func() {
		<-core.Config.Signal.Channel
		core.Config.Signal.Run = false
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		s.Shutdown(ctx)
	}()

	core.Clients.Log.Sugar().Debugw("auth server started", "port", strconv.Itoa(core.Config.HTTPD.Port))

	var err error
	if core.Config.HTTPD.SSL == true {
		err = s.ListenAndServeTLS(core.Config.HTTPD.SSLcert, core.Config.HTTPD.SSLkey)
	} else {
		err = s.ListenAndServe()
	}

	if err.Error() == "http: Server closed" {
		return nil
	}

	return err
}

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("Config Load failed: %s", err.Error())
	}

	core, err := lib.NewCore(conf)
	if err != nil {
		log.Fatalf("Application Core Services failed: %s\n", err.Error())
	}

	errorLevel := 0

	if err := run(core); err != nil {
		core.Clients.Log.Sugar().Errorf("Application Run Loop failed: %s", err.Error())
		errorLevel = 1
	}

	if err := core.Stop(); err != nil {
		core.Clients.Log.Sugar().Errorf("Core Stop failed: %s", err.Error())
		errorLevel = 1
	}

	config.ShutdownClients(core.Clients)

	os.Exit(errorLevel)
}

// EOF
