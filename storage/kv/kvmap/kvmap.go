// Package kvmap implements an in-memory key-value store backed by a Go map.
package kvmap

import (
	"sync"
)

// KVMap is an in-memory key-value store backed by a Go map.
type KVMap struct {
	mu sync.RWMutex
	m  map[string][]byte
}

// New creates a new in-memory key-value store.
func New() *KVMap {
	return &KVMap{m: make(map[string][]byte)}
}
