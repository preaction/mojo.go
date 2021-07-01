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
Range: bytes=1-2,-3

`)

	req := &mojo.Request{}
	req.Read(raw)

	r := req.Headers.Ranges()
	if len(r) != 2 {
		t.Errorf("Range header not split. Got: %d; Expect: 2", len(r))
		return
	}
	if r[0] != [2]int{1, 2} {
		t.Errorf("Start/End range not parsed correctly. Got: %#v; Expect: [2]int{1, 2}", r[0])
	}
	if r[1] != [2]int{-1, 3} {
		t.Errorf("Start/End range not parsed correctly. Got: %#v; Expect: [2]int{-1, 3}", r[1])
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
