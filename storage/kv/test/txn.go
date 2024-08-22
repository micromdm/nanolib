package test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/micromdm/nanolib/storage/kv"
)

func TestTxnSimple(t *testing.T, ctx context.Context, b kv.BucketTxnBucket) {
	// first, set a value in the "parent" bucket
	err := b.Set(ctx, "test-txn-key-1", []byte("test-txn-val-1"))
	if err != nil {
		t.Fatal(err)
	}

	// sanity check by reading the value we just set
	val, err := b.Get(ctx, "test-txn-key-1")
	if err != nil {
		t.Fatal(err)
	}
	if have, want := val, []byte("test-txn-val-1"); !bytes.Equal(have, want) {
		t.Errorf("have: %v, want: %v", string(have), string(want))
	}

	// create a txn
	bt, err := b.BeginBucketTxn(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// sanity check by reading the value we just set within the txn
	val, err = bt.Get(ctx, "test-txn-key-1")
	if err != nil {
		t.Fatal(err)
	}
	if have, want := val, []byte("test-txn-val-1"); !bytes.Equal(have, want) {
		t.Errorf("have: %v, want: %v", string(have), string(want))
	}

	// now, reset the key within the txn ...
	err = bt.Set(ctx, "test-txn-key-1", []byte("test-txn-val-2"))
	if err != nil {
		t.Fatal(err)
	}

	// ... but rollback the transaction.
	err = bt.Rollback(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// read the value we just reset in the parent and make sure it hasn't changed
	val, err = b.Get(ctx, "test-txn-key-1")
	if err != nil {
		t.Fatal(err)
	}
	if have, want := val, []byte("test-txn-val-1"); !bytes.Equal(have, want) {
		t.Errorf("have: %v, want: %v", string(have), string(want))
	}

	// read the value we just reset in the txn and make sure it was rolled back
	val, err = bt.Get(ctx, "test-txn-key-1")
	if err != nil {
		t.Fatal(err)
	}
	if have, want := val, []byte("test-txn-val-1"); !bytes.Equal(have, want) {
		t.Errorf("have: %v, want: %v", string(have), string(want))
	}

	// okay, let's re-reset the value again
	err = bt.Set(ctx, "test-txn-key-1", []byte("test-txn-val-2"))
	if err != nil {
		t.Fatal(err)
	}

	// now, commit the change
	err = bt.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// and make sure the "parent" bucket received that changed
	val, err = b.Get(ctx, "test-txn-key-1")
	if err != nil {
		t.Fatal(err)
	}

	if have, want := val, []byte("test-txn-val-2"); !bytes.Equal(have, want) {
		t.Errorf("have: %v, want: %v", string(have), string(want))
	}

	// lets make a new txn
	bt, err = b.BeginBucketTxn(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// set a value
	err = bt.Set(ctx, "test-txn-key-2", []byte("test-txn-val-3"))
	if err != nil {
		t.Fatal(err)
	}

	// sanity check by reading the value we just set in the within the txn
	val, err = bt.Get(ctx, "test-txn-key-2")
	if err != nil {
		t.Fatal(err)
	}
	if have, want := val, []byte("test-txn-val-3"); !bytes.Equal(have, want) {
		t.Errorf("have: %v, want: %v", string(have), string(want))
	}

	// now, rollback our changes:
	err = bt.Rollback(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// and try and read the values we just set (but discarded)
	// should error with a key not found
	_, err = bt.Get(ctx, "test-txn-key-2")
	if !errors.Is(err, kv.ErrKeyNotFound) {
		t.Fatal(err)
	}

	// .. same for the parent bucket
	_, err = b.Get(ctx, "test-txn-key-2")
	if !errors.Is(err, kv.ErrKeyNotFound) {
		t.Fatal(err)
	}
}
