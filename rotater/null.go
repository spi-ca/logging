package rotater

import (
	"sync"
)

type nullWriteSyncer struct {
	setOnceOnclose sync.Once
	onceOnclose    sync.Once
	onClose        func()
}

// NewNull create a blackhole writer.
func NewNull() RotateSyncer {
	return &nullWriteSyncer{}
}

func (s *nullWriteSyncer) SetOnClose(closeFunc func()) {
	s.setOnceOnclose.Do(func() {
		s.onClose = closeFunc
	})
}

func (s *nullWriteSyncer) Rotate() error                { return nil }
func (s *nullWriteSyncer) Write(bs []byte) (int, error) { return len(bs), nil }
func (s *nullWriteSyncer) Sync() error                  { return nil }
func (s *nullWriteSyncer) Close() error {
	s.onceOnclose.Do(func() {
		if s.onClose != nil {
			s.onClose()
			s.onClose = nil
		}
	})
	return nil
}
