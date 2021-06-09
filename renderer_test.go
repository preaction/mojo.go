package mojo_test

import (
	"github.com/preaction/mojo.go"
	"testing"
)

func TestGoRendererRender(t *testing.T) {
	r := mojo.GoRenderer{}
	r.Add("foo", "bar")
	out := r.Render("foo", mojo.Context{})
	if out != "bar" {
		t.Errorf(`Render("foo") != "bar"`)
	}
}
