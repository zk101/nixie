package telemetry

import (
	"github.com/golang/protobuf/proto"
	pbTelemetry "github.com/zk101/nixie/proto/telemetry"
)

// Unpack a telemetry message
func (m *Model) Unpack(data []byte) error {
	telemetryPB := pbTelemetry.TelemetryMsg{}

	if err := proto.Unmarshal(data, &telemetryPB); err != nil {
		return err
	}

	*m = Model(telemetryPB)

	return nil
}

// EOF
