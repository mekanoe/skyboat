package httputil

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// GetJSON gets the JSON out of the post body
func GetJSON(ctx *fasthttp.RequestCtx, iface interface{}) error {
	return json.Unmarshal(ctx.PostBody(), iface)
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
