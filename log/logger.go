package log

// Pacakge log is embedded (not imported) from:
// https://github.com/jessepeterson/go-log

// Logger is a generic logging interface for a structured, levelled, nest-able logger.
type Logger interface {
	// Info logs using the info level.
	Info(...interface{})

	// Debug logs using the debug level.
	Debug(...interface{})

	// With returns a new nested Logger.
	// Usually for adding logging additional context to an existing logger.
	With(...interface{}) Logger
}
