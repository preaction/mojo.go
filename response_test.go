package mojo_test

import (
	"testing"

	"github.com/preaction/mojo.go"
	"github.com/preaction/mojo.go/testmojo"
)

func TestNewResponse(t *testing.T) {
	res := mojo.NewResponse()
	if res.Headers == nil {
		t.Errorf("mojo.NewResponse() did not initialize Headers")
	}
}

func TestResponseRead(t *testing.T) {
	raw := testmojo.BuildHTTPResponse(t, `HTTP/1.1 200 OK
Content-Length: 5

Hello`)
	res := mojo.NewResponse()
	res.Read(raw)
	if res.Content.String() != "Hello" {
		t.Errorf("Read() response did not build content asset")
	}
	if res.Headers.Header("Content-Length") != "5" {
		t.Errorf("Read() response did not build headers")
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

	if res.Content.String() == "" {
		t.Errorf("JSON() left empty content")
	} else if res.Content.String() != expect {
		t.Errorf("JSON() content incorrect; Expect: %s; Got: %s", expect, res.Content)
	}
}
