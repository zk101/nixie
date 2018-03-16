package config

// Controls holds application control vars
type Controls struct {
	CAcertPath        string
	IOtimeout         int
	ListenAddr        string
	ListenPort        int
	SSL               bool
	SSLkey            string
	SSLcert           string
	OverrideConsul    bool
	OverrideEnv       bool
	PrefixConsul      string
	PrefixEnv         string
	PresenceExpiry    uint32
	QueueSize         int
	ServiceConsul     bool
	ServiceName       string
	ServiceAddrFilter string
	ServiceTags       string
	WorkerCount       int
}

// DefaultControls sets up default Controls
func DefaultControls() Controls {
	return Controls{
		IOtimeout:         5,
		ListenPort:        10000,
		PresenceExpiry:    300,
		QueueSize:         20,
		ServiceAddrFilter: "^127\\.0\\.0\\.1",
		WorkerCount:       10,
	}
}

// EOF
