package logging

import "go.uber.org/zap/zapcore"

var (

	// LogCommonFormat is a common log entry format.
	LogCommonFormat = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// LogOnlyMessageFormat is a reduced log entry format.
	LogOnlyMessageFormat = zapcore.EncoderConfig{
		TimeKey:       "",
		LevelKey:      "L",
		NameKey:       "",
		CallerKey:     "",
		MessageKey:    "M",
		StacktraceKey: "",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			switch l {
			case zapcore.DebugLevel:
				enc.AppendString("(-)")
			case zapcore.InfoLevel:
			case zapcore.WarnLevel:
				enc.AppendString("(*)")
			case zapcore.ErrorLevel:
				enc.AppendString("(!)")
			case zapcore.DPanicLevel:
				fallthrough
			case zapcore.PanicLevel:
				enc.AppendString("(!!)")
			case zapcore.FatalLevel:
				enc.AppendString("(!!!)")
			default:
				// nothing
			}
		},
	}
)
