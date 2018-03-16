package config

// Controls holds application control vars
type Controls struct {
	AuthTimeout       uint32
	CAcertPath        string
	OverrideConsul    bool
	OverrideEnv       bool
	PrefixConsul      string
	PrefixEnv         string
	PresenceExpiry    uint32
	ServiceConsul     bool
	ServiceName       string
	ServiceAddrFilter string
	ServiceTags       string
}

// DefaultControls sets up default Controls
func DefaultControls() Controls {
	return Controls{
		AuthTimeout:       5,
		PresenceExpiry:    300,
		ServiceAddrFilter: "^127\\.0\\.0\\.1",
	}
}

// EOF
