package mojo_test

import (
	"testing"
	"time"

	"github.com/preaction/mojo.go"
	// XXX: Move test -> mojotest
	mojotest "github.com/preaction/mojo.go/test"
)

func TestHeadersRanges(t *testing.T) {
	raw := mojotest.BuildHTTPRequest(t, `GET /foo HTTP/1.1
Range: bytes=1-2

`)

	req := &mojo.Request{}
	req.Read(raw)

	start, end := req.Headers.Range()
	if start != 1 {
		t.Errorf("Start range not parsed correctly. Got: %d; Expect: 1", start)
	}
	if end != 2 {
		t.Errorf("End range not parsed correctly. Got: %d; Expect: 2", end)
	}
}

func TestHeadersIfModifiedSince(t *testing.T) {
	raw := mojotest.BuildHTTPRequest(t, `GET /foo HTTP/1.1
If-Modified-Since: Tue, 31 Dec 2999 23:30:00 GMT

`)

	req := &mojo.Request{}
	req.Read(raw)

	since := req.Headers.IfModifiedSince()
	gmt := time.FixedZone("GMT", 0)
	expect := time.Date(2999, 12, 31, 23, 30, 0, 0, gmt)
	if !expect.Equal(since) {
		t.Errorf("If-Modified-Since not parsed correctly. Got: %#v; Expect: %#v", since, expect)
	}
}

func TestHeadersIfNoneMatch(t *testing.T) {
	raw := mojotest.BuildHTTPRequest(t, `GET /foo HTTP/1.1
If-None-Match: acabacab

`)

	req := &mojo.Request{}
	req.Read(raw)

	etags := req.Headers.IfNoneMatch()
	expect := "acabacab"
	if etags != expect {
		t.Errorf("If-None-Match not parsed correctly Got: %#v; Expect: %#v", etags, expect)
	}
}
