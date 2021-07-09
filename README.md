# mojo

package mojo is a web application framework, ported from the
Mojolicious web framework for the Perl language.

In the Mojolicious framework, an author builds an Application object,
registers handlers for one or more Routes, and then starts the
application using its Start() method.

```go
package main
import "github.com/preaction/mojo.go"
import "fmt"
func main() {
	app := mojo.NewApplication()
	app.Renderer.AddTemplate("greet", "Hello, <% .Stash.name %>!\n")
	app.Routes.Get( "/:name", mojo.Stash{"name": "World"}).To( GreetHandler )
	app.Start()
}
func GreetHandler( c *mojo.Context ) {
	c.Render( "greet" )
}
```

## Sub Packages

* [mojo](./mojo)

* [test](./test)

* [util](./util)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
