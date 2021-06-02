package mojo_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/preaction/mojo.go"
	mojotest "github.com/preaction/mojo.go/test"
)

func TestApplication(t *testing.T) {
	app := mojo.Application{}
	app.Routes.Get("/").To(func(c *mojo.Context) { c.Res.Body = "Hello, World!" })

	mt := mojotest.Tester{T: t, App: &app}
	mt.GetOk("/", "Can get root").StatusIs(200)
}

func TestApplicationHookBeforeDispatch(t *testing.T) {
	app := mojo.Application{}
	app.Hook(mojo.BeforeDispatch, func(c *mojo.Context) { c.Stash["who"] = "Mojolicious" })
	app.Hook(mojo.AfterDispatch, func(c *mojo.Context) {
		c.Res.Body = fmt.Sprintf("%s\n... and Goodbye!", c.Res.Body)
	})
	app.Routes.Get("/").To(func(c *mojo.Context) {
		c.Res.Body = fmt.Sprintf("Hello, %s!", c.Stash["who"])
	})

	mt := mojotest.Tester{T: t, App: &app}
	mt.GetOk("/", "Can get root").StatusIs(200)
	mt.TextIs("Hello, Mojolicious!\n... and Goodbye!")
}

func ExampleApplicationHelloWorld() {
	app := mojo.Application{}
	app.Routes.Get("/").To(func(*mojo.Context) { fmt.Print("Hello, World!") })

	// XXX: Create mojo.URL wrapper
	// Everything needs a damned wrapper to be even remotely usable...
	url, _ := url.ParseRequestURI("/")
	c := mojo.Context{
		Req: &mojo.Request{Method: "GET", URL: url},
		Stash: map[string]interface{}{
			"path": url.Path,
		},
	}

	app.Routes.Dispatch(&c)

	// Output:
	// Hello, World!
}
