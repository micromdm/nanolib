// Package kvdiskv wraps diskv to a standard interface for a key-value store.
package kvdiskv

import (
	"github.com/peterbourgon/diskv/v3"
)

// FlatTransform is a diskv TransformFunction. From diskv; the
// Simplest transform function: put all the data files into the base dir.
var FlatTransform = func(s string) []string { return []string{} }

// KVDiskv wraps diskv to implement an on-disk key-value store.
type KVDiskv struct {
	diskv *diskv.Diskv
}

// New creates a new on-disk key-value store backed by dv.
func New(dv *diskv.Diskv) *KVDiskv {
	if dv == nil {
		panic("nil diskv")
	}
	return &KVDiskv{diskv: dv}
}
