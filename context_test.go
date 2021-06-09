package mojo_test

import (
	"testing"

	"github.com/preaction/mojo.go"
)

func TestContextParam(t *testing.T) {
	req := mojo.NewRequest("GET", "/")
	// XXX: Add Params argument to mojo.NewRequest
	req.Params = map[string][]string{"foo": []string{"bar"}}

	c := mojo.Context{Req: req, Stash: map[string]interface{}{}}
	if c.Param("foo") != "bar" {
		t.Errorf("Param not read from request")
	}

	c.Stash["foo"] = "buzz"
	if c.Param("foo") != "buzz" {
		t.Errorf("Param not overridden by stash")
	}
}

func TestContextRenderString(t *testing.T) {
	app := mojo.NewApplication()
	app.Renderer.(*mojo.GoRenderer).Add("foo", "bar")

	req := mojo.NewRequest("GET", "/")
	// XXX: Add Params argument to mojo.NewRequest
	req.Params = map[string][]string{"foo": []string{"bar"}}

	c := app.BuildContext(req, &mojo.Response{})
	out := c.RenderString("foo")
	if out != "bar" {
		t.Errorf(`RenderString("foo") != "bar"; Got: %s`, out)
	}
}
