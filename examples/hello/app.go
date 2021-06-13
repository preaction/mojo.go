package main

import (
	"fmt"

	"github.com/preaction/mojo.go"
)

func main() {
	app := mojo.NewApplication()
	app.Routes.Get("/:who").To(Greeting)
	app.Start()
}

func Greeting(c *mojo.Context) {
	who := c.Param("who")
	if who == "" {
		who = "World"
	}
	c.Res.Content = fmt.Sprintf("Hello, %s!", who)
}
