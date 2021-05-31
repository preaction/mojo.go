package mojo

import (
	"net/http"
)

type Response struct {
	raw     http.ResponseWriter
	msg     Message
	Code    int
	Message string
}
