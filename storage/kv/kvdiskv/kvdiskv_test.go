package kvdiskv

import (
	"context"
	"testing"

	"github.com/micromdm/nanolib/storage/kv/test"
	"github.com/peterbourgon/diskv/v3"
)

func newDV(t *testing.T) *diskv.Diskv {
	return diskv.New(diskv.Options{
		BasePath:     t.TempDir(),
		Transform:    FlatTransform,
		CacheSizeMax: 1024 * 1024,
	})

}

func TestKVMap(t *testing.T) {
	ctx := context.Background()
	test.TestBucketSimple(t, ctx, New(newDV(t)))
	test.TestKeysTraversing(t, ctx, New(newDV(t)))
}
