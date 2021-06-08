package mojo

import (
	"net"
	"net/http"
)

type Server struct {
	raw http.Server
	App *Application
}

func (srv *Server) Serve(l net.Listener) error {
	srv.raw.Handler = srv
	return srv.raw.Serve(l)
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	req.Read(r)
	res := &Response{Writer: w}
	c := srv.App.BuildContext(req, res)
	// XXX: Capture panics?
	srv.App.Handler(c)
}
