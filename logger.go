// Package logging is convenient wrapper of zap logger.
// It provides grouping, rotation.
package logging // import "github.com/spi-ca/logging"

import (
	"go.uber.org/zap"
	"log"
)

type (
	// logger container
	loggerImpl struct {
		zap.SugaredLogger
		name string
	}
)

// Named adds a sub-scope to the logger's name.
func (l *loggerImpl) Named(name string) Logger {
	return &loggerImpl{
		SugaredLogger: *l.SugaredLogger.Named(name),
		name:          l.name + "." + name,
	}
}

// Name returns logger name
func (l *loggerImpl) Name() string {
	return l.name
}

// NewStdLogAt returns *log.Logger which writes to supplied the logger at
// required level.
func (l *loggerImpl) ToStdLogAt(level Level) (*log.Logger, error) {
	return zap.NewStdLogAt(l.Desugar(), level)
}
