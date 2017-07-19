package restokit // import "skyboat.io/x/restokit"

import (
	"fmt"
	"net"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"skyboat.io/x/etc"
)

type Middleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

// Restokit is the REST framework common building block.
// The system involves simple codegen tricks.
type Restokit struct {
	Router   *fasthttprouter.Router
	Server   *fasthttp.Server
	Listener net.Listener

	middleware []Middleware

	addr string
}

var (
	serverName = fmt.Sprintf("spaceplane restokit/%s (%s)", etc.Version, etc.Ref)
)

// NewRestokit creates a new restokit with the specified address
func NewRestokit(addr string) *Restokit {
	r := &Restokit{
		Router: fasthttprouter.New(),
		Server: &fasthttp.Server{
			Name: serverName,
		},
		addr: addr,
	}

	return r
}

// AddGlobalMiddleware to the middleware stack. Only works before starting.
func (r *Restokit) AddGlobalMiddleware(fn Middleware) {
	r.middleware = append(r.middleware, fn)
}

func (r *Restokit) middlewareStack(initialHandler fasthttp.RequestHandler) fasthttp.RequestHandler {
	handler := initialHandler
	for _, mw := range r.middleware {
		handler = mw(handler)
	}
	return handler
}

// Start starts the server as built.
func (r *Restokit) Start() error {
	var err error

	r.Server.Handler = r.middlewareStack(r.Router.Handler)

	if r.Listener == nil {
		err = r.Server.ListenAndServe(r.addr)
	} else {
		err = r.Server.Serve(r.Listener)
	}

	return err
}
