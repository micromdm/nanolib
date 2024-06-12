// Package kvprefix provides a pseudo key-value store which uses a prefix
// string for keys over an existing store.
package kvprefix

import "github.com/micromdm/nanolib/storage/kv"

// KVPrefix is a a pseudo key-value store which uses a prefix
// string for all keys over an existing store.
type KVPrefix struct {
	prefix string
	store  kv.KeysPrefixTraversingBucket
}

// New creates a new prefix store.
func New(prefix string, b kv.KeysPrefixTraversingBucket) *KVPrefix {
	return &KVPrefix{prefix: prefix, store: b}
}
