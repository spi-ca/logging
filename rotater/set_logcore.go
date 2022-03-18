package rotater

import (
	"sync"
	"sync/atomic"
)

// RotateSyncerSet is registry of RotateSyncer
type RotateSyncerSet struct {
	storage sync.Map
}

// Delete deletes the value for a key.
func (s *RotateSyncerSet) Delete(key RotateSyncer) {
	s.storage.Delete(key)
}

// Exist returns whether value was found in the map.
func (s *RotateSyncerSet) Exist(key RotateSyncer) (ok bool) {
	_, ok = s.storage.Load(key)
	return
}

// SetNx returns false value was found in the map.
// Otherwise, it stores and returns true.
func (s *RotateSyncerSet) SetNx(key RotateSyncer) bool {
	_, exist := s.storage.LoadOrStore(key, 0)
	return !exist
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (s *RotateSyncerSet) Range(f func(key RotateSyncer)) {
	s.storage.Range(s.rangeWrap(f))
}

// Store sets the value for a key.
func (s *RotateSyncerSet) Store(key RotateSyncer) {
	s.storage.Store(key, 0)
}

// Len returns sizeof the map.
func (s *RotateSyncerSet) Len() int {
	var count uint64
	s.Range(func(conn RotateSyncer) {
		atomic.AddUint64(&count, 1)
	})
	return int(count)
}

func (s *RotateSyncerSet) rangeWrap(f func(key RotateSyncer)) func(key, value any) bool {
	ok := true
	return func(key, value any) bool {
		f(key.(RotateSyncer))
		return ok
	}
}
