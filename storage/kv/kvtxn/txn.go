package kvtxn

import (
	"context"

	"github.com/micromdm/nanolib/storage/kv"
)

// Commit sends the staged operations to the wrapped KV store.
// Staged key write locks are unlocked.
//
// Note that this is a layer over a non-transactional KV store. Thus it
// does not support "atomic" commits. Some staged operations may fail
// and will return early leaving the underlying KV store in an
// inconsistent state (e.g. with staged operations half-applied).
// Note also that if there is an error then the stage will contain the
// errored and remaining operations (which could theorectically be
// re-tried with another commit attempt).
func (b *KVTxn) Commit(ctx context.Context) error {
	b.stageLock.Lock()
	defer b.stageLock.Unlock()
	return b.stageCommit(ctx)
}

// Rollback resets (removes) the staged operations and unlocks staged locks.
func (b *KVTxn) Rollback(context.Context) error {
	b.stageLock.Lock()
	defer b.stageLock.Unlock()
	// discard any transaction operations
	b.stageReset()
	return nil
}

// BeginKeysPrefixTraversingBucketTxn creates a new in-memory transacting key-value store that wraps the same store that b wraps.
// Auto-commit is turned off for the new store (allowing staged operations).
func (b *KVTxn) BeginKeysPrefixTraversingBucketTxn(context.Context) (kv.KeysPrefixTraversingBucketTxnCompleter, error) {
	return new(b.store, b.keyLock, false), nil
}

// BeginCRUDBucketTxn creates a new in-memory transacting key-value store that wraps the same store that b wraps.
// Auto-commit is turned off for the new store (allowing staged operations).
func (b *KVTxn) BeginCRUDBucketTxn(context.Context) (kv.CRUDBucketTxnCompleter, error) {
	return new(b.store, b.keyLock, false), nil
}

// BeginBucketTxn creates a new in-memory transacting key-value store that wraps the same store that b wraps.
// Auto-commit is turned off for the new store (allowing staged operations).
func (b *KVTxn) BeginBucketTxn(context.Context) (kv.BucketTxnCompleter, error) {
	return new(b.store, b.keyLock, false), nil
}
