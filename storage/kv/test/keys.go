package test

import (
	"bytes"
	"context"
	"testing"

	"github.com/micromdm/nanolib/storage/kv"
)

// TestKeysTraversing tests retrieving keys from stores.
// Note because we're enumating (all) keys in a store and testing any
// remainders b should not have any keys already set.
func TestKeysTraversing(t *testing.T, ctx context.Context, b kv.KeysPrefixTraversingBucket) {
	kvMap := map[string][]byte{
		"hello": []byte("world"),
		"foo":   []byte("bar"),
		"help":  []byte("i need somebody"),
	}

	// put the data into b
	err := kv.SetMap(ctx, b, kvMap)
	if err != nil {
		t.Fatal(err)
	}

	kvMap2 := copyMap(kvMap)

	// iterate over keys and remove entries from the copied map
	for k := range b.Keys(ctx, nil) {
		if _, ok := kvMap2[k]; ok {
			delete(kvMap2, k)
		} else {
			// test shuold be started with a new pristine bucket
			// so we should never get this
			t.Errorf("unexpected key returned from bucket: %s", k)
		}
	}

	if len(kvMap2) > 0 {
		t.Error("leftover keys (not all keys returned)")
	}

	kvMap3 := copyMap(kvMap)

	// iterate over prefix keys and remove entries from the copied map
	for k := range b.KeysPrefix(ctx, "hel", nil) {
		if _, ok := kvMap3[k]; ok {
			delete(kvMap3, k)
		} else {
			// test shuold be started with a new pristine bucket
			// so we should never get this
			t.Errorf("unexpected key returned from bucket: %s", k)
		}
	}

	if value, ok := kvMap3["foo"]; len(kvMap3) != 1 || !ok || !bytes.Equal(value, kvMap["foo"]) {
		t.Error("incorrect leftover value(s) (should just be foo key)")
	}
}

func copyMap(in map[string][]byte) (out map[string][]byte) {
	out = make(map[string][]byte, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
