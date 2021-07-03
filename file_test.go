package mojo_test

import (
	"bytes"
	"errors"
	"io"
	"os"
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

func TestFileOpen(t *testing.T) {
	temp, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatalf("Could not create temp file for testing Open(): %v", err)
	}
	file := mojo.NewFile(temp.Name()).Open()
	if v, ok := file.(io.Reader); !ok {
		t.Errorf("Open() does not return io.Reader. Got: %t", v)
	}
	if v, ok := file.(io.Writer); !ok {
		t.Errorf("Open() does not return io.Writer. Got: %t", v)
	}
	if v, ok := file.(io.Seeker); !ok {
		t.Errorf("Open() does not return io.Seeker. Got: %t", v)
	}
}

func TestTempFile(t *testing.T) {
	file := mojo.TempFile()
	_, err := os.Open(file.String())
	if err != nil {
		e := err.(*os.PathError)
		if errors.Is(e.Unwrap(), os.ErrNotExist) {
			t.Errorf("TempFile() did not create file: %v", err)
		} else {
			t.Errorf("TempFile() Open() caused unknown error: %v", err)
		}
	}
}

func TestFileSlurp(t *testing.T) {
	temp, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatalf("Could not create temp file for testing Slurp(): %v", err)
	}
	err = os.WriteFile(temp.Name(), []byte("foobar"), 0600)
	if err != nil {
		t.Fatalf("Could not write to temp file for testing Slurp(): %v", err)
	}

	file := mojo.NewFile(temp.Name())
	content := file.Slurp()
	if !bytes.Equal(content, []byte("foobar")) {
		t.Errorf("Slurp() read incorrect content. Got: %s; Expect: foobar", content)
	}
}

func TestFileSpurt(t *testing.T) {
	temp, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatalf("Could not create temp file for testing Spurt(): %v", err)
	}

	file := mojo.NewFile(temp.Name())
	file.Spurt([]byte("foobar"))

	content, err := os.ReadFile(temp.Name())
	if err != nil {
		t.Errorf("Error reading file written by Spurt(): %v", err)
	} else if !bytes.Equal(content, []byte("foobar")) {
		t.Errorf("Spurt() wrote incorrect content. Got: %s; Expect: foobar", content)
	}
}
