// Package kvtxn provides an in-memory transactional wrapper for KV stores.
// Note that underlying KV stores are assumed to not support
// multi-operation atomicity. Thus this wrapper cannot guarantee
// transaction atomicity, either.
package kvtxn

import (
	"context"
	"sync"

	"github.com/micromdm/nanolib/storage/kv"
)

// KeyLockManager works like sync.RWMutex but supports per-key locking.
type KeyLockManager interface {
	RLock(key string)
	RUnlock(key string)
	Lock(key string)
	Unlock(key string)
}

// keyOp is a staged operation for a key.
type keyOp struct {
	value []byte
	del   bool // if true this operation signifies a deletion (of a key)
}

// KVTxn is a key-value store wrapper that supports in-memory transactions.
// Note the underlying KV store can still be inconsistent—this wrapper
// does NOT gauarantee any commit atomicity.
//
// KVTxn maintains an in-memory "stage" for write operations
// per-transaction. These staged operations can be rolled-back or
// committed.
// The store uses key-based mutexes for the duration of transactions
// to try to maintain consistency.
type KVTxn struct {
	store       kv.KeysPrefixTraversingBucket
	stageLock   sync.RWMutex
	stageKeyOps map[string]keyOp
	keyLock     KeyLockManager
	autoCommit  bool
}

// New creates a new in-memory transacting key-value store that wraps store.
// Note that a single in-memory lock manager is created so transaction
// locking will only be scoped to this newly created store.
func New(store kv.KeysPrefixTraversingBucket) *KVTxn {
	// create a new store with auto-commit on.
	return new(store, NewInmemLockManager(), true)
}

// new is a helper for creating KVTxns that wraps store.
func new(store kv.KeysPrefixTraversingBucket, keyLock KeyLockManager, autoCommit bool) *KVTxn {
	if store == nil {
		panic("nil store")
	}
	if keyLock == nil {
		panic("nil key lock manager")
	}
	return &KVTxn{
		store:       store,
		stageKeyOps: make(map[string]keyOp),
		keyLock:     keyLock,
		autoCommit:  autoCommit,
	}
}

// stageGet retreives a key from the staged key operations.
func (b *KVTxn) stageGet(key string) (value []byte, del bool, found bool) {
	keyOp, ok := b.stageKeyOps[key]
	return keyOp.value, keyOp.del, ok
}

// stageSet sets a value for key in the staged key operations.
func (b *KVTxn) stageSet(key string, value []byte) {
	b.stageKeyOps[key] = keyOp{value: value}
}

// stageHas checks that a key can be found in the staged key operations.
func (b *KVTxn) stageHas(key string) (has, found bool) {
	keyOp, ok := b.stageKeyOps[key]
	return !keyOp.del, ok
}

// stageDelete stages a key deletion in the staged key operations.
func (b *KVTxn) stageDelete(key string) {
	b.stageKeyOps[key] = keyOp{del: true}
}

// stageReset resets the staged operations.
func (b *KVTxn) stageReset() {
	for k := range b.stageKeyOps {
		// make sure we unlock any keys in the stage
		b.keyLock.Unlock(k)
	}
	b.stageKeyOps = make(map[string]keyOp)
}

// hasOp checks if there is an operation staged for key.
// A read lock is obtained for the stage lookup.
func (b *KVTxn) hasOp(key string) (ok bool) {
	b.stageLock.RLock()
	_, ok = b.stageKeyOps[key]
	b.stageLock.RUnlock()
	return
}

// stageCommit commits (sends) the staged operations to the wrapped KV store.
func (b *KVTxn) stageCommit(ctx context.Context) error {
	var err error
	for key, op := range b.stageKeyOps {
		if err == nil {
			if op.del {
				err = b.store.Delete(ctx, key)
			} else {
				err = b.store.Set(ctx, key, op.value)
			}
		}
		b.keyLock.Unlock(key)
		// if we had no error, remove the operation
		if err == nil {
			delete(b.stageKeyOps, key)
		}
	}
	return err
}
