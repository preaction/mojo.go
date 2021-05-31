package util

import (
	"testing"
)

func TestStringSliceHas(t *testing.T) {
	strs := StringSlice{"foo", "bar", "baz"}

	has := strs.Has("bar")
	if !has {
		t.Errorf("Could not find string '%s' in %v", "bar", strs)
	}
}
