package kvtxn

import (
	"context"
	"testing"

	"github.com/micromdm/nanolib/storage/kv/kvmap"
	"github.com/micromdm/nanolib/storage/kv/test"
)

func TestKVTxn(t *testing.T) {
	b := New(kvmap.New())
	ctx := context.Background()
	test.TestBucketSimple(t, ctx, b)
	test.TestKeysTraversing(t, ctx, b)
	test.TestTxnSimple(t, ctx, b)
	b = New(kvmap.New()) // clear test data
	t.Run("TestKVTxnKeys", func(t *testing.T) { test.TestKVTxnKeys(t, ctx, b) })
}
