package main

import (
	"github.com/preaction/mojo.go"
)

func main() {
	app := mojo.NewApplication()
	app.Routes.Get("/").To(func(c *mojo.Context) { c.Res.Body = "Hello, World!" })
	app.Start()
}
