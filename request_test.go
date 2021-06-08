package mojo_test

import (
	"testing"

	"github.com/preaction/mojo.go"
	// XXX: Move test -> mojotest
	mojotest "github.com/preaction/mojo.go/test"
)

func TestRequestRead(t *testing.T) {
	raw := mojotest.BuildHTTPRequest(t, "GET /foo HTTP/1.1\n\n")

	req := &mojo.Request{}
	req.Read(raw)

	if req.Method != "GET" {
		t.Errorf("Method not read correctly. Expected: %s, Got: %s", "GET", req.Method)
	}
}

func TestRequestFormParams(t *testing.T) {
	raw := mojotest.BuildHTTPRequest(t, `POST /foo?bar=baz&fizz=no HTTP/1.1
Content-Type: application/x-www-form-urlencoded
Content-Length: 9

fizz=buzz`)

	req := &mojo.Request{}
	req.Read(raw)

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
