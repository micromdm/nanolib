package kv

import (
	"context"
	"errors"
)

var ErrKeyNotFound = errors.New("key not found")

// ROBucket defines simple read-only operations for key-value stores.
type ROBucket interface {
	// Has checks that key can be found.
	Has(ctx context.Context, key string) (found bool, err error)

	// Get retrieves value at key.
	// If key is not found then ErrKeyNotFound should be
	// returned in the error chain.
	Get(ctx context.Context, key string) (value []byte, err error)
}

// RWBucket defines simple write operations for key-value stores.
type RWBucket interface {
	// Set sets key to value.
	Set(ctx context.Context, key string, value []byte) error

	// Delete deletes key.
	// An error should not be returned if key does not exist.
	Delete(ctx context.Context, key string) error
}

// Bucket defines simple read-write operations for key-value stores.
type Bucket interface {
	ROBucket
	RWBucket
}
