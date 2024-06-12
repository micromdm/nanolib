// Package ctxlog allows logging of key-value pairs stored with a context.
//
// At a high level for a package that wants to log data using a [context]:
//
//  1. Create a [CtxKVFunc] function that can turn a context into
//     a slice of interfaces. Typically this function will read a context
//     key and output a key-value pair for log lines. A helper for these
//     functions is [SimpleStringFunc].
//  2. Attach this function to a context using [AddFunc].
//  3. At some point later attach a value to a context using
//     [context.WithValue].
//  4. Then, when ready to log, create a "context aware logger" using
//     [Logger] which will execute the CtxKVFuncs and return a logger
//     With the key-values from the executed functions.
//
// In this way we've decoupled the specific loggers, data types, keys,
// context keys, and other data from themselves but still have convenient
// access to log data from a context. This also allows flexible layering
// of logging locations and middleware that a context may pass through.
package ctxlog

import (
	"context"
	"sync"

	"github.com/micromdm/nanolib/log"
)

// CtxKVFunc creates logger key-value pairs from a context.
// CtxKVFuncs should aim to be be as efficient as possible—ideally only
// doing the minimum to read context values and generate KV pairs. Each
// associated CtxKVFunc is called every time we adapt a logger with
// Logger.
type CtxKVFunc func(context.Context) []interface{}

// ctxKeyFuncs is the context key for storing and retriveing
// a funcs{} struct on a context.
type ctxKeyFuncs struct{}

// funcs holds the associated CtxKVFunc functions to run.
type funcs struct {
	sync.RWMutex
	funcs []CtxKVFunc
}

// AddFunc appends f to to ctx.
func AddFunc(ctx context.Context, f CtxKVFunc) context.Context {
	if ctx == nil {
		return ctx
	}
	ctxFuncs, ok := ctx.Value(ctxKeyFuncs{}).(*funcs)
	if !ok || ctxFuncs == nil {
		ctxFuncs = &funcs{}
	}
	ctxFuncs.Lock()
	ctxFuncs.funcs = append(ctxFuncs.funcs, f)
	ctxFuncs.Unlock()
	return context.WithValue(ctx, ctxKeyFuncs{}, ctxFuncs)
}

// Logger executes the associated KV funcs on ctx and returns a new Logger with the results.
// If ctx is nil or if there are no KV funcs associated then logger is returned as-is.
func Logger(ctx context.Context, logger log.Logger) log.Logger {
	if ctx == nil {
		return logger
	}
	ctxFuncs, ok := ctx.Value(ctxKeyFuncs{}).(*funcs)
	if !ok || ctxFuncs == nil {
		return logger
	}
	var acc []interface{}
	ctxFuncs.RLock()
	for _, f := range ctxFuncs.funcs {
		acc = append(acc, f(ctx)...)
	}
	ctxFuncs.RUnlock()
	return logger.With(acc...)
}

// SimpleStringFunc is a helper that makes a simple CtxKVFunc that
// returns a key-value pair if found on the context.
func SimpleStringFunc(logKey string, ctxKey interface{}) CtxKVFunc {
	return func(ctx context.Context) (out []interface{}) {
		v, _ := ctx.Value(ctxKey).(string)
		if v != "" {
			out = []interface{}{logKey, v}
		}
		return
	}
}
