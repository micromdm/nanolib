package kvprefix

import (
	"context"
)

// Keys returns all keys in the underlying key-value store.
// The returned keys have no ordering guaratees.
// The keys channel should be closed if cancel was provided and closed.
// Beware of deadlocks with underlying implementations.
func (b *KVPrefix) Keys(ctx context.Context, cancel <-chan struct{}) <-chan string {
	return b.KeysPrefix(ctx, "", cancel)
}

// Keys returns all keys starting with prefix in the underlying key-value store.
// The returned keys have no ordering guaratees.
// The keys channel should be closed if cancel was provided and closed.
// Beware of deadlocks with underlying implementations.
func (b *KVPrefix) KeysPrefix(ctx context.Context, prefix string, cancel <-chan struct{}) <-chan string {
	r := make(chan string)
	go func() {
		defer close(r)
		for k := range b.store.KeysPrefix(ctx, b.prefix+prefix, cancel) {
			select {
			case <-cancel:
				return
			case r <- k[len(b.prefix):]:
			}
		}
	}()
	return r
}
