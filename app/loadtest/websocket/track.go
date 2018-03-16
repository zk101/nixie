package websocket

import (
	"errors"
	"time"
)

// Track holds data required for tracking a websocket message and returning data
type Track struct {
	Reply chan ReadData
}

// AddTrack is a helper function to add a Track to the trackMap and returns the Track struct
func (c *Client) AddTrack(key string) Track {
	c.trackMap[key] = Track{
		Reply: make(chan ReadData, 1),
	}

	return c.trackMap[key]
}

// ReadBlock waits for a reply or timeout on a Track struct
func (t *Track) ReadBlock() (ReadData, error) {
	msgData := ReadData{}
	breakLoop := 0
	runLoop := true

	for runLoop {
		select {
		case msgData = <-t.Reply:
			runLoop = false

		default:
			if breakLoop > 10000 {
				return msgData, errors.New("tracker readblock reply timeout")
			}
			breakLoop++
			time.Sleep(time.Millisecond)
		}
	}

	return msgData, nil
}

// EOF
