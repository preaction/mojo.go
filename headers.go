package mojo

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
