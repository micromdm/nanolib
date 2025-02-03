package test

import (
	"sync"

	"github.com/micromdm/nanolib/log"
)

// Log is a log line with all logged context and log level.
type Log struct {
	Debug bool
	Log   []interface{}
}

// Logger is a logger that accumulates log lines.
type Logger struct {
	mu       sync.RWMutex
	ctx      []interface{}
	lastLogs *[]Log

	// KeepLastWith turns on keeping and sharing of the last log lines
	// for loggers in the With() chain.
	KeepLastWith bool
}

// Info logs using the info level appending to the test logger's stored logs.
func (l *Logger) Info(args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.lastLogs == nil {
		l.lastLogs = &[]Log{}
	}
	*l.lastLogs = append(*l.lastLogs, Log{
		Debug: false,
		Log:   append(l.ctx, args...),
	})
}

// Debug logs using the debug level level appending to the test logger's stored logs.
func (l *Logger) Debug(args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.lastLogs == nil {
		l.lastLogs = &[]Log{}
	}
	*l.lastLogs = append(*l.lastLogs, Log{
		Debug: true,
		Log:   append(l.ctx, args...),
	})
}

// With returns a new nested Logger.
func (l *Logger) With(args ...interface{}) log.Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.lastLogs == nil {
		l.lastLogs = &[]Log{}
	}
	newLogger := &Logger{ctx: append(l.ctx, args...)}
	if l.KeepLastWith {
		// note we copy the pointer of the lastLogs here to maintain the
		// same log lines for any logger in this With() chain.
		newLogger.lastLogs = l.lastLogs
	}
	return newLogger
}

// Last returns the last logged line.
func (l *Logger) Last() *Log {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.lastLogs == nil {
		l.lastLogs = &[]Log{}
	}
	if len(*l.lastLogs) < 1 {
		return nil
	}
	return &(*l.lastLogs)[len(*l.lastLogs)-1]
}

// LastKey returns the value at the string key in the last log line.
// We make the assumption logs are key-value tuples using string keys.
func (l *Logger) LastKey(key string) (last *Log, value interface{}, found bool) {
	last = l.Last()
	if last == nil {
		return
	}

	// at least two values required (for key and value)
	if len(last.Log) < 2 {
		return
	}

	// find our first matching key and return its value
	for i := 0; i < len(last.Log); i = i + 2 {
		logKey, ok := last.Log[i].(string)
		if !ok {
			continue
		}

		if logKey == key {
			found = true
			value = last.Log[i+1]
			break
		}
	}

	return
}

// LastKeyStringValue is like LastKey but assumes you only want a string value.
func (l *Logger) LastKeyStringValue(key string) (value string, found bool) {
	var v interface{}
	_, v, found = l.LastKey(key)
	value, _ = v.(string)
	return
}
