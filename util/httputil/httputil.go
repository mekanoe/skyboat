package httputil // import "skyboat.io/x/util/httputil"

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

// GetJSON gets the JSON out of the post body
func GetJSON(ctx *fasthttp.RequestCtx, iface interface{}) error {
	return json.Unmarshal(ctx.PostBody(), iface)
}

// Error is a convienience wrapper around logging and JSON output
func Error(msg string, err error, ctx *fasthttp.RequestCtx, log *logrus.Entry) {
	log.WithError(err).Errorln(msg)
	ctx.SetStatusCode(500)
	Write(ctx, map[string]interface{}{"success": false, "err": msg})
}

// Write writes JSON. Pretty simple.
func Write(ctx *fasthttp.RequestCtx, iface interface{}) error {
	data, err := json.Marshal(iface)
	if err != nil {
		return err
	}

	_, err = ctx.Write(data)
	return err
}
