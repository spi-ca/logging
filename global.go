package logging

import (
	"errors"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spi-ca/logging/rotater"
	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
)

var defaultErrorOutputOptions []zap.Option

func init() {
	if _, err := zap.RedirectStdLogAt(zap.L(), zapcore.DebugLevel); err != nil {
		panic(err)
	}
}

// New create a sub-logger from root logger with specified name.
func New(parent *zap.SugaredLogger, moduleName string, options ...zap.Option) *zap.SugaredLogger {
	var subLogger *zap.Logger
	if parent == nil {
		subLogger = zap.L().Named(moduleName)
	} else {
		subLogger = parent.Desugar().Named(moduleName)
	}

	subLogger.WithOptions(options...)

	return subLogger.Sugar()
}

// NewOtherLogger create a seperated-root-logger.
func NewOtherLogger(
	formatter zapcore.Encoder,
	moduleName, logFilename, logDir string,
	rotateOption []rotater.Option,
	logLevel zapcore.Level,
	fields ...zapcore.Field,
) (logger *zap.SugaredLogger, closer func() error, err error) {
	loglevel := zap.NewAtomicLevelAt(logLevel)
	logWriter, err := rotater.NewLogWriter(logFilename, logDir, rotateOption...)
	if err != nil {
		return
	}
	core := zapcore.NewCore(formatter, logWriter, loglevel)
	closer = logWriter.Close
	logger = zap.New(core, defaultErrorOutputOptions...).
		Named(moduleName).With(fields...).Sugar()
	return
}

// NewOtherLoggerWithOption create a seperated-root-logger with zap-logger option.
func NewOtherLoggerWithOption(
	formatter zapcore.Encoder,
	moduleName, logFilename, logDir string,
	rotateOption []rotater.Option,
	logLevel zapcore.Level,
	options []zap.Option,
	fields ...zapcore.Field,
) (logger *zap.SugaredLogger, closer func() error, err error) {
	loglevel := zap.NewAtomicLevelAt(logLevel)
	logWriter, err := rotater.NewLogWriter(logFilename, logDir, rotateOption...)
	if err != nil {
		return
	}
	core := zapcore.NewCore(formatter, logWriter, loglevel)
	closer = logWriter.Close
	options = append(defaultErrorOutputOptions, options...)
	logger = zap.New(core, options...).
		Named(moduleName).With(fields...).Sugar()
	return
}

// ReplaceGlobalHookLogger replaces log.Default() logger
func ReplaceGlobalHookLogger(name string, verbose bool, maxBackup uint, loggingDirectory, filename string, loggerLevel Level, simple bool) (logger HookLogger, canceler func(), err error) {
	var (
		logWrapper loggerWithHookImpl
		formatter  zapcore.Encoder
	)
	if simple {
		formatter = zapcore.NewConsoleEncoder(LogOnlyMessageFormat)
	} else {
		formatter = zapcore.NewConsoleEncoder(LogCommonFormat)
	}

	var zapLoger *zap.SugaredLogger
	// 전역 로거 초기화
	zapLoger, canceler, err = replaceGlobalLogger(
		verbose,
		formatter,
		name,
		filename,
		loggingDirectory,
		[]rotater.Option{
			rotatelogs.WithMaxAge(-1),
			rotatelogs.WithRotationCount(maxBackup),
		},
		loggerLevel,
		zap.Hooks(logWrapper.hook),
	)
	defer func() {
		if err != nil && canceler != nil {
			canceler()
			canceler = nil
		}
	}()
	if err != nil {
		//do nothing
	} else if zapLoger == nil {
		err = errors.New("not initialized")
	} else {
		logWrapper.SugaredLogger = *zapLoger
		logWrapper.name = name
		logger = &logWrapper
	}
	return
}

func replaceGlobalLogger(
	verbose bool,
	formatter zapcore.Encoder,
	mainLogName, logFilename, logDir string,
	rotateOption []rotater.Option,
	logLevel zapcore.Level,
	additionalOptions ...zap.Option,
) (logger *zap.SugaredLogger, back func(), err error) {
	level := zap.NewAtomicLevelAt(logLevel)

	var defaultWriter rotater.RotateSyncer
	if defaultWriter, err = rotater.NewLogWriter(logFilename, logDir, rotateOption...); err != nil {
		return
	}
	if defaultErrorOutputOptions == nil {
		defaultErrorOutputOptions = []zap.Option{zap.ErrorOutput(defaultWriter)}
	}
	options := defaultErrorOutputOptions
	if verbose {
		options = append(options, zap.AddStacktrace(zap.NewAtomicLevelAt(zap.PanicLevel)))
	}
	// reset log option slice
	options = append(options, additionalOptions...)
	log := zap.New(zapcore.NewCore(formatter, defaultWriter, level), options...).Named(mainLogName)

	var (
		closers []func()
		closer  = func() {
			for i := len(closers) - 1; i >= 0; i-- {
				closers[i]()
			}
		}
	)
	defer func() {
		if err != nil {
			closer()
		}
	}()
	closers = append(closers, zap.ReplaceGlobals(log))

	var rollback func()
	if rollback, err = zap.RedirectStdLogAt(log, zapcore.DebugLevel); err != nil {
		return
	}
	closers = append(closers, rollback)

	return log.Sugar(), closer, nil
}
