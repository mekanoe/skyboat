package middleware

import "github.com/valyala/fasthttp"

// NoLogging turns off logging via middleware rather than manually setting the internal user value.
func NoLogging(f fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetUserValue("log:silent", true)
		f(ctx)
	}
}
