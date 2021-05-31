package mojo

import (
	"net/http"
)

type Server struct {
	raw http.Server
	App *Application
}

func (srv *Server) ListenAndServe() error {
	srv.raw.Handler = srv
	return srv.raw.ListenAndServe()
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tx := Transaction{
		Req: Request{raw: r},
		Res: Response{raw: w},
	}

	// XXX: Capture panics?
	srv.App.Handler(&tx)
}
