package logging

import (
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

type (

	// Logger represents logging interface.
	Logger interface {
		// DPanic uses fmt.Sprint to construct and log a message. In development, the
		// logger then panics. (See DPanicLevel for details.)
		DPanic(args ...any)
		// DPanicf uses fmt.Sprintf to log a templated message. In development, the
		// logger then panics. (See DPanicLevel for details.)
		DPanicf(template string, args ...any)
		// Debug uses fmt.Sprint to construct and log a message.
		Debug(args ...any)
		// Debugf uses fmt.Sprintf to log a templated message.
		Debugf(template string, args ...any)
		// Error uses fmt.Sprint to construct and log a message.
		Error(args ...any)
		// Errorf uses fmt.Sprintf to log a templated message.
		Errorf(template string, args ...any)
		// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
		Fatal(args ...any)
		// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
		Fatalf(template string, args ...any)
		// Info uses fmt.Sprint to construct and log a message.
		Info(args ...any)
		// Infof uses fmt.Sprintf to log a templated message.
		Infof(template string, args ...any)
		// Named adds a sub-scope to the logger's name.
		Named(name string) Logger
		// Name returns logger name
		Name() string
		// Panic uses fmt.Sprint to construct and log a message, then panics.
		Panic(args ...any)
		// Panicf uses fmt.Sprintf to log a templated message, then panics.
		Panicf(template string, args ...any)
		// Sync flushes any buffered log entries.
		Sync() error
		// Warn uses fmt.Sprint to construct and log a message.
		Warn(args ...any)
		// Warnf uses fmt.Sprintf to log a templated message.
		Warnf(template string, args ...any)
		// ToStdLogAt returns *log.Logger which writes to supplied the logger at
		// required level.
		ToStdLogAt(level Level) (*log.Logger, error)
	}

	// HookLogger is a Logging interface with a hooking capability.
	HookLogger interface {
		Logger
		// SetHook specify a log entry wrapper.
		SetHook(hook LoggerHook) (err error)
	}

	// LoggerHook is an alias for the hooking function.
	LoggerHook = func(level Level, logger, message string, at time.Time) (err error)

	Level = zapcore.Level
)

const (
	// LevelDebug logs are typically voluminous, and are usually disabled in
	// production.
	LevelDebug = zapcore.DebugLevel
	// LevelInfo is the default logging priority.
	LevelInfo = zapcore.InfoLevel
	// LevelWarn logs are more important than Info, but don't need individual
	// human review.
	LevelWarn = zapcore.WarnLevel
	// LevelError logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	LevelError = zapcore.ErrorLevel
	// LevelDPanic logs are particularly important errors. In development the
	// logger panics after writing the message.
	LevelDPanic = zapcore.DPanicLevel
	// LevelPanic logs a message, then panics.
	LevelPanic = zapcore.PanicLevel
	// LevelFatal logs a message, then calls os.Exit(1).
	LevelFatal = zapcore.FatalLevel
)
