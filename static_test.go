package mojo_test

import (
	"fmt"
	"io/fs"
	"net/http"
	"testing"
	"testing/fstest"
	"time"

	"github.com/preaction/mojo.go"
	mojotest "github.com/preaction/mojo.go/test"
	"github.com/preaction/mojo.go/util"
)

func TestStatic(t *testing.T) {
	gmt := time.FixedZone("GMT", 0)
	expect := time.Date(2999, 12, 31, 23, 30, 0, 0, gmt)
	etag := util.MD5Sum(expect.Format(http.TimeFormat))

	testfs := fstest.MapFS{
		"hello.txt": &fstest.MapFile{
			Data:    []byte("Hello, World"),
			Mode:    0644,
			ModTime: expect,
		},
	}
	s := mojo.Static{Paths: []fs.FS{testfs}}
	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/hello.txt"))
	served := s.Dispatch(c)
	if !served {
		t.Fatalf("Static dispatch did not serve request")
	}
	if c.Res.Code != 200 {
		t.Errorf("Static dispatch set incorrect response code. Got: %d, Expect: %d", c.Res.Code, 200)
	}
	if c.Res.Content.String() != "Hello, World" {
		t.Errorf("Static dispatch got incorrect file. Got: %s, Expect: %s", c.Res.Content.String(), "Hello, World")
	}

	if !c.Res.Headers.LastModified().Equal(expect) {
		t.Errorf("Static dispatch sent incorrect Last-Modified. Got: %s, Expect: %s", c.Res.Headers.LastModified(), expect)
	}

	if c.Res.Headers.Etag() != etag {
		t.Errorf("Static dispatch sent incorrect Etag. Got: %s, Expect: %s", c.Res.Headers.Etag(), etag)
	}
}

func TestStaticPaths(t *testing.T) {
	firstfs := fstest.MapFS{
		"hello.txt": &fstest.MapFile{
			Data: []byte("Hello, Gophers"),
			Mode: 0644,
		},
	}
	secondfs := fstest.MapFS{
		"hello.txt": &fstest.MapFile{
			Data: []byte("Hello, World"),
			Mode: 0644,
		},
		"robots.txt": &fstest.MapFile{
			Data: []byte("Allow: *"),
			Mode: 0644,
		},
	}
	s := mojo.Static{Paths: []fs.FS{firstfs, secondfs}}

	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/hello.txt"))
	served := s.Dispatch(c)
	if !served {
		t.Fatalf("Static dispatch did not serve request")
	}
	if c.Res.Content.String() != "Hello, Gophers" {
		t.Errorf("Static dispatch got incorrect file. Got: %s, Expect: %s", c.Res.Content.String(), "Hello, Gophers")
	}

	c = mojotest.NewContext(t, mojo.NewRequest("GET", "/robots.txt"))
	served = s.Dispatch(c)
	if !served {
		t.Fatalf("Static dispatch did not serve request")
	}
	if c.Res.Content.String() != "Allow: *" {
		t.Errorf("Static dispatch got incorrect file. Got: %s, Expect: %s", c.Res.Content.String(), "Hello, Gophers")
	}
}

func TestStaticRange(t *testing.T) {
	testfs := fstest.MapFS{
		"hello.txt": &fstest.MapFile{
			Data: []byte("Hello, World"),
			Mode: 0644,
		},
	}
	s := mojo.Static{Paths: []fs.FS{testfs}}
	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/hello.txt"))
	c.Req.Headers["Range"] = []string{"bytes=0-4"}
	served := s.Dispatch(c)
	if !served {
		t.Fatalf("Static dispatch did not serve request")
	}
	if c.Res.Code != 206 {
		t.Errorf("Static dispatch set incorrect response code. Got: %d, Expect: %d", c.Res.Code, 206)
	}
	if c.Res.Content.String() != "Hello" {
		t.Errorf("Static dispatch got incorrect range. Got: %s, Expect: %s", c.Res.Content.String(), "Hello")
	}
}

func TestStaticCacheModifiedSince(t *testing.T) {
	now := time.Now().Round(0).UTC()
	testfs := fstest.MapFS{
		"old.txt": &fstest.MapFile{
			Data:    []byte("Old"),
			Mode:    0644,
			ModTime: now.Add(-time.Hour),
		},
		"new.txt": &fstest.MapFile{
			Data:    []byte("New"),
			Mode:    0644,
			ModTime: now,
		},
	}
	s := mojo.Static{Paths: []fs.FS{testfs}}

	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/old.txt"))
	c.Req.Headers.Add("If-Modified-Since", now.Format(http.TimeFormat))
	served := s.Dispatch(c)
	if !served {
		t.Fatalf("Static dispatch did not serve request")
	}
	if c.Res.Code != 304 {
		t.Errorf("Static dispatch set incorrect response code. Got: %d, Expect: %d", c.Res.Code, 304)
	}
	if c.Res.Content.String() != "" {
		t.Errorf("Static dispatch served file: %s", c.Res.Content.String())
	}

	c = mojotest.NewContext(t, mojo.NewRequest("GET", "/new.txt"))
	c.Req.Headers.Add("If-Modified-Since", now.Format(http.TimeFormat))
	served = s.Dispatch(c)
	if !served {
		t.Fatalf("Static dispatch did not serve request")
	}
	if c.Res.Code != 200 {
		t.Errorf("Static dispatch set incorrect response code. Got: %d, Expect: %d", c.Res.Code, 200)
	}
	if c.Res.Content.String() != "New" {
		t.Errorf("Static dispatch got incorrect file. Got: %s, Expect: %s", c.Res.Content.String(), "New")
	}
}

func TestStaticCacheEtag(t *testing.T) {
	now := time.Now().Round(0)
	etag := util.MD5Sum(now.Add(-time.Hour).Format(http.TimeFormat))
	testfs := fstest.MapFS{
		"old.txt": &fstest.MapFile{
			Data:    []byte("Old"),
			Mode:    0644,
			ModTime: now.Add(-time.Hour),
		},
		"new.txt": &fstest.MapFile{
			Data:    []byte("New"),
			Mode:    0644,
			ModTime: now,
		},
	}
	s := mojo.Static{Paths: []fs.FS{testfs}}

	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/old.txt"))
	c.Req.Headers.Add("If-None-Match", fmt.Sprintf("\"%s\"", etag))
	served := s.Dispatch(c)
	if !served {
		t.Fatalf("Static dispatch did not serve request")
	}
	if c.Res.Code != 304 {
		t.Errorf("Static dispatch set incorrect response code. Got: %d, Expect: %d", c.Res.Code, 304)
	}
	if c.Res.Content.String() != "" {
		t.Errorf("Static dispatch served file: %s", c.Res.Content.String())
	}

	c = mojotest.NewContext(t, mojo.NewRequest("GET", "/new.txt"))
	c.Req.Headers.Add("If-None-Match", fmt.Sprintf("\"%s\"", etag))
	served = s.Dispatch(c)
	if !served {
		t.Fatalf("Static dispatch did not serve request")
	}
	if c.Res.Code != 200 {
		t.Errorf("Static dispatch set incorrect response code. Got: %d, Expect: %d", c.Res.Code, 200)
	}
	if c.Res.Content.String() != "New" {
		t.Errorf("Static dispatch got incorrect file. Got: %s, Expect: %s", c.Res.Content.String(), "New")
	}
}
