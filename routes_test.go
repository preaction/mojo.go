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

func TestRoutesStash(t *testing.T) {
	gotStash := mojo.Stash{}
	expectValue := "bar"

	handlerCalled := false
	handler := func(c *mojo.Context) {
		gotStash = c.Stash
		handlerCalled = true
	}

	app := mojo.NewApplication()
	app.Routes.Get("/foo", mojo.Stash{"foo": "bar"}).To(handler)

	req := mojo.NewRequest("GET", "/foo")
	res := &mojo.Response{}
	c := app.BuildContext(req, res)

	app.Routes.Dispatch(c)

	if !handlerCalled {
		t.Errorf("Route handler not called by Dispatch()")
	}
	val, ok := gotStash["foo"]
	if !ok {
		t.Errorf("Handler stash does not contain route default")
		return
	}
	str, ok := val.(string)
	if !ok {
		t.Errorf("Handler stash value is not a string")
		return
	}
	if str != expectValue {
		t.Errorf(`Stash["foo"] != %s`, expectValue)
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
