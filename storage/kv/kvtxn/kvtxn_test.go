package kvtxn

import (
	"context"
	"testing"

	"github.com/micromdm/nanolib/storage/kv"
	"github.com/micromdm/nanolib/storage/kv/kvmap"
	"github.com/micromdm/nanolib/storage/kv/test"
)

func TestKVTxn(t *testing.T) {
	b := New(kvmap.New())
	ctx := context.Background()
	test.TestBucketSimple(t, ctx, b)
	test.TestKeysTraversing(t, ctx, b)
	test.TestTxnSimple(t, ctx, b)
}

func slicesEqual[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestKVTxnKeys(t *testing.T) {
	u := kvmap.New()
	b := New(u)
	ctx := context.Background()
	bt, err := b.BeginKeysPrefixTraversingBucketTxn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = bt.Set(ctx, "hello", []byte("world"))
	if err != nil {
		t.Fatal(err)
	}
	// make sure we have what we set in the txn
	keys := kv.AllKeys(ctx, bt)
	if want, have := []string{"hello"}, keys; !slicesEqual(want, have) {
		t.Errorf("want: %v, have: %v", want, have)
	}
	// delete the key
	err = bt.Delete(ctx, "hello")
	if err != nil {
		t.Fatal(err)
	}
	// check that we don't see it
	keys = kv.AllKeys(ctx, bt)
	if want, have := []string{}, keys; !slicesEqual(want, have) {
		t.Errorf("want: %v, have: %v", want, have)
	}
	// set a value on the non-txn store (auto-commit)
	err = u.Set(ctx, "hello", []byte("dlrow"))
	if err != nil {
		t.Fatal(err)
	}
	// again check that our txn does not see it
	keys = kv.AllKeys(ctx, bt)
	if want, have := []string{}, keys; !slicesEqual(want, have) {
		t.Errorf("want: %v, have: %v", want, have)
	}
}
