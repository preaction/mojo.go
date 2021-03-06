package mojo

import (
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"strings"
)

// Renderer is an interface for template renderers. Implement this
// interface to integrate with another template system.
type Renderer interface {
	AddTemplate(name string, content string)
	AddHelper(name string, f interface{})
	AddFS(fs fs.FS)
	Render(name string, c *Context) string
}

// GoRenderer implements the Renderer interface using Go's built-in HTML
// Template system, changing the delimiters from "{{ ... }}" to "<% ...
// %>".
type GoRenderer struct {
	fs        []fs.FS
	helpers   map[string]interface{}
	templates map[string]string
}

// AddHelper adds a template function with the given name.
func (ren *GoRenderer) AddHelper(name string, f interface{}) {
	if ren.helpers == nil {
		ren.helpers = map[string]interface{}{}
	}
	ren.helpers[name] = f
}

// template initializes a new Template object with the appropriate
// settings.
func (ren *GoRenderer) template(name string) *template.Template {
	// Initialize a template with the correct settings.
	return template.New(name).Delims("<%", "%>")
}

// AddTemplate adds a template to the cache.
func (ren *GoRenderer) AddTemplate(name string, content string) {
	// XXX: Do we keep this or do something else to inject templates?
	if ren.templates == nil {
		ren.templates = map[string]string{}
	}
	ren.templates[name] = content
}

// AddPath adds a path to look up templates.
func (ren *GoRenderer) AddPath(f File) {
	ren.fs = append([]fs.FS{os.DirFS(f.String())}, ren.fs...)
}

// AddFS adds a filesystem to look up templates.
func (ren *GoRenderer) AddFS(f fs.FS) {
	ren.fs = append([]fs.FS{f}, ren.fs...)
}

// Render renders the named template using the data in the given
// context.
func (ren *GoRenderer) Render(name string, c *Context) string {
	// Look up content in the cache
	content, ok := ren.templates[name]
	// If missing, look up template from the available filesystems
	if !ok {
		for _, f := range ren.fs {
			bytes, err := fs.ReadFile(f, name)
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				panic(fmt.Sprintf("Could not read template file: %s", err))
			} else if err != nil {
				continue
			}
			content = string(bytes)
			ren.AddTemplate(name, content)
			break
		}
	}

	t := ren.template(name).Funcs(ren.helpers)
	template.Must(t.Parse(content))

	str := strings.Builder{}
	t.Execute(&str, c)

	return str.String()
}
