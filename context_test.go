package mojo_test

import (
	"testing"

	"github.com/preaction/mojo.go"
)

func TestContextParam(t *testing.T) {
	req := mojo.Request{
		Params: mojo.Parameters{
			"foo": []string{"bar"},
		},
	}

	c := mojo.Context{Req: &req, Stash: map[string]interface{}{}}
	if c.Param("foo") != "bar" {
		t.Errorf("Param not read from request")
	}

	c.Stash["foo"] = "buzz"
	if c.Param("foo") != "buzz" {
		t.Errorf("Param not overridden by stash")
	}
}
