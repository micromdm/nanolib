package kvdiskv

import "context"

// Keys returns all keys in the diskv store.
// The returned keys have no ordering guaratees.
// The keys channel should be closed if cancel was provided and closed.
func (b *KVDiskv) Keys(_ context.Context, cancel <-chan struct{}) <-chan string {
	return b.diskv.Keys(cancel)
}

// Keys returns all keys starting with prefix in the diskv store.
// The returned keys have no ordering guaratees.
// The keys channel should be closed if cancel was provided and closed.
func (b *KVDiskv) KeysPrefix(_ context.Context, prefix string, cancel <-chan struct{}) <-chan string {
	return b.diskv.KeysPrefix(prefix, cancel)
}
