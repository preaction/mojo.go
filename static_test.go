package mojo_test

import (
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	"github.com/preaction/mojo.go"
	mojotest "github.com/preaction/mojo.go/test"
	"github.com/preaction/mojo.go/util"
)

func TestStatic(t *testing.T) {
	testfs := fstest.MapFS{
		"hello.txt": &fstest.MapFile{
			Data: []byte("Hello, World"),
			Mode: 0644,
		},
	}
	s := mojo.Static{Paths: []fs.FS{testfs}}
	// Test basic response handling
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
	now := time.Now().Round(0)
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
	c.Req.Headers["If-Modified-Since"] = []string{now.Format(time.RFC850)}
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
	c.Req.Headers["If-Modified-Since"] = []string{now.Format(time.RFC850)}
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

func TestStaticCacheETag(t *testing.T) {
	now := time.Now().Round(0)
	etag := util.MD5Sum(now.Add(-time.Hour).Format(time.RFC3339))
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
	c.Req.Headers["If-None-Match"] = []string{etag}
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
	c.Req.Headers["If-None-Match"] = []string{etag}
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
