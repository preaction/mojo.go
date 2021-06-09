package mojo_test

import (
	"testing"

	"github.com/preaction/mojo.go"
)

func TestFile(t *testing.T) {
	f := mojo.NewFile("foo", "bar", "baz")
	if f.String() != "foo/bar/baz" {
		t.Errorf(`NewFile("foo", "bar", "baz") != "foo/bar/baz"; Got: %s`, f.String())
	}
}

func TestFileDirname(t *testing.T) {
	f := mojo.NewFile("foo", "bar", "baz")
	if f.Dirname().String() != "foo/bar" {
		t.Errorf(`Dirname() != "foo/bar"; Got: %s`, f.String())
	}
}

func TestFileChild(t *testing.T) {
	f := mojo.NewFile("foo", "bar")
	c := f.Child("baz")
	if c.String() != "foo/bar/baz" {
		t.Errorf(`Child("baz") != "foo/bar/baz"; Got: %s`, c.String())
	}
	c = f.Child("fizz/buzz")
	if c.String() != "foo/bar/fizz/buzz" {
		t.Errorf(`Child("fizz/buzz") != "foo/bar/fizz/buzz"; Got: %s`, c.String())
	}
}

func TestFileSibling(t *testing.T) {
	f := mojo.NewFile("foo", "bar")
	c := f.Sibling("fizz", "buzz")
	if c.String() != "foo/fizz/buzz" {
		t.Errorf(`Sibling("fizz", "buzz") != "foo/fizz/buzz"; Got: %s`, c.String())
	}
	c = f.Sibling("fizz/buzz")
	if c.String() != "foo/fizz/buzz" {
		t.Errorf(`Sibling("fizz/buzz") != "foo/fizz/buzz"; Got: %s`, c.String())
	}
}
