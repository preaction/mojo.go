package mojo

import (
	"crypto/md5"
	"encoding/base32"
	"errors"
	"fmt"
	"io/fs"
	"time"
)

// APP_START is the time the application was started. This is used by
// the Static renderer to determine modified times for in-memory files.
var APP_START time.Time

func init() {
	APP_START = time.Now()
}

// Static handles requests for static files, including support for Range and HTTP
// caching (If-Modified-Since and If-None-Match).
type Static struct {
	Paths []fs.FS
}

// Dispatch tries to find a static file to handle the request. Returns
// true if a static file was found and served.
func (st *Static) Dispatch(c *Context) bool {
	// Remove the leading "/"
	path := c.Req.URL.Path[1:]
	return st.Serve(c, path)
}

// Serve tries to serve the static file at the given relative path.
// Returns true if a static file was found and served.
func (st *Static) Serve(c *Context, path string) bool {
	var file fs.File
	found := false
	for _, f := range st.Paths {
		var err error
		file, err = f.Open(path)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("Error opening file %s: %s\n", path, err)
			// XXX: Log an error
		}
		if err != nil {
			continue
		}
		// Found a readable file
		found = true
		break
	}
	if !found {
		return false
	}

	// Handle If-None-Match/If-Modified-Since
	if c.Req.Headers.Exists("If-None-Match") || c.Req.Headers.Exists("If-Modified-Since") {
		fstat, err := file.Stat()
		if err == nil {
			mtime := fstat.ModTime()
			// XXX: Create an MD5Sum utility method
			sum := md5.Sum([]byte(mtime.Format(time.RFC3339)))
			etag := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(sum[:])
			if c.Req.Headers.Exists("If-None-Match") && c.Req.Headers.IfNoneMatch() == etag {
				c.Res.Code = 304
				return true
			}
			if c.Req.Headers.Exists("If-Modified-Since") && c.Req.Headers.IfModifiedSince().After(mtime) {
				c.Res.Code = 304
				return true
			}
		}
	}

	c.Res.Content = NewAsset(file)

	// Handle Range request
	if c.Req.Headers.Exists("Range") {
		start, end := c.Req.Headers.Range()
		// XXX: Validate ranges else return 416
		c.Res.Code = 206
		c.Res.Content.Range(start, end)
		return true
	}

	// Serve entire file
	c.Res.Code = 200
	return true
}
