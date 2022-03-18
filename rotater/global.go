package rotater

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"os"
	"path/filepath"
)

var loggers RotateSyncerSet

// NewLogWriter create a RotateSyncer writer for logging.Logger.
func NewLogWriter(FileName string, logDir string, options ...Option) (RotateSyncer, error) {
	switch FileName {
	case "Stdout":
		return NewLocked(os.Stdout), nil
	case "Stderr":
		return NewLocked(os.Stderr), nil
	case "Null":
		return NewNull(), nil
	default:
		logpath := FileName
		if logDir != "" && !filepath.IsAbs(FileName) {
			logpath, _ = filepath.Abs(filepath.Join(logDir, FileName))
		}
		options = append(options, rotatelogs.WithLinkName(logpath))
		if logWriter, err := NewRotater(logpath+".%Y%m%d", options...); err != nil {
			return nil, err
		} else {
			loggers.Store(logWriter)
			logWriter.SetOnClose(func() { loggers.Delete(logWriter) })
			return logWriter, nil
		}
	}
}

// Rotate will rotate all registered logger .
func Rotate() {
	loggers.Range(func(rotater RotateSyncer) {
		_ = rotater.Sync()
		_ = rotater.Rotate()
	})
}

// Close will close all registered logger .
func Close() {
	loggers.Range(func(rotater RotateSyncer) {
		_ = rotater.Sync()
		_ = rotater.Rotate()
	})
}
