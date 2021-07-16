package mojo

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/preaction/mojo.go/util"
)

// Headers represents HTTP Headers, which are case-insensitive
type Headers map[string][]string

// Exists returns true if the header exists
func (h Headers) Exists(name string) bool {
	canonName := http.CanonicalHeaderKey(name)
	_, ok := h[canonName]
	return ok && h[canonName] != nil
}

// Header returns the first value for the given header, or the empty
// string if it does not exist.
func (h Headers) Header(name string) string {
	if !h.Exists(name) {
		return ""
	}
	return h[http.CanonicalHeaderKey(name)][0]
}

// EveryHeader returns all values for the given Header, or an empty
// array if it does not exist
func (h Headers) EveryHeader(name string) []string {
	if !h.Exists(name) {
		return []string{}
	}
	return h[name]
}

// Pairs returns an array of arrays of name, value strings
func (h Headers) Pairs() [][2]string {
	pairs := [][2]string{}
	for name, values := range h {
		for _, value := range values {
			pairs = append(pairs, [2]string{name, value})
		}
	}
	return pairs
}

// Add adds one or more lines to a header.
func (h Headers) Add(header string, values ...string) {
	name := http.CanonicalHeaderKey(header)
	if _, ok := h[name]; !ok {
		h[name] = []string{}
	}
	h[name] = append(h[name], values...)
}

// IfNoneMatch returns an ETag string if the request contains an
// If-None-Match header. Otherwise, returns the empty string.
func (h Headers) IfNoneMatch() string {
	// There can be multiple matched separated by commas
	val := h.Header("If-None-Match")
	if val == "" {
		return ""
	}
	tags := strings.Split(val, ",")
	etag := strings.Trim(tags[0], `" `)
	return etag
}

// IfModifiedSince returns a Time if the request contains an
// If-Modified-Since header. Otherwise, returns time.Now() to prevent
// a 304 Not Modified response. Use Exists("If-Modified-Since") to
// detect this header, if needed.
// Note: If-Modified-Since should be ignored if If-None-Match exists:
// https://datatracker.ietf.org/doc/html/rfc7232#section-3.3
func (h Headers) IfModifiedSince() time.Time {
	header := h.Header("If-Modified-Since")
	if header == "" {
		return time.Now()
	}
	t, err := http.ParseTime(header)
	if err != nil {
		// XXX: Log error?
		return time.Now()
	}
	return t
}

// Ranges returns the start and end range requested from the Range
// header (in bytes). If the header does not exist, returns -1, -1.
func (h Headers) Range() (int64, int64) {
	header := h.Header("Range")
	if !strings.HasPrefix(header, "bytes=") {
		return -1, -1
	}
	parts := strings.Split(header[6:], "-")
	var start, end int64
	if parts[0] == "" {
		start = -1
	} else {
		start, _ = strconv.ParseInt(parts[0], 10, 64)
	}
	if parts[1] == "" {
		end = -1
	} else {
		end, _ = strconv.ParseInt(parts[1], 10, 64)
	}

	return start, end
}

// LastModified returns a Time if the request contains an LastModified
// header. Otherwise, returns time.Time zero value.  Use time.IsZero()
// or Exists("If-Modified-Since") to detect this, if needed.
func (h Headers) LastModified() time.Time {
	header := h.Header("Last-Modified")
	if header == "" {
		return time.Time{}
	}
	t, err := http.ParseTime(header)
	if err != nil {
		// XXX: Log error?
		return time.Time{}
	}
	return t
}

// Etag returns the value of any ETag header, or the empty string.
func (h Headers) Etag() string {
	etag := h.Header("Etag")
	if etag == "" {
		return ""
	}
	return strings.Trim(etag, `" `)
}

// Authorization returns the base64-decoded value of the Authorization
// header, if any.
func (h Headers) Authorization() string {
	raw := h.Header("Authorization")
	if !strings.HasPrefix(raw, "Basic ") {
		return ""
	}
	encoded := strings.TrimPrefix(raw, "Basic ")
	return util.B64Decode(encoded)
}
