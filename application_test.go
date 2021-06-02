package mojo_test

import (
	"fmt"
	"net/url"

	"github.com/preaction/mojo.go"
)

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
