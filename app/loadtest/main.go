package main

import (
	"errors"
	"log"
	"math"
	"os"

	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/config"
	"github.com/zk101/nixie/app/loadtest/data"
	"github.com/zk101/nixie/app/loadtest/lib"
	"github.com/zk101/nixie/app/loadtest/test/profiles"
	"github.com/zk101/nixie/app/loadtest/worker"
)

func run(core *lib.Core) error {
	if core.Config.Test.RPS < 1 {
		return errors.New("rps needs to be 1 for greater")
	}

	if core.Config.Test.NumWorkers < 1 {
		return errors.New("number of workers needs to be 1 for greater")
	}

	if core.Config.Test.PairCount < 1 {
		return errors.New("pair count needs to be 1 for greater")
	}

	if core.Config.Test.PairCount > core.Config.Test.NumWorkers {
		return errors.New("pair count needs to be less than or equal too the number of workers")
	}

	if mod := math.Mod(float64(core.Config.Test.NumWorkers), float64(core.Config.Test.PairCount)); mod != 0 {
		return errors.New("number of workers needs to be divisible by the number of peers")
	}

	testList := profiles.GetTestProfiles(core)

	if _, check := testList[core.Config.Test.TestSelect]; check == false {
		return errors.New("selected test does not exist")
	}

	peerData := make(map[string]*data.Peer, 0)
	localData := make(map[string]*data.Local, 0)
	for x := uint32(0); x < core.Config.Test.NumWorkers; x++ {
		workerUUID := uuid.NewV4().String()
		localData[workerUUID] = data.CreateLocal()
		localData[workerUUID].WorkerID = workerUUID
		peerData[workerUUID] = &data.Peer{}

		if uint32(len(peerData)) == core.Config.Test.PairCount {
			for key := range peerData {
				localData[key].Peers = peerData
			}

			peerData = make(map[string]*data.Peer, 0)
		}
	}

	workerConf, err := worker.NewClient(core, testList[core.Config.Test.TestSelect])
	if err != nil {
		return err
	}

	for _, value := range localData {
		go workerConf.Worker(value)
	}

	<-core.Config.Signal.Channel
	workerConf.SwitchRunFlag()
	workerConf.WG.Wait()

	return nil
}

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
