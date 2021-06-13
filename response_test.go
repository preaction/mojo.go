package mojo_test

import (
	"testing"

	"github.com/preaction/mojo.go"
)

func TestNewResponse(t *testing.T) {
	res := mojo.NewResponse()
	if res.Headers == nil {
		t.Errorf("mojo.NewResponse() did not initialize Headers")
	}
}

func TestResponseJSON(t *testing.T) {
	type Employee struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	fry := Employee{Name: "Philip J. Fry", Email: "orangejoe@planex.nny"}
	expect := `{"name":"Philip J. Fry","email":"orangejoe@planex.nny"}`

	res := mojo.NewResponse()
	res.JSON(fry)

	if res.Content == "" {
		t.Errorf("JSON() left empty content")
	} else if res.Content != expect {
		t.Errorf("JSON() content incorrect; Expect: %s; Got: %s", expect, res.Content)
	}
}
