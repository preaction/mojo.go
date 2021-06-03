# mojo

package mojo is a web application framework, ported from the
Mojolicious web framework for the Perl language.

In the Mojolicious framework, an author builds an Application object,
registers handlers for one or more Routes, and then starts the
application using its (TODO) Start() method.

```go
package main
import "github.com/preaction/mojo.go"
func main() {
	app := mojo.NewApplication()
	app.Routes.Get( "/" ).To( RootHandler )
	app.Start()
}
func RootHandler( c *mojo.Context ) {
	c.Res.Body = "Hello, World!"
}
```

## Sub Packages

* [mojo](./mojo)

* [test](./test)

* [util](./util)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
