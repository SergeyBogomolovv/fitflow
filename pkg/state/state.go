package state

import (
	"sync"
)

type state struct {
	mu   sync.RWMutex
	data map[int64]any
}

type State interface {
	Set(key int64, value any)
	Get(key int64) any
	Delete(key int64)
}

func NewState() State {
	return &state{data: make(map[int64]any)}
}

func (s *state) Set(key int64, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *state) Get(key int64) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[key]
}

func (s *state) Delete(key int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
