package websocket

import (
	"github.com/golang/protobuf/proto"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// workerRead is a goroutine to listen to a websocket
func (c *Client) workerRead() {
	if c.conn == nil {
		return
	}

	c.wg.Add(1)

	readChan := make(chan ReadData, 1)

	go func() {
		for {
			if c.conn == nil {
				break
			}

			data := ReadData{}
			data.MsgType, data.MsgData, data.Err = c.conn.ReadMessage()

			readChan <- data
		}
	}()

	for c.run {
		select {
		case <-c.quitChan:
			break

		case data := <-readChan:
			if data.Err != nil {
				c.log.Sugar().Errorw("websocket read error", "error", data.Err.Error())
				break
			}

			wsmsg := pbmsgwrap.MsgWrap{}
			if err := proto.Unmarshal(data.MsgData, &wsmsg); err != nil {
				c.log.Sugar().Errorw("websocket read error", "error", err.Error())
				break
			}

			if track, ok := c.trackMap[wsmsg.GetMsgTrack()]; ok {
				track.Reply <- data
				delete(c.trackMap, wsmsg.GetMsgTrack())
				break
			}

			if autoAction, ok := c.autoMap[pbmsgwrap.MsgType_name[int32(wsmsg.GetMsgType())]]; ok {
				if err := autoAction(data.MsgData); err != nil {
					c.log.Sugar().Errorw("websocket autoaction process error", "autoaction", pbmsgwrap.MsgType_name[int32(wsmsg.GetMsgType())], "error", err.Error())
				}
			}
		}
	}

	c.wg.Done()
}

// workerWrite is a gorouting to write to a websocket
func (c *Client) workerWrite() {
	if c.conn == nil {
		return
	}

	c.wg.Add(1)

	for c.run {
		select {
		case <-c.quitChan:
			break

		case data := <-c.writeChan:
			if err := c.conn.WriteMessage(data.msgType, data.msgData); err != nil {
				data.errChan <- err
			}

			data.errChan <- nil

			// DEBUG
			//fmt.Printf("Write: %+v\n", data)
		}
	}

	c.wg.Done()
}

// EOF
