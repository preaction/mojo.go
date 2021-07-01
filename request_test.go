package mojo_test

import (
	"net/url"
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
		t.Errorf("Method not read correctly. Expected: %#v, Got: %#v", "GET", req.Method)
	}

	expectURL, _ := url.ParseRequestURI("/foo")
	if *req.URL != *expectURL {
		t.Errorf("URL not read correctly. Expected: %#v, Got: %#v", expectURL, req.URL)
	}
}

func TestRequestReadHeaders(t *testing.T) {
	raw := mojotest.BuildHTTPRequest(t, `GET /foo HTTP/1.1
Content-Length: 0
Host: example.com

`)

	req := &mojo.Request{}
	req.Read(raw)

	// Testing direct access to the map tests that case is preserved in
	// the map
	if len(req.Headers["Host"]) == 0 || req.Headers["Host"][0] != "example.com" {
		t.Errorf(`Host header not read correctly. Expected: []string{"example.com"}, Got: %v`, req.Headers["Host"])
	}
	if len(req.Headers["Content-Length"]) == 0 || req.Headers["Content-Length"][0] != "0" {
		t.Errorf(`Content-Length header not read correctly. Expected: []string{"0"}, Got: %v`, req.Headers["Content-Length"])
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

func TestRequestJSON(t *testing.T) {
	raw := mojotest.BuildHTTPRequest(t, `POST /foo?bar=baz&fizz=no HTTP/1.1
Content-Type: application/json
Content-Length: 39

{"foo":"foo","baz":"baz","fizz":"fizz"}`)

	req := &mojo.Request{}
	req.Read(raw)

	type TestJSON struct {
		Foo  string `json:"foo"`
		Baz  string `json:"baz"`
		Fizz string `json:"fizz"`
	}

	body := TestJSON{}
	err := req.JSON(&body)
	if err != nil {
		t.Errorf("Error parsing JSON: %v", err)
	}
	if body.Foo != "foo" || body.Baz != "baz" || body.Fizz != "fizz" {
		t.Errorf("JSON parse incorrect: %v", body)
	}

}
