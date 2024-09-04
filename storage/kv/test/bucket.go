package test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/micromdm/nanolib/storage/kv"
)

func TestBucketSimple(t *testing.T, ctx context.Context, b kv.CRUDBucket) {
	const testKey1 = "test_key_1"
	const testKey2 = "test_key_2"
	const testValue1 = "test_value_1"

	for i := 0; i < 5; i++ {
		si := "_iter_" + strconv.Itoa(i)

		// first, delete the key to make sure its gone.
		err := b.Delete(ctx, testKey1+si)
		if err != nil {
			t.Fatal(err)
		}

		// then, do it again to try and catch implementations that
		// throw errors for non-existent keys
		err = b.Delete(ctx, testKey1+si)
		if err != nil {
			t.Fatalf("should not throw error for non-existing key: %v", err)
		}

		// then try and read the missing key, to get the explicit error
		_, err = b.Get(ctx, testKey1+si)
		if err == nil {
			t.Error("should be an error")
		}
		if !errors.Is(err, kv.ErrKeyNotFound) {
			t.Errorf("error should be an ErrKeyNotFound, but found: %v", err)
		}

		// basic write test
		err = b.Set(ctx, testKey2+si, []byte(testValue1))
		if err != nil {
			t.Fatal(err)
		}

		// basic check
		found, err := b.Has(ctx, testKey2+si)
		if err != nil {
			t.Fatal(err)
		}
		if !found {
			t.Errorf("key not found: %s", testKey2+si)
		}

		// basic read test
		v, err := b.Get(ctx, testKey2+si)
		if err != nil {
			t.Fatal(err)
		}
		if have, want := string(v), testValue1; have != want {
			t.Errorf("have = %q, want = %q", have, want)
		}

		// cleanup
		err = b.Delete(ctx, testKey2+si)
		if err != nil {
			t.Fatal(err)
		}
	}
}
