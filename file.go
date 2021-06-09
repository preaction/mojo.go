package mojo

import (
	"strings"
)

// File represents a path on the filesystem
type File struct {
	path []string
}

func splitParts(pathParts ...string) []string {
	splitParts := []string{}
	for _, p := range pathParts {
		// XXX: Use os.PathSeparator, but also try to DTRT on WinNT...
		for _, pp := range strings.Split(p, "/") {
			splitParts = append(splitParts, pp)
		}
	}
	return splitParts
}

// NewFile creates a new File with the given path parts
func NewFile(pathParts ...string) File {
	return File{path: splitParts(pathParts...)}
}

// Dirname returns a new File representing the directory of the current
// file
func (f File) Dirname() File {
	if len(f.path) > 1 {
		return File{path: f.path[:len(f.path)-1]}
	}
	return File{path: []string{}}
}

// Child returns a new File for a child of the current file. This does not
// check for existence.
func (f File) Child(pathParts ...string) File {
	path := append(f.path, splitParts(pathParts...)...)
	return File{path: path}
}

// Sibling returns a new File for a sibling of the current file.
// Equivalent to file.Dirname().Child(...).
func (f File) Sibling(pathParts ...string) File {
	return f.Dirname().Child(pathParts...)
}

// String returns the file path as a string
func (f File) String() string {
	// XXX: Use os.PathSeparator
	// XXX: Add IsAbs
	return strings.Join(f.path, "/")
}
