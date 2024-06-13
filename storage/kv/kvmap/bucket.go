package kvmap

import (
	"context"
	"fmt"

	"github.com/micromdm/nanolib/storage/kv"
)

// Get retrieves the value at key in the Go map.
// If key is not found then a wrapped ErrKeyNotFound will be returned.
func (s *KVMap) Get(_ context.Context, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	if !ok {
		// generate specific error type to comply with interface
		return nil, fmt.Errorf("%w: %s", kv.ErrKeyNotFound, key)
	}
	return v, nil
}

// Set sets key to value in the Go map.
func (s *KVMap) Set(_ context.Context, key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
	return nil
}

// Has checks that key is found in the Go map.
func (s *KVMap) Has(_ context.Context, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[key]
	return ok, nil
}

// Delete deletes key in the Go map.
func (s *KVMap) Delete(_ context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
	return nil
}
