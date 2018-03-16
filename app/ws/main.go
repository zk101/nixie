package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"github.com/zk101/nixie/app/ws/config"
	"github.com/zk101/nixie/app/ws/lib"
)

// deadliner struct is a wrapper around net.Conn to provide Deadline setting helpers
type deadliner struct {
	net.Conn
	t time.Duration
}

// Write sets a write deadline for Conn
func (d deadliner) Write(p []byte) (int, error) {
	if err := d.Conn.SetWriteDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Write(p)
}

// Read sets a read deadline for Conn
func (d deadliner) Read(p []byte) (int, error) {
	if err := d.Conn.SetReadDeadline(time.Now().Add(d.t)); err != nil {
		return 0, err
	}
	return d.Conn.Read(p)
}

// run is the main program loop
func run(core *lib.Core) error {
	pollConfig := netpoll.Config{OnWaitError: core.NetpollError}

	poller, err := netpoll.New(&pollConfig)
	if err != nil {
		return err
	}

	ioTimeout := time.Millisecond * time.Duration(core.Config.Controls.IOtimeout)

	handle := func(conn net.Conn) {
		safeConn := deadliner{conn, ioTimeout}

		_, err := ws.Upgrade(safeConn)
		if err != nil {
			core.Clients.Log.Sugar().Warnw("upgrade error", "service_id", core.Clients.ServiceID, "error", err.Error(), "remote_addr", conn.RemoteAddr().String())
			conn.Close()
			return
		}

		core.Clients.Log.Sugar().Debugw("websocket connection established", "service_id", core.Clients.ServiceID, "remote_addr", conn.RemoteAddr().String())

		user := core.Clients.Connection.Register(safeConn)
		desc := netpoll.Must(netpoll.HandleRead(conn))

		poller.Start(desc, func(ev netpoll.Event) {
			if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
				poller.Stop(desc)
				core.Clients.Connection.Remove(user)
				core.Clients.Log.Sugar().Debugw("connection closed", "service_id", core.Clients.ServiceID, "reason", "client disconnect", "remote_addr", conn.RemoteAddr().String())
				return
			}
			core.Clients.GoPool.Schedule(func() {
				if err := user.Receive(); err != nil {
					poller.Stop(desc)
					core.Clients.Connection.Remove(user)
					core.Clients.Log.Sugar().Warnw("connection closed", "service_id", core.Clients.ServiceID, "reason", err.Error(), "remote_addr", conn.RemoteAddr().String())
				}
			})
		})
	}

	var ln net.Listener
	address := core.Config.Controls.ListenAddr + ":" + strconv.Itoa(core.Config.Controls.ListenPort)

	if core.Config.Controls.SSL == true {
		cert, err := tls.LoadX509KeyPair(core.Config.Controls.SSLcert, core.Config.Controls.SSLkey)
		if err != nil {
			return err
		}

		config := tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		ln, err = tls.Listen("tcp", address, &config)
	} else {
		ln, err = net.Listen("tcp", address)
	}
	if err != nil {
		return err
	}

	core.Clients.Log.Sugar().Debugw("ws server started", "service_id", core.Clients.ServiceID, "address", ln.Addr().String())

	acceptDesc := netpoll.Must(netpoll.HandleListener(ln, netpoll.EventRead|netpoll.EventOneShot))
	accept := make(chan error, 1)

	poller.Start(acceptDesc, func(e netpoll.Event) {
		err := core.Clients.GoPool.ScheduleTimeout(time.Millisecond, func() {
			conn, err := ln.Accept()
			if err != nil {
				accept <- err
				return
			}

			accept <- nil
			handle(conn)
		})

		if err == nil {
			err = <-accept
		}

		if err != nil {
			core.Clients.Log.Sugar().Errorw("connection accept error", "service_id", core.Clients.ServiceID, "error", err.Error())
			time.Sleep(time.Millisecond * 5)
		}

		poller.Resume(acceptDesc)
	})

	<-core.Config.Signal.Channel
	core.Config.Signal.Run = false

	return nil
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
		core.Clients.Log.Sugar().Errorf("Appplication Run Loop failed: %s", err.Error())
		errorLevel = 1
	}

	if err := core.Stop(); err != nil {
		core.Clients.Log.Sugar().Errorf("Core Stop failed: %s", err.Error())
		errorLevel = 1
	}

	config.ShutdownClients(core.Clients)

	if err := core.DeleteQueues(); err != nil {
		log.Fatalf("Core Delete Queue failed: %s", err.Error())
	}

	os.Exit(errorLevel)
}

// EOF
