package restokit

import (
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/kayteh/spaceplane/etc"
)

// Restokit is the REST framework common building block.
// The system involves simple codegen tricks.
type Restokit struct {
	Router *fasthttprouter.Router
	Server *fasthttp.Server
	addr   string
}

// NewRestokit creates a new restokit with the specified address
func NewRestokit(addr string) *Restokit {
	resto := &Restokit{
		Router: fasthttprouter.New(),
		addr:   addr,
	}

	srv := &fasthttp.Server{
		Name:    fmt.Sprintf("spaceplane restokit/%s (%s)", etc.Version, etc.Ref),
		Handler: resto.Router.Handler,
	}

	resto.Server = srv
	return resto
}

// Start starts the server as built.
func (r *Restokit) Start() error {
	return r.Server.ListenAndServe(r.addr)
}
