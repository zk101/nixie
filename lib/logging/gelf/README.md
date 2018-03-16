#### GELF Syncer for Zap

This code is modified from https://github.com/c-atarella/zap-graylog-writer, who in turn based the code from https://github.com/robertkowalski/graylog-golang before being mangled by zk :-)

http://docs.graylog.org/en/2.1/pages/gelf.html
https://en.wikipedia.org/wiki/Syslog

#### Chunk Size
This will depend on the network MTU.  So either 1420, or 8154 are good starting points.  This is the MTU minus the frame.

#### Usage

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = gelf.TimeKey
		encoderCfg.MessageKey = gelf.MessageKey
		encoderCfg.EncodeLevel = gelfLevelEncoder
		encoderCfg.EncodeName = gelfNameEncoder

		w := zapcore.AddSync(gelf.New(gelf.DefaultConfig("127.0.0.1")))

		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), w, loglevel)
		logger = zap.New(core).WithOptions(zap.Fields(
			zap.String(gelf.VersionTag, gelf.Version),
			zap.String(gelf.HostTag, "App ID or Host ID"),
		))

#### Usage Notes

The above Usage references two functions which override internal zap functions.  They are listed below for reference.  They are not included in the gelf library as I did not want the gelf-zap lib to be dependant on the zap lib to be dependant on the gelf-zap... Potential recursive dependancy.

#### gelfLevelEncoder

		func gelfLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt(gelf.ZapLevelToGelfLevel(int32(l)))
		}

#### gelfNameEncoder

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
