package kvdiskv

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/micromdm/nanolib/storage/kv"
)

// Get retreives the value at key in the diskv store.
// If key is not found then a wrapped ErrKeyNotFound will be returned.
func (b *KVDiskv) Get(_ context.Context, key string) ([]byte, error) {
	r, err := b.diskv.Read(key)
	if errors.Is(err, os.ErrNotExist) {
		// replace error type to comply with interface
		return r, fmt.Errorf("%w: %s", kv.ErrKeyNotFound, key)
	}
	return r, err
}

// Set sets key to value in the diskv store.
func (b *KVDiskv) Set(_ context.Context, key string, value []byte) error {
	return b.diskv.Write(key, value)
}

// Has checks that key is found in the diskv store.
func (b *KVDiskv) Has(_ context.Context, key string) (bool, error) {
	return b.diskv.Has(key), nil
}

// Delete deletes key in the diskv store.
func (b *KVDiskv) Delete(_ context.Context, key string) error {
	err := b.diskv.Erase(key)
	if errors.Is(err, os.ErrNotExist) {
		// hide this specific error to comply with interface
		return nil
	}
	return err
}
