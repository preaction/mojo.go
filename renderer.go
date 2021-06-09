package mojo

import (
	"fmt"
	"html/template"
	"strings"
)

// Renderer is an interface for template renderers. Implement this
// interface to integrate with another template system.
type Renderer interface {
	Render(what string, c Context) *string
}

// GoRenderer implements the Renderer interface using Go's built-in HTML
// Template system, changing the delimiters from "{{ ... }}" to "<% ...
// %>".
type GoRenderer struct {
	Paths     []File
	templates map[string]*template.Template
}

// template initializes a new Template object with the appropriate
// settings. Since Go requires that all functions be registered before
// the template is parsed, we must do this as late as possible...
func (ren *GoRenderer) template(name string) *template.Template {
	// Initialize a template with the correct settings.
	// XXX: Add FuncMap
	return template.New(name).Delims("<%", "%>")
}

// Add adds a template to the cache.
func (ren *GoRenderer) Add(name string, content string) {
	// XXX: Do we keep this or do something else to inject templates?
	if ren.templates == nil {
		ren.templates = map[string]*template.Template{}
	}
	t := ren.template(name)
	ren.templates[name] = template.Must(t.Parse(content))
}

// Render renders the named template using the data in the given
// context.
func (ren *GoRenderer) Render(name string, c Context) string {
	// Look up template in the cache
	t, ok := ren.templates[name]
	// XXX: If missing, look up template on the filesystem
	if !ok {
		panic(fmt.Sprintf("No template %s", name))
	}

	str := strings.Builder{}
	t.Execute(&str, c)

	return str.String()
}
