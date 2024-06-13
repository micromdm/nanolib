package kvprefix

import (
	"context"
)

// Get retrieves the value at key in the underlying store.
// The key is preprended with the prefix.
func (b *KVPrefix) Get(ctx context.Context, key string) ([]byte, error) {
	return b.store.Get(ctx, b.prefix+key)
}

// Set sets key to value in the underlying store.
// The key is preprended with the prefix.
func (b *KVPrefix) Set(ctx context.Context, key string, value []byte) error {
	return b.store.Set(ctx, b.prefix+key, value)
}

// Has checks that key is found the underlying store.
// The key is preprended with the prefix.
func (b *KVPrefix) Has(ctx context.Context, key string) (bool, error) {
	return b.store.Has(ctx, b.prefix+key)
}

// Delete deletes key the underlying store.
// The key is preprended with the prefix.
func (b *KVPrefix) Delete(ctx context.Context, key string) error {
	return b.store.Delete(ctx, b.prefix+key)
}
