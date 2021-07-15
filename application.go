package mojo

import (
	"fmt"
	"os"
)

// Application is the main type for an application to use. Applications
// register route handlers and hook handlers, and contain all the
// global application configuration and tools that can be used by those
// handlers.
type Application struct {
	Home     File
	Routes   Routes
	Static   *Static
	Log      Log
	hooks    map[Hook][]HookHandler
	Commands map[string]Command
	Renderer Renderer
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
	// request has been dispatched to the Static or Routes dispatcher.
	// Good for rewriting outgoing responses and other post-processing
	// tasks.
	AfterDispatch Hook = "AfterDispatch"
	// AfterStatic is a hook that is called after the
	// Static renderer finds a file and prepares a response. Good for
	// post-processing static file responses.
	AfterStatic Hook = "AfterDispatch"
)

// NewApplication builds a basic Mojo application with the default set
// of commands and (TODO) plugins.
func NewApplication() *Application {
	home := os.Getenv("MOJO_HOME")
	if home == "" {
		var err error
		home, err = os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("Could not determine app home directory: %s. Set MOJO_HOME to override the default.", err))
		}
	}

	app := &Application{
		Commands: map[string]Command{},
		Renderer: &GoRenderer{},
		Log:      NewLog(),
		Static:   &Static{},
	}
	app.Commands["help"] = &HelpCommand{App: app}
	app.Commands["version"] = &VersionCommand{App: app}
	app.Commands["daemon"] = &DaemonCommand{App: app}
	app.Static.AddPath(NewFile(home).Child("public"))

	return app
}

// XXX: Use go embed to embed templates and static files

// BuildContext fills in the context from the given Request and Response
// objects, including setting the default stash values from the
// Application and any stash values that come from the Request
func (app *Application) BuildContext(req *Request, res *Response) *Context {
	if req == nil {
		req = &Request{}
	}
	if res == nil {
		res = NewResponse()
	}
	c := &Context{Req: req, Res: res, App: app, Stash: map[string]interface{}{}}

	// XXX: Add defaults from application
	// Set default stash values from request
	// XXX: Add URL to Request
	c.Stash["path"] = c.Req.URL.Path

	return c
}

// Hook registers a Hook handler.
func (app *Application) Hook(hook Hook, handler HookHandler) {
	if app.hooks == nil {
		app.hooks = map[Hook][]HookHandler{}
	}
	if _, ok := app.hooks[hook]; ok {
		app.hooks[hook] = append(app.hooks[hook], handler)
	} else {
		app.hooks[hook] = []HookHandler{handler}
	}
}

// emit emits a hook event with the given context
// XXX: We may need different arguments for future hooks
func (app *Application) emit(hook Hook, c *Context) {
	hooks, ok := app.hooks[hook]
	if !ok {
		return
	}
	for _, hook := range hooks {
		hook(c)
	}
}

// Handler handles the request for the given Context. This includes
// dispatching the BeforeDispatch and AfterDispatch hooks, trying the
// Static dispatch and Routes dispatch, and writing the response to the
// user (if it hasn't been already)
func (app *Application) Handler(c *Context) {
	app.emit(BeforeDispatch, c)
	if app.Static.Dispatch(c) {
		app.emit(AfterStatic, c)
	} else {
		app.Routes.Dispatch(c)
	}
	app.emit(AfterDispatch, c)

	// Write the response
	if !c.rendered {
		c.Render("")
	}
	if c.Res.Code == 0 {
		c.Res.Code = 200
	}
	if c.Res.Writer != nil {
		// XXX: Copy response headers
		c.Res.Writer.WriteHeader(c.Res.Code)
		// XXX: Build Body from whatever parts we have
		c.Res.Content.Serve(c.Res.Writer)
	}
}

// Start invokes the Application's commands using the arguments given on
// the command-line.
func (app *Application) Start() {
	name := os.Args[1]
	cmd, ok := app.Commands[name]
	if !ok {
		fmt.Printf("Command not found: %s\n", name)
		os.Exit(1)
	}
	cmd.Run(os.Args[2:])
}
