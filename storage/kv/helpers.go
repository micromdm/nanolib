package kv

import (
	"context"
	"fmt"
)

// SetMap iterates over m to set the keys in b and returns any error immediately.
func SetMap(ctx context.Context, b RWBucket, m map[string][]byte) error {
	var err error
	for k, v := range m {
		if err = b.Set(ctx, k, v); err != nil {
			return fmt.Errorf("setting %s: %w", k, err)
		}
	}
	return nil
}

// GetMap iterates over keys to get the values in b and returns any error immediately.
func GetMap(ctx context.Context, b ROBucket, keys []string) (map[string][]byte, error) {
	var err error
	ret := make(map[string][]byte)
	for _, k := range keys {
		if ret[k], err = b.Get(ctx, k); err != nil {
			return ret, fmt.Errorf("getting %s: %w", k, err)
		}
	}
	return ret, nil
}

// DeleteSlice deletes s keys from b and returns any error immediately.
func DeleteSlice(ctx context.Context, b RWBucket, s []string) error {
	var err error
	for _, i := range s {
		if err = b.Delete(ctx, i); err != nil {
			return fmt.Errorf("deleting %s: %w", i, err)
		}
	}
	return nil
}

// AllKeys collects and returns all keys in b in a slice.
// Warning: this buffers all keys in b. For large stores this may be prohibitive.
func AllKeys(ctx context.Context, b KeysTraverser) (r []string) {
	for k := range b.Keys(ctx, nil) {
		r = append(r, k)
	}
	return
}

// AllKeysPrefix collects and returns a slice of all keys starting with prefix in b.
// Warning: this buffers the found keys in b. For large stores this may be prohibitive.
func AllKeysPrefix(ctx context.Context, b KeysPrefixTraverser, prefix string) (r []string) {
	for k := range b.KeysPrefix(ctx, prefix, nil) {
		r = append(r, k)
	}
	return
}
