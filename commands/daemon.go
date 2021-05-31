package main

import (
	"github.com/preaction/mojo.go"
)

func main() {
	srv := mojo.Server{}
	srv.ListenAndServe()
}
