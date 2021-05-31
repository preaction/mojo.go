package mojo_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/preaction/mojo.go"
)

func TestRequestRead(t *testing.T) {
	url, err := url.ParseRequestURI("/foo")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}

	raw := http.Request{Method: "GET", URL: url}

	req := &mojo.Request{}
	req.Read(&raw)

	if req.Method != "GET" {
		t.Errorf("Method not read correctly. Expected: %s, Got: %s", "GET", req.Method)
	}
}
