package kvtxn

import (
	"context"

	"github.com/micromdm/nanolib/storage/kv"
)

// Get retrieves value at key.
// A previously staged key may be returned.
func (b *KVTxn) Get(ctx context.Context, key string) ([]byte, error) {
	if !b.hasOp(key) {
		b.keyLock.RLock(key)
		defer b.keyLock.RUnlock(key)
	}
	if !b.autoCommit {
		b.stageLock.RLock()
		defer b.stageLock.RUnlock()
		if value, del, found := b.stageGet(key); found {
			if del {
				// found a stage operation that deleted this key
				return nil, kv.ErrKeyNotFound
			}
			return value, nil
		}
	}
	// fallback to underlying store
	return b.store.Get(ctx, key)
}

// Set sets key to value in the staged operations.
// This change may be auto-commited.
func (b *KVTxn) Set(ctx context.Context, key string, value []byte) error {
	if !b.hasOp(key) {
		b.keyLock.Lock(key)
	}
	b.stageLock.Lock()
	defer b.stageLock.Unlock()
	b.stageSet(key, value)
	if b.autoCommit {
		return b.stageCommit(ctx)
	}
	return nil
}

// Has checks that key can be found.
// A previously staged key may be returned.
func (b *KVTxn) Has(ctx context.Context, key string) (bool, error) {
	if !b.hasOp(key) {
		b.keyLock.RLock(key)
		defer b.keyLock.RUnlock(key)
	}
	if !b.autoCommit {
		b.stageLock.RLock()
		defer b.stageLock.RUnlock()
		if has, found := b.stageHas(key); found {
			return has, nil
		}
	}
	// fallback to underlying store
	return b.store.Has(ctx, key)
}

// Delete deletes key in the staged operations.
// This change may be auto-commited.
func (b *KVTxn) Delete(ctx context.Context, key string) error {
	if !b.hasOp(key) {
		b.keyLock.Lock(key)
	}
	b.stageLock.Lock()
	defer b.stageLock.Unlock()
	b.stageDelete(key)
	if b.autoCommit {
		return b.stageCommit(ctx)
	}
	return nil
}
