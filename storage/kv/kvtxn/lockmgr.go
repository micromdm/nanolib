package kvtxn

import (
	"sync"
)

// InmemLockManager is a lock manager that supports locking on keys (strings).
// In-memory native map based.
type InmemLockManager struct {
	locks    map[string]*sync.RWMutex
	counters map[string]int
	m        sync.Mutex
}

// NewInmemLockManager creates a new key lock manager.
func NewInmemLockManager() *InmemLockManager {
	return &InmemLockManager{
		locks:    make(map[string]*sync.RWMutex),
		counters: make(map[string]int),
	}
}

// RLock locks key lock in klm for reading.
// RLock on a sync.RWMutex is called under the hood.
func (klm *InmemLockManager) RLock(key string) {
	klm.m.Lock()

	lock, ok := klm.locks[key]
	if !ok || lock == nil {
		lock = &sync.RWMutex{}
		klm.locks[key] = lock
	}

	klm.counters[key]++

	klm.m.Unlock()

	lock.RLock()
}

// RUnlock undoes a single RLock call for key in klm.
// RUnlock on a sync.RWMutex is called under the hood.
func (klm *InmemLockManager) RUnlock(key string) {
	klm.m.Lock()
	defer klm.m.Unlock()

	lock, ok := klm.locks[key]
	if !ok || lock == nil {
		// no lock present, remove the keys anyway
		delete(klm.counters, key)
		delete(klm.locks, key)
		return
	}

	klm.counters[key]--
	lock.RUnlock()

	if klm.counters[key] <= 0 {
		delete(klm.counters, key)
		delete(klm.locks, key)
	}
}

// Lock locks key for writing in klm.
// Lock on a sync.RWMutex is called under the hood.
func (klm *InmemLockManager) Lock(key string) {
	klm.m.Lock()

	lock, ok := klm.locks[key]
	if !ok || lock == nil {
		lock = &sync.RWMutex{}
		klm.locks[key] = lock
	}

	klm.counters[key]++

	klm.m.Unlock()

	lock.Lock()
}

// Unlock unlocks key for writing in klm.
// Unlock on a sync.RWMutex is called under the hood.
func (klm *InmemLockManager) Unlock(key string) {
	klm.m.Lock()
	defer klm.m.Unlock()

	lock, ok := klm.locks[key]
	if !ok || lock == nil {
		// no lock present, remove the keys anyway
		delete(klm.counters, key)
		delete(klm.locks, key)
		return
	}

	klm.counters[key]--
	lock.Unlock()

	if klm.counters[key] <= 0 {
		delete(klm.counters, key)
		delete(klm.locks, key)
	}
}
