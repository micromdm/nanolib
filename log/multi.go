package log

// Pacakge log is embedded (not imported) from:
// https://github.com/jessepeterson/go-log

// MultiLogger logs to multiple Loggers.
type MultiLogger struct {
	loggers []Logger
}

// NewMultiLogger creates a new MultiLogger with Loggers
func NewMultiLogger(loggers ...Logger) *MultiLogger {
	if len(loggers) < 1 {
		panic("must have at least one logger")
	}
	return &MultiLogger{loggers: loggers}
}

// Info logs the info level to each logger.
func (ml *MultiLogger) Info(keyvals ...interface{}) {
	for _, logger := range ml.loggers {
		logger.Info(keyvals...)
	}
}

// Info logs the debug level to each logger.
func (ml *MultiLogger) Debug(keyvals ...interface{}) {
	for _, logger := range ml.loggers {
		logger.Debug(keyvals...)
	}
}

// With creates a new MultiLogger using passing context to each logger.
func (ml *MultiLogger) With(keyvals ...interface{}) Logger {
	nml := new(MultiLogger)
	nml.loggers = make([]Logger, len(ml.loggers), cap(ml.loggers))
	for i, logger := range ml.loggers {
		nml.loggers[i] = logger.With(keyvals...)
	}
	return nml
}
