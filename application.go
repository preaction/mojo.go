package mojo

// Application is the main type for an application to use.
type Application struct {
	Routes Routes
}

// XXX: Use go embed to embed templates and static files

func (app *Application) BuildContext(c *Context) *Context {
	c.Stash = map[string]interface{}{}
	// XXX: Add defaults from application
	// Set default stash values from request
	// XXX: Add URL to Request
	c.Stash["path"] = c.Req.URL.Path

	// XXX: Invoke afterBuildContext hook

	return c
}

func (app *Application) Handler(c *Context) {
	// XXX: Invoke hooks
	// Look up route and call handler
	app.Routes.Dispatch(c)

	// Write the response
	if c.Res.Code == 0 {
		c.Res.Code = 200
	}
	// XXX: Copy response headers
	// XXX: Embed writer?
	c.Res.Writer.WriteHeader(c.Res.Code)
	// XXX: Build Body from whatever parts we have
	c.Res.Writer.Write([]byte(c.Res.Body))
}
