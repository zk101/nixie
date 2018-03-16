package user

import (
	"errors"
	"time"

	"github.com/zk101/nixie/models/couchbase/auth/presence"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// checkPresence tests if a presence document requires touching
func (d *Data) checkPresence(wsmsg *pbmsgwrap.MsgWrap) error {
	if wsmsg.GetMsgType() == pbmsgwrap.MsgType_MSG_NULL {
		return nil
	}

	if d.key == "" {
		return errors.New("connection key is not setup require null message to preceed with all further websocket messaging")
	}

	modelPresence := presence.New()

	if err := modelPresence.Fetch(d.cbPool, d.key, false); err != nil {
		return errors.New("get model presence failed")
	}

	if d.timeCheck > 0 && time.Now().Unix() >= d.timeCheck {
		if err := modelPresence.Touch(d.cbPool); err != nil {
			return err
		}
		d.io.Lock()
		d.timeCheck = time.Now().Unix() + 150
		d.io.Unlock()
	}

	return nil
}

// EOF
