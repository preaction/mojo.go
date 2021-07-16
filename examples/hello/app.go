package main

import (
	"fmt"

	"github.com/preaction/mojo.go"
)

func main() {
	app := mojo.NewApplication()
	app.Routes.Get("/:who", mojo.Stash{"who": "World"}).To(Greeting)
	app.Start()
}

func Greeting(c *mojo.Context) {
	who := c.Param("who")
	c.Res.Text(fmt.Sprintf("Hello, %s!", who))
}
