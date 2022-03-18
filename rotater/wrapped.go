package rotater

import (
	"sync"
)

type lockedWriteSyncer struct {
	setOnceOnclose sync.Once
	onceOnclose    sync.Once
	onClose        func()
	sync.Mutex
	ws WriteSyncer
}

// NewLocked create a writer.
func NewLocked(ws WriteSyncer) RotateSyncer {
	if lws, ok := ws.(*lockedWriteSyncer); ok {
		// no need to layer on another lock
		return lws
	}
	return &lockedWriteSyncer{ws: ws}
}

func (s *lockedWriteSyncer) SetOnClose(closeFunc func()) {
	s.setOnceOnclose.Do(func() {
		s.onClose = closeFunc
	})
}

func (s *lockedWriteSyncer) Rotate() error {
	return s.Sync()
}

func (s *lockedWriteSyncer) Write(bs []byte) (int, error) {
	s.Lock()
	defer s.Unlock()
	return s.ws.Write(bs)
}

func (s *lockedWriteSyncer) Sync() error {
	s.Lock()
	defer s.Unlock()
	return s.ws.Sync()
}

func (s *lockedWriteSyncer) Close() error {
	s.Lock()
	defer s.Unlock()
	s.onceOnclose.Do(func() {
		if s.onClose != nil {
			s.onClose()
			s.onClose = nil
		}
	})
	return s.ws.Close()
}
