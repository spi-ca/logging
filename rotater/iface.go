package rotater

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
)

type (
	// WriteSyncer is a WriteCloser interface with synchronize capability.
	WriteSyncer interface {
		io.WriteCloser
		Sync() error
	}

	// RotateSyncer is a WriteSyncer interface with file rotate capability.
	RotateSyncer interface {
		WriteSyncer
		SetOnClose(func())
		Rotate() error
	}

	// Option is an alias for the rotatelogs.Option
	Option = rotatelogs.Option
)
