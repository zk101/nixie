package telemetry

import (
	"errors"

	"github.com/golang/protobuf/proto"
	pbTelemetry "github.com/zk101/nixie/proto/telemetry"
)

// Pack a telemetry message
func (m *Model) Pack() ([]byte, error) {
	if m.User == "" {
		return []byte{}, errors.New("user can not be empty")
	}

	if m.Client == "" {
		return []byte{}, errors.New("client can not be empty")
	}

	if m.Version == "" {
		return []byte{}, errors.New("version can not be empty")
	}

	telemetryData := pbTelemetry.TelemetryMsg(*m)
	msg, err := proto.Marshal(&telemetryData)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

// EOF
