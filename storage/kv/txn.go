package kv

import "context"

// TxnCompleter completes transactions.
type TxnCompleter interface {
	// Commit permanently applies transaction changes.
	// Some implementations may not continue using the transaction.
	Commit(ctx context.Context) error

	// Rollback discards transaction changes.
	// Some implementations may not continue using the transaction.
	Rollback(ctx context.Context) error
}

// CRUDBucketTxnCompleter is a transacting key-value store.
type CRUDBucketTxnCompleter interface {
	CRUDBucket
	TxnCompleter
}

// CRUDBucketTxnBeginner can start transactions.
type CRUDBucketTxnBeginner interface {
	// BeginCRUDBucketTxn creates a new transaction that can later be completed.
	BeginCRUDBucketTxn(ctx context.Context) (CRUDBucketTxnCompleter, error)
}

// BucketTxnBucket is a key-value store that can start transactions.
// If a transaction has not begun then individual operations should auto-commit.
type TxnCRUDBucket interface {
	CRUDBucketTxnBeginner
	CRUDBucketTxnCompleter
}

// KeysPrefixTraversingBucketTxnCompleter is a transacting key-value store.
// This store can traverse keys including using a prefix.
type KeysPrefixTraversingBucketTxnCompleter interface {
	KeysPrefixTraversingBucket
	TxnCompleter
}

// KeysPrefixTraversingBucketTxnBeginner can start transactions.
type KeysPrefixTraversingBucketTxnBeginner interface {
	// BeginKeysPrefixTraversingBucketTxn creates a new transaction that can later be completed.
	BeginKeysPrefixTraversingBucketTxn(ctx context.Context) (KeysPrefixTraversingBucketTxnCompleter, error)
}

// TxnKeysPrefixTraversingBucket is a key-value store that can start transactions.
// If a transaction has not begun then individual operations should auto-commit.
type TxnKeysPrefixTraversingBucket interface {
	KeysPrefixTraversingBucketTxnBeginner
	KeysPrefixTraversingBucketTxnCompleter
}

// BucketTxnCompleter is a transacting key-value store.
type BucketTxnCompleter interface {
	Bucket
	TxnCompleter
}

// BucketTxnBeginner can start transactions.
type BucketTxnBeginner interface {
	// BeginBucketTxn creates a new transaction that can later be completed.
	BeginBucketTxn(ctx context.Context) (BucketTxnCompleter, error)
}

// TxnBucket is a key-value store that can start transactions.
// If a transaction has not begun then individual operations should auto-commit.
type TxnBucket interface {
	BucketTxnBeginner
	BucketTxnCompleter
}
