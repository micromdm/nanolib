package kvtxn

import (
	"context"

	"github.com/micromdm/nanolib/storage/kv"
)

// NopTxn wraps KV bucket storage to support transaction calls.
// However it does not support transactions: KV storage calls are
// passed directly through to the wrapped storage.
// WARNING: this is dangerous.
type NopTxn struct {
	kv.KeysPrefixTraversingBucket
}

// NewNopTxn creates a new pass-through transaction wrapper.
func NewNopTxn(b kv.KeysPrefixTraversingBucket) *NopTxn {
	return &NopTxn{KeysPrefixTraversingBucket: b}
}

// BeginBucketTxn simply returns b.
func (b *NopTxn) BeginBucketTxn(context.Context) (kv.BucketTxnCompleter, error) {
	return b, nil
}

// BeginKeysPrefixTraversingBucketTxn simply returns b.
func (b *NopTxn) BeginKeysPrefixTraversingBucketTxn(context.Context) (kv.KeysPrefixTraversingBucketTxnCompleter, error) {
	return b, nil
}

// Commit does nothing.
func (b *NopTxn) Commit(context.Context) error {
	return nil
}

// Rollback does nothing.
func (b *NopTxn) Rollback(context.Context) error {
	return nil
}
