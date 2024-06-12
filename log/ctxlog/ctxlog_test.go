package ctxlog

import (
	"context"
	"testing"

	"github.com/micromdm/nanolib/log"
	"github.com/micromdm/nanolib/log/test"
)

func TestCtxLog(t *testing.T) {
	ctx := context.Background()
	type ctxKeyTest1 struct{}
	const test1Value = "test1Value"

	// create a new logger
	var logger log.Logger = new(test.Logger)

	// add some initial context
	logger = logger.With("context1Key", "context1Value")

	// associate a (simple) KV func with our context key
	ctx = AddFunc(ctx, SimpleStringFunc("test1Key", ctxKeyTest1{}))

	// stick a test value on the context with the KV func context key
	ctx = context.WithValue(ctx, ctxKeyTest1{}, test1Value)

	// create a logger with context.. er, context
	ctxLogger := Logger(ctx, logger)

	// log a line
	ctxLogger.Info("log1Key", "log1Value")

	// cast to get back to our test logger (with useful methods)
	actuallyTestLogger := ctxLogger.(*test.Logger)

	// verify the original context
	test.TestLastLogKeyValueMatches(t, actuallyTestLogger, "context1Key", "context1Value")

	// check the contextually-logged log KV
	test.TestLastLogKeyValueMatches(t, actuallyTestLogger, "test1Key", "test1Value")

	// and of course the actual log line args
	test.TestLastLogKeyValueMatches(t, actuallyTestLogger, "log1Key", "log1Value")
}
