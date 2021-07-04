package mojo

import (
	"fmt"
	"io"
	"os"
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

// Open returns an io.ReadWriteSeeker for the given file
func (f File) Open() io.ReadWriteSeeker {
	reader, err := os.OpenFile(f.String(), os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}
	return reader
}

// Stat returns info about the current path
func (f File) Stat() os.FileInfo {
	info, err := os.Stat(f.String())
	if err != nil {
		panic(err)
	}
	return info
}

// TempFile creates a new file in the system's default temp directory
// and returns the path to that file.
func TempFile() File {
	file, err := os.CreateTemp("", "*")
	if err != nil {
		panic(fmt.Sprintf("Could not create temp file: %v", err))
	}
	// XXX: Automatic clean up?
	return NewFile(file.Name())
}

// Slurp reads all data from the file
func (f File) Slurp() []byte {
	file := f.Open()
	buf, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("Could not slurp file: %v", err))
	}
	return buf
}

// Spurt writes data to the file, replacing any existing content.
func (f File) Spurt(buf []byte) {
	err := os.WriteFile(f.String(), buf, 0666)
	if err != nil {
		panic(fmt.Sprintf("Could not spurt file: %v", err))
	}
}
