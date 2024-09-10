package kvtxn

import (
	"context"
	"testing"

	"github.com/micromdm/nanolib/storage/kv/kvmap"
	"github.com/micromdm/nanolib/storage/kv/test"
)

func TestNop(t *testing.T) {
	ctx := context.Background()
	b := NewNopTxn(kvmap.New())
	bt, err := b.BeginCRUDBucketTxn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	test.TestBucketSimple(t, ctx, b)
	test.TestKeysTraversing(t, ctx, b)

	// We cannot run `test.TestTxnSimple()` because NopTxn does not
	// actually support commit/rollback (i.e. caching transaction data)
	// which is specifically tested for.
	// test.TestTxnSimple(t, ctx, b)

	// test just to make sure no error
	err = bt.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// test just to make sure no error
	err = bt.Rollback(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
