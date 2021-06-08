package mojo_test

import (
	"net/http/httptest"
	"testing"

	"github.com/preaction/mojo.go"
)

func TestRoutesMatch(t *testing.T) {
	app := mojo.NewApplication()
	r := app.Routes
	r.Get("/foo")

	req := mojo.NewRequest("GET", "/foo")
	res := &mojo.Response{}
	c := app.BuildContext(req, res)

	r.Dispatch(c)

	if c.Match == nil {
		t.Errorf("No route found for request: %v", c.Req)
	}
}

func TestRoutesHandler(t *testing.T) {
	handlerCalled := false
	handler := func(c *mojo.Context) {
		handlerCalled = true
	}

	app := mojo.NewApplication()
	app.Routes.Get("/foo").To(handler)

	req := mojo.NewRequest("GET", "/foo")
	res := &mojo.Response{}
	c := app.BuildContext(req, res)

	app.Routes.Dispatch(c)

	if !handlerCalled {
		t.Errorf("Route handler not called by Dispatch()")
	}
}

func TestRoutesPlaceholder(t *testing.T) {
	name := ""
	handlerCalled := false
	handler := func(c *mojo.Context) {
		name = c.Stash["name"].(string)
		handlerCalled = true
	}

	app := mojo.NewApplication()
	app.Routes.Get("/foo/:name").To(handler)

	req := mojo.NewRequest("GET", "/foo/morbo")
	res := &mojo.Response{Writer: httptest.NewRecorder()}
	c := app.BuildContext(req, res)

	app.Routes.Dispatch(c)

	if !handlerCalled {
		t.Errorf("Route handler not called by Dispatch()")
	}
	if name != "morbo" {
		t.Errorf("Route placeholder not added to stash")
	}
}
