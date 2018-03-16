package chat

import (
	"errors"
	"time"

	"github.com/zk101/nixie/app/loadtest/data"
)

// Wait sends a search request
func Wait(local *data.Local) error {
	if local.WS == nil {
		return errors.New("websocket not connected")
	}

	for key, value := range local.Peers {
		if local.WorkerID == key {
			continue
		}

		loopBreak := 0
		for {
			if value.Friends == false {
				if loopBreak > 10000 {
					return errors.New("chat wait gave up waiting for peers")
				}

				time.Sleep(time.Millisecond)
				loopBreak++
				continue
			}

			break
		}
	}

	return nil
}

// EOF
