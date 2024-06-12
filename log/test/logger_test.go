package test

import "testing"

func TestLogger(t *testing.T) {
	logger := new(Logger)

	logger.Info("msg", "hello")

	TestLastLogKeyValueMatches(t, logger, "msg", "hello")
}
