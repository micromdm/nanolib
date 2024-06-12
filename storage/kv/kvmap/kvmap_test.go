package kvmap

import (
	"context"
	"testing"

	"github.com/micromdm/nanolib/storage/kv/test"
)

func TestKVMap(t *testing.T) {
	ctx := context.Background()
	test.TestBucketSimple(t, ctx, New())
	test.TestKeysTraversing(t, ctx, New())
}
