package telemetry

import (
	"encoding/json"
	"errors"

	gwebsock "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"github.com/zk101/nixie/app/loadtest/data"
	mMsgwrap "github.com/zk101/nixie/models/protobuf/ws/msgwrap"
	mTelemetry "github.com/zk101/nixie/models/sql/telemetry/telemetry"
	pbmsgwrap "github.com/zk101/nixie/proto/ws/msgwrap"
)

// testTelemetryData
type testTelemetryData struct {
	Test1 string
	Test2 int
	Test3 bool
}

// Send sends a Telemetry message
func Send(local *data.Local) error {
	if local.WS == nil {
		return errors.New("websocket not connected")
	}

	testData := testTelemetryData{
		Test1: "test string",
		Test2: 12345,
		Test3: true,
	}

	telemetryData, err := json.Marshal(&testData)
	if err != nil {
		return err
	}

	telemetryModel := mTelemetry.New()
	telemetryModel.User = local.User
	telemetryModel.Version = "0.1alpha"
	telemetryModel.Client = "loadtest"
	telemetryModel.Data = string(telemetryData)

	msgTelemetry, err := telemetryModel.Pack()
	if err != nil {
		return err
	}

	msgTrack := uuid.NewV4().String()

	msgWrap, err := mMsgwrap.Pack(pbmsgwrap.MsgType_MSG_TELEMETRY, pbmsgwrap.MsgSec_SEC_SIGN, local.Key, local.Sign, local.Cipher, msgTrack, &msgTelemetry)
	if err != nil {
		return err
	}

	return local.WS.Write(gwebsock.BinaryMessage, *msgWrap)
}

// EOF
