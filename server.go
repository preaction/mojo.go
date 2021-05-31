package mojo

import (
	"net/http"
)

type Server struct {
	raw http.Server
}

func (srv *Server) ListenAndServe() error {
	srv.raw.Handler = srv
	return srv.raw.ListenAndServe()
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := Request{raw: r}

	// XXX: This should be handled by the Application class
	// XXX: Instantiate appropriate controller type
	c := Controller{Req: &req, Res: &Response{}}

	// XXX: Look up route and call handler
	c.Res.Code = 200

	// Write the response
	// XXX: Copy response headers
	w.WriteHeader(c.Res.Code)
	w.Write([]byte("Hello, World!"))
}
