package mojo_test

import (
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/preaction/mojo.go"
	mojotest "github.com/preaction/mojo.go/test"
)

func TestGoRendererRender(t *testing.T) {
	r := mojo.GoRenderer{}
	r.AddTemplate("foo", "bar")
	out := r.Render("foo", &mojo.Context{})
	if out != "bar" {
		t.Errorf(`Render("foo") != "bar"`)
	}
}

func TestGoRendererHelpers(t *testing.T) {
	r := mojo.GoRenderer{}
	r.AddHelper("greet", func(who string) string {
		return fmt.Sprintf("Hello, %s!", who)
	})
	r.AddTemplate("foo", `<% greet .Stash.who %>`)

	c := mojotest.NewContext(t, mojo.Stash{"who": "PHILIP J. FRY"})
	out := r.Render("foo", c)
	if out != "Hello, PHILIP J. FRY!" {
		t.Errorf(`Render("foo") failed. Expect: "Hello, PHILIP J. FRY!"; Got: %v`, out)
	}
}

func TestGoRendererFS(t *testing.T) {
	testFS := fstest.MapFS{
		"foo.html.gt": &fstest.MapFile{
			Data: []byte("Hello!"),
		},
	}
	r := mojo.GoRenderer{}
	r.AddFS(testFS)
	c := mojotest.NewContext(t)
	out := r.Render("foo.html.gt", c)
	if out != "Hello!" {
		t.Errorf(`Render("foo.html.gt") failed. Expect: "Hello!"; Got: %v`, out)
	}
}

func TestGoRendererFSOrder(t *testing.T) {
	lastFS := fstest.MapFS{
		"foo.html.gt": &fstest.MapFile{
			Data: []byte("ERROR!"),
		},
		"bar.html.gt": &fstest.MapFile{
			Data: []byte("Goodbye!"),
		},
	}
	firstFS := fstest.MapFS{
		"foo.html.gt": &fstest.MapFile{
			Data: []byte("Hello!"),
		},
	}
	r := mojo.GoRenderer{}
	r.AddFS(lastFS)
	r.AddFS(firstFS)
	c := mojotest.NewContext(t)
	out := r.Render("foo.html.gt", c)
	if out != "Hello!" {
		t.Errorf(`Render("foo.html.gt") failed. Expect: "Hello!"; Got: %v`, out)
	}
	out = r.Render("bar.html.gt", c)
	if out != "Goodbye!" {
		t.Errorf(`Render("bar.html.gt") failed. Expect: "Goodbye!"; Got: %v`, out)
	}
}
