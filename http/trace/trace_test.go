package trace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/micromdm/nanolib/log/test"
)

type contextCapture struct {
	body []byte
	ctx  context.Context
}

func (c *contextCapture) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.ctx = r.Context()
	w.Write(c.body)
}

func TestHTTPCtxLog(t *testing.T) {
	ctx := context.Background()

	r, err := http.NewRequestWithContext(ctx, "GET", "/x-test-path", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	logger := &test.Logger{KeepLastWith: true}

	next := &contextCapture{body: []byte("x-test-body")}

	handler := NewTraceLoggingHandler(next, logger, func(_ *http.Request) string { return "x-test-trace" })

	handler.ServeHTTP(w, r)

	if want, have := http.StatusOK, w.Code; want != have {
		t.Errorf("http status: want: %q, have: %q", want, have)
	}

	if want, have := "x-test-body", w.Body.String(); want != have {
		t.Errorf("body bytes: want: %q, have: %q", want, have)
	}

	// verify the trace ID is on the captured request context
	if want, have := "x-test-trace", GetTraceID(next.ctx); want != have {
		t.Errorf("body bytes: want: %q, have: %q", want, have)
	}

	// also check the trace ID made it through to the log line
	logVal, _ := logger.LastKeyStringValue("trace_id")
	if want, have := "x-test-trace", logVal; want != have {
		t.Errorf("body bytes: want: %q, have: %q", want, have)
	}

	logVal, _ = logger.LastKeyStringValue("method")
	if want, have := "GET", logVal; want != have {
		t.Errorf("body bytes: want: %q, have: %q", want, have)
	}

	logVal, _ = logger.LastKeyStringValue("path")
	if want, have := "/x-test-path", logVal; want != have {
		t.Errorf("body bytes: want: %q, have: %q", want, have)
	}
}
