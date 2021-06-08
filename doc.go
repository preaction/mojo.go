// package mojo is a web application framework, ported from the
// Mojolicious web framework for the Perl language.
//
// In the Mojolicious framework, an author builds an Application object,
// registers handlers for one or more Routes, and then starts the
// application using its Start() method.
//
//		package main
//		import "github.com/preaction/mojo.go"
//		import "fmt"
//		func main() {
//			app := mojo.NewApplication()
//			app.Routes.Get( "/:name", mojo.Stash{"name": "World"}).To( GreetHandler )
//			app.Start()
//		}
//		func GreetHandler( c *mojo.Context ) {
//			c.Res.Body = fmt.Sprintf( "Hello, %s!", c.Stash["name"] )
//		}
//
package mojo
