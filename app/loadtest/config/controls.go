package config

// Controls holds application control vars
type Controls struct {
	CAcertPath        string
	OverrideConsul    bool
	OverrideEnv       bool
	PrefixConsul      string
	PrefixEnv         string
	ServiceConsul     bool
	ServiceName       string
	ServiceAddrFilter string
	ServiceTags       string
}

// DefaultControls sets up default Controls
func DefaultControls() Controls {
	return Controls{
		ServiceAddrFilter: "^127\\.0\\.0\\.1",
	}
}

// EOF
