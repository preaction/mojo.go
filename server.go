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

func (srv *Server) BuildRequest(r *http.Request) *Request {
	return &Request{
		raw: r,
		URL: r.URL,
	}
}

func (srv *Server) BuildResponse(w http.ResponseWriter) *Response {
	return &Response{Writer: w}
}

func (srv *Server) BuildContext(w http.ResponseWriter, r *http.Request) *Context {
	c := &Context{
		Req: srv.BuildRequest(r),
		Res: srv.BuildResponse(w),
	}
	return srv.App.BuildContext(c)
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// XXX: Capture panics?
	srv.App.Handler(srv.BuildContext(w, r))
}
