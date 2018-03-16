package logging

// Config struct stores configuration for logging
type Config struct {
	Method           string
	LogLevel         string
	ConsoleColour    bool
	StackTrace       bool
	StackTraceLevel  string
	Filename         string
	FileMaxSize      int
	FileMaxBackup    int
	FileMaxAge       int
	GELFhost         string
	GELFport         int
	GELFmaxChunkSize int
	GELFcompression  string
	GELFappID        string
}

// DefaultConfig returns a Config struct with default data
func DefaultConfig() Config {
	return Config{
		Method:   "console",
		LogLevel: "debug",
	}
}

// EOF
