package mojo

// Parameters represents URL or form parameters
type Parameters map[string][]string

// Exists returns true if the parameter exists
func (p Parameters) Exists(name string) bool {
	_, ok := p[name]
	return ok && p[name] != nil
}

// Names returns the names of all parameters
func (p Parameters) Names() []string {
	names := make([]string, 0, len(p))
	for k := range p {
		names = append(names, k)
	}
	return names
}

// Param returns the first value for the given parameter, or the empty
// string if it does not exist.
func (p Parameters) Param(name string) string {
	if !p.Exists(name) {
		return ""
	}
	return p[name][0]
}

// EveryParam returns all values for the given parameter, or an empty
// array if it does not exist
func (p Parameters) EveryParam(name string) []string {
	if !p.Exists(name) {
		return []string{}
	}
	return p[name]
}
