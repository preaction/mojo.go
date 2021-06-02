package mojo

import (
	"net/http"
)

type Response struct {
	Writer  http.ResponseWriter
	msg     Message
	Code    int
	Message string
}
