package mojo_test

import (
	"bufio"
	"net/http"
	"net/url"
	"strings"
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

func TestRequestFormParams(t *testing.T) {
	raw := `POST /foo?bar=baz&fizz=no HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 9

fizz=buzz`
	rawReq, err := http.ReadRequest(bufio.NewReader(strings.NewReader(raw)))
	if err != nil {
		t.Fatalf("Could not read request: %v", err)
	}

	req := &mojo.Request{}
	req.Read(rawReq)

	if req.BodyParams.Param("fizz") != "buzz" {
		t.Errorf("POST body parameter not correct")
	}
	if req.QueryParams.Param("bar") != "baz" {
		t.Errorf("POST query parameter not correct")
	}
	if req.QueryParams.Param("fizz") != "no" {
		t.Errorf("POST query parameter not correct")
	}
	if req.Param("bar") != "baz" {
		t.Errorf("POST query parameter does not take precedence")
	}
	if req.Param("fizz") != "buzz" {
		t.Errorf("POST body parameter does not take precedence")
	}
}
