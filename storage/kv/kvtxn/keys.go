package kvtxn

import (
	"context"
)

// Keys returns all keys in the underlying key-value store merging with the operations stage.
// The returned keys have no ordering guaratees.
// The keys channel should be closed if cancel was provided and closed.
// Beware of deadlocks with underlying implementations.
// Note that key-based stage locks are not consulted.
func (b *KVTxn) Keys(ctx context.Context, cancel <-chan struct{}) <-chan string {
	return b.keysWithStagedKeys(b.store.Keys(ctx, cancel), cancel)
}

// Keys returns all keys starting with prefix in the underlying key-value store merging with the operations stage.
// The returned keys have no ordering guaratees.
// The keys channel should be closed if cancel was provided and closed.
// Beware of deadlocks with underlying implementations.
// Note that key-based stage locks are not consulted.
func (b *KVTxn) KeysPrefix(ctx context.Context, prefix string, cancel <-chan struct{}) <-chan string {
	return b.keysWithStagedKeys(b.store.KeysPrefix(ctx, prefix, cancel), cancel)
}

// stageKeys returns a slice of all staged keys.
// Keys that have a delete operation are not included if noDel is true.
func (b *KVTxn) stageKeys(skipDeleted bool) []string {
	var r []string
	for k, v := range b.stageKeyOps {
		if v.del && skipDeleted {
			continue
		}
		r = append(r, k)
	}
	return r
}

// keysWithStagedKeys returns a merged set from inKeys and keys from staged operations.
func (b *KVTxn) keysWithStagedKeys(inKeys <-chan string, cancel <-chan struct{}) <-chan string {
	r := make(chan string)
	go func() {
		defer close(r)
		for k := range inKeys {
			b.stageLock.RLock()
			_, found := b.stageHas(k)
			b.stageLock.RUnlock()
			if found {
				// skip this key, it's in the stage
				continue
			}
			select {
			case <-cancel:
				return
			case r <- k:
			}
		}
		b.stageLock.RLock()
		// retreive all of our staged keys (minus the staged deletions)
		for _, k := range b.stageKeys(true) {
			select {
			case <-cancel:
				return
			case r <- k:
			}
		}
		b.stageLock.RUnlock()
	}()
	return r
}
