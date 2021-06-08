package mojo_test

import (
	"fmt"
	"testing"

	"github.com/preaction/mojo.go"
	mojotest "github.com/preaction/mojo.go/test"
)

func TestApplication(t *testing.T) {
	app := mojo.NewApplication()
	app.Routes.Get("/").To(func(c *mojo.Context) { c.Res.Body = "Hello, World!" })

	mt := mojotest.NewTester(t, app)
	mt.GetOk("/", "Can get root").StatusIs(200)
}

func TestApplicationHookBeforeDispatch(t *testing.T) {
	app := mojo.NewApplication()
	app.Hook(mojo.BeforeDispatch, func(c *mojo.Context) { c.Stash["who"] = "Mojolicious" })
	app.Hook(mojo.AfterDispatch, func(c *mojo.Context) {
		c.Res.Body = fmt.Sprintf("%s\n... and Goodbye!", c.Res.Body)
	})
	app.Routes.Get("/").To(func(c *mojo.Context) {
		c.Res.Body = fmt.Sprintf("Hello, %s!", c.Stash["who"])
	})

	mt := mojotest.NewTester(t, app)
	mt.GetOk("/", "Can get root").StatusIs(200)
	mt.TextIs("Hello, Mojolicious!\n... and Goodbye!")
}

func ExampleApplicationHelloWorld() {
	app := mojo.NewApplication()
	r := app.Routes.Get("/:name", mojo.Stash{"name": "World"})
	r.To(func(c *mojo.Context) {
		c.Res.Body = fmt.Sprintf("Hello, %s!\n", c.Stash["name"])
	})

	req := mojo.NewRequest("GET", "/")
	res := &mojo.Response{}
	c := app.BuildContext(req, res)
	app.Handler(c)
	fmt.Print(c.Res.Body)

	req = mojo.NewRequest("GET", "/Gophers")
	res = &mojo.Response{}
	c = app.BuildContext(req, res)
	app.Handler(c)
	fmt.Print(c.Res.Body)

	// Output:
	// Hello, World!
	// Hello, Gophers!
}
