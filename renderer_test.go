package mojo_test

import (
	"fmt"
	"testing"

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
