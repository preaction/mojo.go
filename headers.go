package mojo

import (
	"strconv"
	"strings"
	"time"
)

// Headers represents HTTP Headers, which are case-insensitive
type Headers map[string][]string

// Exists returns true if the header exists
func (h Headers) Exists(name string) bool {
	// XXX: toLower
	_, ok := h[name]
	return ok && h[name] != nil
}

// Header returns the first value for the given header, or the empty
// string if it does not exist.
func (h Headers) Header(name string) string {
	// XXX: toLower
	if !h.Exists(name) {
		return ""
	}
	return h[name][0]
}

// EveryHeader returns all values for the given Header, or an empty
// array if it does not exist
func (h Headers) EveryHeader(name string) []string {
	// XXX: toLower
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
// If-Modified-Since header. Otherwise, returns Time's zero value (see
// Time.IsZero() to detect this, or use Exists("If-Modified-Since")).
// Note: If-Modified-Since should be ignored if If-None-Match exists:
// https://datatracker.ietf.org/doc/html/rfc7232#section-3.3
func (h Headers) IfModifiedSince() time.Time {
	header := h.Header("If-Modified-Since")
	if header == "" {
		return time.Now()
	}
	t, err := time.Parse(time.RFC1123, header)
	if err != nil {
		// XXX: Log error?
		return time.Now()
	}
	return t
}

// Ranges returns an array of arrays containing start,end ranges
// requested from the Range header (in bytes). If the header does not
// exist, returns an empty array.
func (h Headers) Ranges() [][2]int64 {
	header := h.Header("Range")
	if !strings.HasPrefix(header, "bytes=") {
		return [][2]int64{}
	}
	parts := strings.Split(header[6:], ",")
	ranges := make([][2]int64, len(parts))
	for i, str := range parts {
		parts := strings.Split(str, "-")
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			start = -1
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			end = -1
		}
		ranges[i] = [2]int64{int64(start), int64(end)}
	}
	return ranges
}
