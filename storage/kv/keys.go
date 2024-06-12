package kv

import "context"

// KeysTraverser can traverse keys.
type KeysTraverser interface {
	// Keys returns all keys in the key-value store.
	// The returned keys have no ordering guaratees.
	// The keys channel should be closed if cancel was provided and closed.
	// Beware of deadlocks with underlying implementations.
	Keys(ctx context.Context, cancel <-chan struct{}) <-chan string
}

// KeysPrefixTraverser can traverse keys using a prefix.
type KeysPrefixTraverser interface {
	// Keys returns all keys starting with prefix in the key-value store.
	// The returned keys have no ordering guaratees.
	// The keys channel should be closed if cancel was provided and closed.
	// Beware of deadlocks with underlying implementations.
	KeysPrefix(ctx context.Context, prefix string, cancel <-chan struct{}) <-chan string
}

// KeysTraversingBucket is a key-value store that can traverse keys.
type KeysTraversingBucket interface {
	Bucket
	KeysTraverser
}

// KeysPrefixTraversingBucket is a key-value store that can traverse keys.
// Inlcuding using a prefix.
type KeysPrefixTraversingBucket interface {
	Bucket
	KeysTraverser
	KeysPrefixTraverser
}
