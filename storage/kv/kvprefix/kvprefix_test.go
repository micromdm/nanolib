package kvprefix

import (
	"context"
	"testing"

	"github.com/micromdm/nanolib/storage/kv/kvmap"
	"github.com/micromdm/nanolib/storage/kv/test"
)

func TestKVPrefix(t *testing.T) {
	ctx := context.Background()
	b := kvmap.New()
	prefixBucket1 := New("kvprefix1.", b)

	// run the standard kv tests
	test.TestBucketSimple(t, ctx, prefixBucket1)
	test.TestKeysTraversing(t, ctx, New("kvprefix2.", b))

	// set a value in our prefixed store
	err := prefixBucket1.Set(ctx, "lorem", []byte("ipsum"))
	if err != nil {
		t.Fatal(err)
	}

	// check that the value is represented in the underlying store.
	getVal, err := b.Get(ctx, "kvprefix1.lorem")
	if err != nil {
		t.Fatal(err)
	}
	if have, want := string(getVal), "ipsum"; have != want {
		t.Errorf("have = %q, want = %q", have, want)
	}
}
