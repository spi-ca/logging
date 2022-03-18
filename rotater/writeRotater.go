package rotater

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"sync"
)

type rotateSyncer struct {
	setOnceOnclose sync.Once
	onceOnclose    sync.Once
	onClose        func()
	*rotatelogs.RotateLogs
}

// NewRotater create a RotateSyncer writer.
func NewRotater(filename string, options ...Option) (RotateSyncer, error) {
	if rotateLogger, err := rotatelogs.New(filename, options...); err != nil {
		return nil, err
	} else {
		return &rotateSyncer{RotateLogs: rotateLogger}, nil
	}
}
func (r *rotateSyncer) SetOnClose(closeFunc func()) {
	r.setOnceOnclose.Do(func() {
		r.onClose = closeFunc
	})
}

func (r *rotateSyncer) Rotate() error {
	return r.RotateLogs.Rotate()
}
func (r *rotateSyncer) Close() error {
	r.onceOnclose.Do(func() {
		if r.onClose != nil {
			r.onClose()
			r.onClose = nil
		}
	})
	return r.RotateLogs.Close()
}

func (r *rotateSyncer) Sync() error {
	return nil
}

func (s *rotateSyncer) Write(bs []byte) (int, error) {
	return s.RotateLogs.Write(bs)
}
