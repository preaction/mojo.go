package mojo_test

import (
	"testing"

	"github.com/preaction/mojo.go"
	mojotest "github.com/preaction/mojo.go/test"
)

func TestRoutesMatch(t *testing.T) {
	router := &mojo.Routes{}
	router.Get("/foo")

	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/foo"))
	router.Dispatch(c)

	if c.Match == nil {
		t.Errorf("No route found for request: %v", c.Req)
	}
}

func TestRoutesHandler(t *testing.T) {
	handlerCalled := false
	handler := func(c *mojo.Context) {
		handlerCalled = true
	}

	router := &mojo.Routes{}
	router.Get("/foo").To(handler)

	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/foo"))
	router.Dispatch(c)

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

	router := &mojo.Routes{}
	router.Get("/foo", mojo.Stash{"foo": "bar"}).To(handler)

	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/foo"))
	router.Dispatch(c)

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

func TestRoutesStandardPlaceholder(t *testing.T) {
	name := ""
	handlerCalled := false
	handler := func(c *mojo.Context) {
		name = c.Stash["name"].(string)
		handlerCalled = true
	}

	router := &mojo.Routes{}
	router.Get("/foo/:name").To(handler)

	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/foo/morbo"))
	router.Dispatch(c)

	if !handlerCalled {
		t.Errorf("Route handler not called by Dispatch()")
	}
	if name != "morbo" {
		t.Errorf("Route placeholder not added to stash")
	}
}

func TestRoutesDelimitedPlaceholder(t *testing.T) {
	name := ""
	handlerCalled := false
	handler := func(c *mojo.Context) {
		name = c.Stash["name"].(string)
		handlerCalled = true
	}

	router := &mojo.Routes{}
	router.Get("/hello_<:name>").To(handler)

	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/hello_world"))
	router.Dispatch(c)

	if !handlerCalled {
		t.Errorf("Route handler not called by Dispatch()")
	}
	if name != "world" {
		t.Errorf("Route placeholder not added to stash")
	}
}

func TestRoutesOptionalPlaceholder(t *testing.T) {
	gotStash := mojo.Stash{}
	expectValue := "bar"

	handlerCalled := false
	handler := func(c *mojo.Context) {
		gotStash = c.Stash
		handlerCalled = true
	}

	router := &mojo.Routes{}
	router.Get("/foo/:foo", mojo.Stash{"foo": "bar"}).To(handler)

	c := mojotest.NewContext(t, mojo.NewRequest("GET", "/foo/"))
	router.Dispatch(c)

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

	t.Run("Trailing slash is optional", func(t *testing.T) {
		handlerCalled = false
		c := mojotest.NewContext(t, mojo.NewRequest("GET", "/foo"))
		router.Dispatch(c)

		if !handlerCalled {
			t.Errorf("Slash before optional placeholder is not optional")
		}
	})
}
