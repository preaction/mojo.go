package mojo_test

import (
	"fmt"
	"testing"

	"github.com/preaction/mojo.go"
	mojotest "github.com/preaction/mojo.go/test"
)

func TestApplication(t *testing.T) {
	app := mojo.NewApplication()
	app.Routes.Get("/").To(func(c *mojo.Context) { c.Res.Text("Hello, World!") })

	mt := mojotest.NewTester(t, app)
	mt.GetOk("/", "Can get root").StatusIs(200)
}

func TestApplicationHookBeforeDispatch(t *testing.T) {
	app := mojo.NewApplication()
	app.Hook(mojo.BeforeDispatch, func(c *mojo.Context) { c.Stash["who"] = "Mojolicious" })
	app.Hook(mojo.AfterDispatch, func(c *mojo.Context) {
		c.Res.Content.AddChunk([]byte("\n... and Goodbye!"))
	})
	app.Routes.Get("/").To(func(c *mojo.Context) {
		c.Res.Text(fmt.Sprintf("Hello, %s!", c.Stash["who"]))
	})

	mt := mojotest.NewTester(t, app)
	mt.GetOk("/", "Can get root").StatusIs(200)
	mt.TextIs("Hello, Mojolicious!\n... and Goodbye!")
}

func ExampleApplicationHelloWorld() {
	app := mojo.NewApplication()
	app.Renderer.AddTemplate("greet", "Hello, <% .Stash.name %>!\n")

	r := app.Routes.Get("/:name", mojo.Stash{"name": "World"})
	r.To(func(c *mojo.Context) {
		c.Render("greet")
	})

	req := mojo.NewRequest("GET", "/")
	res := mojo.NewResponse()
	c := app.BuildContext(req, res)
	app.Handler(c)
	fmt.Print(c.Res.Content.String())

	req = mojo.NewRequest("GET", "/Gophers")
	res = mojo.NewResponse()
	c = app.BuildContext(req, res)
	app.Handler(c)
	fmt.Print(c.Res.Content.String())

	// Output:
	// Hello, World!
	// Hello, Gophers!
}

func ExampleApplicationJSON() {
	type Employee struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	employees := map[string]Employee{}
	app := mojo.NewApplication()

	r := app.Routes.Put("/employee/:id")
	r.To(func(c *mojo.Context) {
		id := c.Param("id")
		update := Employee{}
		if err := c.Req.JSON(&update); err != nil {
			c.Stash["status"] = 400
			c.Stash["content"] = err
			return
		}
		employees[id] = update
		c.Stash["status"] = 204
	})

	r = app.Routes.Get("/employee/:id")
	r.To(func(c *mojo.Context) {
		id := c.Param("id")
		if _, ok := employees[id]; !ok {
			c.Stash["status"] = 404
			return
		}
		c.Res.JSON(employees[id])
	})

	req := mojo.NewRequest("PUT", "/employee/fry")
	req.Content = mojo.NewAsset(`{"name":"Philip J. Fry","email":"orangejoe@planex.com"}`)
	res := mojo.NewResponse()
	c := app.BuildContext(req, res)
	app.Handler(c)
	fmt.Printf("%+v\n", employees)

	req = mojo.NewRequest("GET", "/employee/fry")
	res = mojo.NewResponse()
	c = app.BuildContext(req, res)
	app.Handler(c)
	fmt.Printf("%s\n", c.Res.Content.String())
	fmt.Printf("%s\n", c.Res.Headers.Header("Content-Type"))

	req = mojo.NewRequest("GET", "/employee/bender")
	res = mojo.NewResponse()
	c = app.BuildContext(req, res)
	app.Handler(c)
	fmt.Printf("%d\n", c.Res.Code)

	// Output:
	// map[fry:{Name:Philip J. Fry Email:orangejoe@planex.com}]
	// {"name":"Philip J. Fry","email":"orangejoe@planex.com"}
	// application/json
	// 404
}
