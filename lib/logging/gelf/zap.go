package gelf

// ZapLevelToGelfLevel maps the zap log levels to the syslog severity levels used for gelf.
// See https://en.wikipedia.org/wiki/Syslog for details.
func ZapLevelToGelfLevel(l int32) int {
	switch l {
	// DebugLevel
	case -1:
		return 7
	// InfoLevel
	case 0:
		return 6
	// WarnLevel
	case 1:
		return 4
	// ErrorLevel
	case 2:
		return 3
	// DPanicLevel
	// PanicLevel
	// FatalLevel
	case 3, 4, 5:
		return 0
	}
	// default, should never happen
	return 1
}

// EOF
