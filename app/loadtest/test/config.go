package test

// Config holds values used to control a Test Profiles execution
type Config struct {
	NumWorkers uint32
	PairCount  uint32
	RPS        uint32
	TestSelect string
	TargetAuth string
}

// DefaultConfig returns a ptr to a Config struct prepared
func DefaultConfig() Config {
	return Config{
		NumWorkers: 1,
		PairCount:  1,
		RPS:        10,
		TestSelect: "Test1",
		TargetAuth: "http://localhost:10001",
	}
}

// EOF
