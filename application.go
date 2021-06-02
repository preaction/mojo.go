package mojo

// Application is the main type for an application to use. Applications
// register route handlers and hook handlers, and contain all the
// global application configuration and tools that can be used by those
// handlers.
type Application struct {
	Routes Routes
	Hooks  map[Hook][]HookHandler
}

// Hook is label for an application event that can have HookHandlers
// assigned to it. These handlers can modify the Context to modify the
// Request or Response.
type Hook string

// HookHandler is a function that is registered to be called for a given
// Hook
type HookHandler func(*Context)

const (
	// BeforeDispatch is a hook that is called before the
	// Application.Dispatch function is called for the current Context
	BeforeDispatch Hook = "BeforeDispatch"
	// AfterDispatch is a hook that is called after the
	// Application.Dispatch function is called for the current Context
	AfterDispatch Hook = "AfterDispatch"
)

// XXX: Use go embed to embed templates and static files

// BuildContext fills in the context from its Request object, including
// setting the default stash values from the Application and any stash
// values that come from the Request
func (app *Application) BuildContext(c *Context) *Context {
	c.Stash = map[string]interface{}{}
	// XXX: Add defaults from application
	// Set default stash values from request
	// XXX: Add URL to Request
	c.Stash["path"] = c.Req.URL.Path

	return c
}

// Hook registers a Hook handler.
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

// emit emits a hook event with the given context
// XXX: We may need different arguments for future hooks
func (app *Application) emit(hook Hook, c *Context) {
	hooks, ok := app.Hooks[hook]
	if !ok {
		return
	}
	for _, hook := range hooks {
		hook(c)
	}
}

// Handler handles the request for the given Context. This includes
// dispatching the BeforeDispatch and AfterDispatch hooks and writing
// the response to the user (if it hasn't been already)
func (app *Application) Handler(c *Context) {
	// Look up route and call handler
	app.emit(BeforeDispatch, c)
	app.Routes.Dispatch(c)
	app.emit(AfterDispatch, c)

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
