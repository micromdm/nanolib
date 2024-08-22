package kv

import "context"

// TxnCompleter completes transactions; readying them for the next transaction.
type TxnCompleter interface {
	// Commit permanently applies the transaction's changes.
	Commit(ctx context.Context) error

	// Rollback discards the transaction's changes.
	Rollback(ctx context.Context) error
}

// BucketTxnCompleter is a transacting key-value store.
type BucketTxnCompleter interface {
	Bucket
	TxnCompleter
}

// BucketTxnBeginner can start transactions.
type BucketTxnBeginner interface {
	// Begin creates a new transaction that can later be completed.
	BeginBucketTxn(ctx context.Context) (BucketTxnCompleter, error)
}

// BucketTxnBucket is a key-value store that can start transactions.
// If a transaction has not begun then individual operations should auto-commit.
type BucketTxnBucket interface {
	BucketTxnBeginner
	BucketTxnCompleter
}

// KeysPrefixTraversingBucketTxnCompleter is a transacting key-value store.
// This store can traverse keys including using a prefix.
type KeysPrefixTraversingBucketTxnCompleter interface {
	KeysPrefixTraversingBucket
	TxnCompleter
}

// KeysPrefixTraversingBucketTxnBeginner can start transactions.
type KeysPrefixTraversingBucketTxnBeginner interface {
	// Begin creates a new transaction that can later be completed.
	BeginKeysPrefixTraversingBucketTxn(ctx context.Context) (KeysPrefixTraversingBucketTxnCompleter, error)
}

// KeysPrefixTraversingBucketTxnBucket is a key-value store that can start transactions.
// If a transaction has not begun then individual operations should auto-commit.
type KeysPrefixTraversingBucketTxnBucket interface {
	KeysPrefixTraversingBucketTxnBeginner
	KeysPrefixTraversingBucketTxnCompleter
}
