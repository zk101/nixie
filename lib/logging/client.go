package logging

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"github.com/zk101/nixie/lib/logging/gelf"
)

// Client is an operational struct which embeds *zap.Logger
type Client struct {
	*zap.Logger
	config *Config
}

// NewClient returns a pointer to a new Client struct
func NewClient(conf *Config, serviceID string) (*Client, error) {
	if serviceID == "" {
		return nil, errors.New("service id must be set")
	}

	if conf == nil {
		c := DefaultConfig()
		conf = &c
	}

	var logger *zap.Logger
	var loglevel zapcore.Level
	encoderCfg := zap.NewProductionEncoderConfig()

	switch conf.LogLevel {
	case "info":
		loglevel = zap.InfoLevel
	case "error":
		loglevel = zap.ErrorLevel
	case "warn":
		loglevel = zap.WarnLevel
	default:
		loglevel = zap.DebugLevel
	}

	switch conf.Method {
	case "console":
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		if conf.ConsoleColour == true {
			encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		encoderCfg.EncodeTime = stdTimeEncoder

		core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), os.Stderr, loglevel)
		logger = zap.New(core).WithOptions(zap.Fields(zap.String("service_id", serviceID)))

	case "file":
		if conf.Filename == "" {
			return nil, errors.New("file logging can not log to an empty filename")
		}
		if conf.FileMaxAge < 1 || conf.FileMaxBackup < 1 || conf.FileMaxSize < 1 {
			return nil, errors.New("additional file options can not be 0 or less")
		}

		encoderCfg.TimeKey = "datetime"
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeTime = stdTimeEncoder

		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.Filename,
			MaxSize:    conf.FileMaxSize,
			MaxBackups: conf.FileMaxBackup,
			MaxAge:     conf.FileMaxAge,
		})

		core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), w, loglevel)
		logger = zap.New(core).WithOptions(zap.Fields(zap.String("service_id", serviceID)))

	case "gelf":
		if conf.GELFhost == "" {
			return nil, errors.New("gelf logging can not log to an empty host")
		}
		if conf.GELFport < 1 {
			return nil, errors.New("gelf port can not be 0 or less")
		}
		if conf.GELFmaxChunkSize < 1 {
			return nil, errors.New("gelf chunk size can not be 0 or less")
		}
		if conf.GELFappID == "" {
			return nil, errors.New("gelf logging can not log with an empty appid")
		}

		var compression int
		switch conf.GELFcompression {
		case "gzip":
			compression = gelf.CompressionGZip
		case "zlib":
			compression = gelf.CompressionZLib
		default:
			compression = gelf.CompressionNone
		}

		encoderCfg.TimeKey = gelf.TimeKey
		encoderCfg.MessageKey = gelf.MessageKey
		encoderCfg.EncodeLevel = gelfLevelEncoder
		encoderCfg.EncodeName = gelfNameEncoder

		w := zapcore.AddSync(gelf.New(gelf.Config{
			Host:         conf.GELFhost,
			Port:         conf.GELFport,
			MaxChunkSize: conf.GELFmaxChunkSize,
			Compression:  compression,
		}))

		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), w, loglevel)
		logger = zap.New(core).WithOptions(zap.Fields(
			zap.String(gelf.VersionTag, gelf.Version),
			zap.String(gelf.HostTag, conf.GELFappID),
			zap.String("service_id", serviceID),
		))

	default:
		return nil, errors.New("log method is not supported")
	}

	if conf.StackTrace == true {
		var stacktracelevel zapcore.Level

		switch conf.StackTraceLevel {
		case "info":
			stacktracelevel = zap.InfoLevel
		case "error":
			stacktracelevel = zap.ErrorLevel
		case "warn":
			stacktracelevel = zap.WarnLevel
		default:
			stacktracelevel = zap.DebugLevel
		}

		logger = logger.WithOptions(zap.AddStacktrace(stacktracelevel))
	}

	client := Client{logger, conf}

	return &client, nil
}

// EOF
