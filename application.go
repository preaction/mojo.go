package mojo

// Application is the main type for an application to use.
type Application struct {
	Routes Routes
}

// XXX: Use go embed to embed templates and static files

func (app *Application) Handler(tx *Transaction) {
	// XXX: Instantiate appropriate controller type
	c := Controller{Req: &tx.Req, Res: &tx.Res}

	// XXX: Add defaults from application
	// Set default stash values from request
	// XXX: Add URL to Request
	c.Stash["path"] = c.Req.raw.URL.Path

	// Look up route and call handler
	app.Routes.Dispatch(&c)

	if c.Res.Code == 0 {
		c.Res.Code = 200
	}

	// Write the response
	// XXX: Copy response headers
	tx.Res.raw.WriteHeader(c.Res.Code)
	tx.Res.raw.Write([]byte("Hello, World!"))
}
