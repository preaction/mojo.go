package mojo

// Application is the main type for an application to use.
type Application struct {
	Routes Routes
	Hooks  map[Hook][]HookHandler
}

type Hook string
type HookHandler func(*Context)

const (
	BeforeDispatch Hook = "BeforeDispatch"
	AfterDispatch  Hook = "AfterDispatch"
)

// XXX: Use go embed to embed templates and static files

func (app *Application) BuildContext(c *Context) *Context {
	c.Stash = map[string]interface{}{}
	// XXX: Add defaults from application
	// Set default stash values from request
	// XXX: Add URL to Request
	c.Stash["path"] = c.Req.URL.Path

	return c
}

func (app *Application) Hook(hook Hook, handler HookHandler) {
	if app.Hooks == nil {
		app.Hooks = map[Hook][]HookHandler{}
	}
	if _, ok := app.Hooks[hook]; ok {
		app.Hooks[hook] = append(app.Hooks[hook], handler)
	} else {
		app.Hooks[hook] = []HookHandler{handler}
	}
}

func (app *Application) Emit(hook Hook, c *Context) {
	hooks, ok := app.Hooks[hook]
	if !ok {
		return
	}
	for _, hook := range hooks {
		hook(c)
	}
}

func (app *Application) Handler(c *Context) {
	// Look up route and call handler
	app.Emit(BeforeDispatch, c)
	app.Routes.Dispatch(c)
	app.Emit(AfterDispatch, c)

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
