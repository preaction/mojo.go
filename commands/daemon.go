package main

import (
	"github.com/preaction/mojo.go"
)

func main() {
	srv := mojo.Server{
		App: &mojo.Application{},
	}
	srv.ListenAndServe()
}
