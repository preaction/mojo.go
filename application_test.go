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
	app.Routes.Get("/").To(func(*mojo.Context) {})
	mt := mojotest.Tester{T: t, App: &app}

	mt.GetOk("/", "Can get root").StatusIs(200)
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
