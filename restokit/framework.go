package restokit // import "skyboat.io/x/restokit"

import (
	"fmt"
	"net"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"github.com/Sirupsen/logrus"
	"github.com/segmentio/ksuid"
	"skyboat.io/x/etc"
	"skyboat.io/x/util"
)

type Middleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

// Restokit is the REST framework common building block.
// The system involves simple codegen tricks.
type Restokit struct {
	Router   *fasthttprouter.Router
	Server   *fasthttp.Server
	Listener net.Listener
	Logger   *logrus.Entry

	HealthCheck    fasthttp.RequestHandler
	ReadinessCheck fasthttp.RequestHandler

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
		Logger:         logrus.New().WithFields(logrus.Fields{}),
		HealthCheck:    defaultHealthCheck,
		ReadinessCheck: defaultReadinessCheck,
		addr:           addr,
	}

	prod, _ := util.Getenvdef("IS_PROD", false).Bool()
	if prod {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	r.AddGlobalMiddleware(r.logging)

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
	r.Router.GET("/+/healthz", r.HealthCheck)
	r.Router.GET("/+/readiness", r.ReadinessCheck)

	if r.Listener == nil {
		err = r.Server.ListenAndServe(r.addr)
	} else {
		err = r.Server.Serve(r.Listener)
	}

	return err
}

func defaultReadinessCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetUserValue("log:silent", true)
	ctx.SetStatusCode(200)
	ctx.WriteString("ok")
}

func defaultHealthCheck(ctx *fasthttp.RequestCtx) {
	ctx.SetUserValue("log:silent", true)
	ctx.SetStatusCode(200)
	ctx.WriteString("ok")
}

func (r *Restokit) logging(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startTime := time.Now()
		reqid, err := ksuid.NewRandom()
		if err != nil {
			r.Logger.WithError(err).Error("Error in logger: ksuid create")
			ctx.Error("unrecoverable error in logger", 500)
		}

		logEntry := r.Logger.WithField("reqid", reqid.String())

		ctx.SetUserValue("log", logEntry)
		ctx.SetUserValue("reqid", reqid)
		ctx.SetUserValue("log:silent", false)

		h(ctx)

		if !ctx.UserValue("log:silent").(bool) {
			logEntry.WithFields(logrus.Fields{
				"url":           string(ctx.URI().Path()),
				"method":        string(ctx.Request.Header.Method()),
				"referer":       string(ctx.Request.Header.Referer()),
				"code":          ctx.Response.StatusCode(),
				"user_agent":    string(ctx.Request.Header.UserAgent()),
				"bytes":         len(ctx.Response.Body()),
				"response_time": time.Since(startTime).Nanoseconds() / 1000,
			}).Infof("HTTP => %d %s %s", ctx.Response.StatusCode(), ctx.Request.Header.Method(), ctx.URI())
		}
	}
}
