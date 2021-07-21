package mojo_test

import (
	"testing"
	"time"

	"github.com/preaction/mojo.go"
	"github.com/preaction/mojo.go/testmojo"
)

func TestHeadersRanges(t *testing.T) {
	raw := testmojo.BuildHTTPRequest(t, `GET /foo HTTP/1.1
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
	raw := testmojo.BuildHTTPRequest(t, `GET /foo HTTP/1.1
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
	raw := testmojo.BuildHTTPRequest(t, `GET /foo HTTP/1.1
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

func TestHeadersLastModified(t *testing.T) {
	raw := testmojo.BuildHTTPResponse(t, `HTTP/1.1 200 OK
Last-Modified: Tue, 31 Dec 2999 23:30:00 GMT
Content-Length: 8

Content
`)

	res := &mojo.Response{}
	res.Read(raw)

	modtime := res.Headers.LastModified()
	gmt := time.FixedZone("GMT", 0)
	expect := time.Date(2999, 12, 31, 23, 30, 0, 0, gmt)
	if !expect.Equal(modtime) {
		t.Errorf("Last-Modified not parsed correctly. Got: %#v; Expect: %#v", modtime, expect)
	}
}

func TestHeadersEtag(t *testing.T) {
	raw := testmojo.BuildHTTPResponse(t, `HTTP/1.1 200 OK 
ETag: "acabacab"
Content-Length: 8

Content
`)

	res := &mojo.Response{}
	res.Read(raw)

	etags := res.Headers.Etag()
	expect := "acabacab"
	if etags != expect {
		t.Errorf("Etag not parsed correctly Got: %#v; Expect: %#v", etags, expect)
	}
}

func TestHeadersAuthorization(t *testing.T) {
	raw := testmojo.BuildHTTPRequest(t, `GET /foo HTTP/1.1
Authorization: Basic QmVuZGVyOnJvY2tz

`)
	req := &mojo.Request{}
	req.Read(raw)

	auth := req.Headers.Authorization()
	expect := "Bender:rocks"
	if auth != expect {
		t.Errorf("Authorization not parsed correctly. Got %v; Expect: %v", auth, expect)
	}
}
