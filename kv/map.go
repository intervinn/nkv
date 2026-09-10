package kv

import (
	"sync"
)

type Entry struct {
	V any
}

type Map struct {
	mu sync.RWMutex
	m  map[string]Entry
}

func NewMap() *Map {
	return &Map{
		m: map[string]Entry{},
	}
}

func (m *Map) Put(k string, v any) {
	m.mu.Lock()
	m.m[k] = Entry{
		V: v,
	}
	m.mu.Unlock()
}

func (m *Map) Get(k string) any {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.m[k]
}

func (m *Map) Has(k string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.m[k]
	return ok
}
