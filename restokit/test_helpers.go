package restokit

import (
	"log"
	"net/http"

	"github.com/hydrogen18/memlistener"
	"github.com/valyala/fasthttp"
)

func verboseHTTP(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
		log.Printf("http => %d %s %s", ctx.Response.StatusCode(), ctx.Method(), ctx.RequestURI())
	}
}

// ScaffoldHTTP creates an in-memory listener restokit with http client.
func ScaffoldHTTP() (*Restokit, *http.Client) {
	iml := memlistener.NewMemoryListener()

	resto := NewRestokit("127.0.0.1:40000")
	resto.Server.Handler = verboseHTTP(resto.Router.Handler)
	resto.Listener = iml

	tport := &http.Transport{}
	tport.Dial = iml.Dial

	client := &http.Client{
		// Force no redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},

		Transport: tport,
	}

	return resto, client
}

// TeardownHTTP for cleanup of test scaffolding
func TeardownHTTP(r *Restokit) error {
	return r.Listener.Close()
}
