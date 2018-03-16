package logging

import (
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
	"github.com/zk101/nixie/lib/logging/gelf"
)

// gelfLevelEncoder calls the gelf ZapLevelToGelfLevel to switch zap levels to syslog levels
func gelfLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt(gelf.ZapLevelToGelfLevel(int32(l)))
}

// gelfNameEncoder sorts naming of structured data to be GELF compliant
func gelfNameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	switch loggerName {
	case "facility", "file", "full_message", "host":
		enc.AppendString(loggerName)
	case "line", "timestamp", "short_message", "version":
		enc.AppendString(loggerName)
	default:
		if fmt.Sprintf("%c", loggerName[1]) != "_" {
			enc.AppendString("_" + loggerName)
		} else {
			enc.AppendString(loggerName)
		}
	}
}

// stdTimeEncoder appends a yyyy/mm/dd hh:mm:ss for time format
func stdTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05"))
}

// EOF
