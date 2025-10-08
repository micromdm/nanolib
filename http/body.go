package http

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// GetAndReplaceBodyBytes returns the body of req as a byte slice.
// If needed the body is replaced with a byte buffer for re-use.
func GetAndReplaceBodyBytes(req *http.Request) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := GetAndReplaceBody(req, buf)
	return buf.Bytes(), err
}

// GetAndReplaceBody writes the body of req to w.
// If needed the body is replaced with a byte buffer for re-use.
func GetAndReplaceBody(req *http.Request, w io.Writer) error {
	if req == nil {
		return errors.New("nil request")
	}
	if req.Body == nil || w == nil {
		return nil
	}

	var err error

	body := req.Body
	var buf *bytes.Buffer

	if req.GetBody != nil {
		// GetBody returns a copy of a the body for reading
		// using it implies we don't need to replace the body
		body, err = req.GetBody()
		if err != nil {
			return fmt.Errorf("getting body: %w", err)
		} else if body == nil {
			return nil
		}
	} else {
		buf = new(bytes.Buffer)
		w = io.MultiWriter(buf, w)
	}

	_, err = io.Copy(w, body)
	if err != nil {
		return fmt.Errorf("copying body: %w", err)
	}
	defer body.Close()

	if buf != nil {
		req.Body = io.NopCloser(buf)
	}

	return nil
}
