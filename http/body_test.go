package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func ExampleGetAndReplaceBodyBytes() {
	buf := bytes.NewBuffer([]byte("Hello, world!"))

	req, err := http.NewRequestWithContext(context.Background(), "GET", "", buf)
	if err != nil {
		log.Fatal(err)
	}

	ret, err := GetAndReplaceBodyBytes(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(ret))
	// Output: Hello, world!
}

func TestCopyAndBody(t *testing.T) {
	testBytes := []byte("Hello, world!")
	buf := bytes.NewBuffer(testBytes)

	req, err := http.NewRequestWithContext(context.Background(), "GET", "", buf)
	if err != nil {
		t.Fatal(err)
	}

	ret, err := GetAndReplaceBodyBytes(req)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(testBytes, ret) {
		t.Error("not equal")
	}

	buf2 := &bytes.Buffer{}
	_, err = io.Copy(buf2, req.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(testBytes, buf2.Bytes()) {
		t.Error("not equal")
	}
}
