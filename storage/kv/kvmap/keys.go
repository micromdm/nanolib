package kvmap

import (
	"context"
	"strings"
)

// Keys returns all keys in the Go map.
// The returned keys have no ordering guaratees.
// The keys channel will be closed if cancel was provided and closed.
// Note that a goroutine which holds a read lock is spawned. This will
// deadlock any writes until the goroutine is done.
func (b *KVMap) Keys(ctx context.Context, cancel <-chan struct{}) <-chan string {
	return b.KeysPrefix(ctx, "", cancel)
}

// Keys returns all keys starting with prefix in the Go map.
// The returned keys have no ordering guaratees.
// The keys channel will be closed if cancel was provided and closed.
// Note that a goroutine which holds a read lock is spawned. This will
// deadlock any writes until the goroutine is done.
func (b *KVMap) KeysPrefix(_ context.Context, prefix string, cancel <-chan struct{}) <-chan string {
	r := make(chan string)
	go func() {
		b.mu.RLock()
		defer b.mu.RUnlock()
		defer close(r)
		for k := range b.m {
			if prefix != "" && !strings.HasPrefix(k, prefix) {
				continue
			}
			select {
			case <-cancel:
				return
			case r <- k:
			}
		}
	}()
	return r
}
