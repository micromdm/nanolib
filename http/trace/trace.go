// Pacakge trace implements a trace logging HTTP request middleware handler.
package trace

import (
	"context"
	"net"
	"net/http"

	"github.com/micromdm/nanolib/log"
	"github.com/micromdm/nanolib/log/ctxlog"
)

type ctxKeyTraceID struct{}

// GetTraceID returns the trace ID from ctx.
func GetTraceID(ctx context.Context) string {
	id, _ := ctx.Value(ctxKeyTraceID{}).(string)
	return id
}

// NewTraceLoggingHandler logs HTTP request details and sets up trace logging.
// The host, method, URL path, user agent are immediately logged.
// The "x_forwarded_for" header is also logged if present.
// The "trace_id" key will be set on the context logger and logged if
// traceIDFn is set.
func NewTraceLoggingHandler(next http.Handler, logger log.Logger, traceIDFn func(*http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if traceIDFn != nil {
			ctx = context.WithValue(ctx, ctxKeyTraceID{}, traceIDFn(r))
			ctx = ctxlog.AddFunc(ctx, ctxlog.SimpleStringFunc("trace_id", ctxKeyTraceID{}))
		}

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		logs := []interface{}{
			"addr", host,
			"method", r.Method,
			"path", r.URL.Path,
			"agent", r.UserAgent(),
		}

		if fwdedFor := r.Header.Get("X-Forwarded-For"); fwdedFor != "" {
			logs = append(logs, "x_forwarded_for", fwdedFor)
		}

		ctxlog.Logger(ctx, logger).Info(logs...)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
