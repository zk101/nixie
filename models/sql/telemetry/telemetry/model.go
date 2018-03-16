package telemetry

import "time"

// Model is an exact copy of the protobuf generated Telemetry struct minus the annotations.  This allows conversion to just work.
type Model struct {
	User    string
	Client  string
	Version string
	Date    int64
	Data    string
}

// New returns a pointer to an empty presence Model
func New() *Model {
	return &Model{
		Date: time.Now().Unix(),
	}
}

// EOF
