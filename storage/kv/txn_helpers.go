package kv

import (
	"context"
	"fmt"
)

// CRUDBucketTxnPerformer is a function that executes key-value operations within a transaction.
type CRUDBucketTxnPerformer func(context.Context, CRUDBucket) error

// PerformCRUDBucketTxn calls f to execute KV operations within a transaction.
// It takes care of beginning a transaction, committing it, or rolling
// it back if f returns an error.
func PerformCRUDBucketTxn(ctx context.Context, beginner CRUDBucketTxnBeginner, f CRUDBucketTxnPerformer) error {
	// note: implementation same/similar to PerformKeysPrefixTraversingBucketTxn and PerformBucketTxn
	b, err := beginner.BeginCRUDBucketTxn(ctx)
	if err != nil {
		return fmt.Errorf("txn begin: %w", err)
	}
	if err = f(ctx, b); err != nil {
		if rbErr := b.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("txn rollback: %w; while trying to handle error: %v", rbErr, err)
		}
		return fmt.Errorf("txn rolled back: %w", err)
	}
	if err = b.Commit(ctx); err != nil {
		return fmt.Errorf("txn commit: %w", err)
	}
	return nil
}

// KeysPrefixTraversingBucketTxnPerformer is a function that executes KV operations within a transaction.
type KeysPrefixTraversingBucketTxnPerformer func(ctx context.Context, txn KeysPrefixTraversingBucket) error

// PerformKeysPrefixTraversingBucketTxn calls f to execute KV operations within a transaction.
// It takes care of beginning a transaction, committing it, or rolling
// it back if f returns an error.
func PerformKeysPrefixTraversingBucketTxn(ctx context.Context, beginner KeysPrefixTraversingBucketTxnBeginner, f KeysPrefixTraversingBucketTxnPerformer) error {
	// note: implementation same/similar to PerformCRUDBucketTxn and PerformBucketTxn
	b, err := beginner.BeginKeysPrefixTraversingBucketTxn(ctx)
	if err != nil {
		return fmt.Errorf("txn begin: %w", err)
	}
	if err = f(ctx, b); err != nil {
		if rbErr := b.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("txn rollback: %w; while trying to handle error: %v", rbErr, err)
		}
		return fmt.Errorf("txn rolled back: %w", err)
	}
	if err = b.Commit(ctx); err != nil {
		return fmt.Errorf("txn commit: %w", err)
	}
	return nil
}

// BucketTxnPerformer is a function that executes KV operations within a transaction.
type BucketTxnPerformer func(ctx context.Context, txn Bucket) error

// PerformBucketTxn calls f to execute KV operations within a transaction.
// It takes care of beginning a transaction, committing it, or rolling
// it back if f returns an error.
func PerformBucketTxn(ctx context.Context, beginner BucketTxnBeginner, f BucketTxnPerformer) error {
	// note: implementation same/similar to PerformCRUDBucketTxn and PerformKeysPrefixTraversingBucketTxn
	b, err := beginner.BeginBucketTxn(ctx)
	if err != nil {
		return fmt.Errorf("txn begin: %w", err)
	}
	if err = f(ctx, b); err != nil {
		if rbErr := b.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("txn rollback: %w; while trying to handle error: %v", rbErr, err)
		}
		return fmt.Errorf("txn rolled back: %w", err)
	}
	if err = b.Commit(ctx); err != nil {
		return fmt.Errorf("txn commit: %w", err)
	}
	return nil
}
