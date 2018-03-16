package gelf

const (
	// VersionTag is mandatory
	VersionTag = "version"
	// Version of supported gelf format
	Version = "1.1"
	// HostTag is mandatory
	HostTag = "host"
	// MessageKey provides the key value for gelf message field
	MessageKey = "short_message"
	// TimeKey provides the key value for gelf time field
	TimeKey = "timestamp"
)

// Compression constants
const (
	CompressionNone = 0
	CompressionGZip = 1
	CompressionZLib = 2
)

// EOF
