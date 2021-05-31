package mojo

// Application is the main type for an application to use.
type Application struct {
}

func (app *Application) Handler(tx *Transaction) {
	// XXX: This should be handled by the Application class
	// XXX: Instantiate appropriate controller type
	c := Controller{Req: &tx.Req, Res: &tx.Res}

	// XXX: Look up route and call handler
	c.Res.Code = 200

	// Write the response
	// XXX: Copy response headers
	tx.Res.raw.WriteHeader(c.Res.Code)
	tx.Res.raw.Write([]byte("Hello, World!"))
}
