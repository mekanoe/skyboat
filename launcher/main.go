// This binary deal directly with the k8s cluster to launch game server instances.
package main // import "skyboat.io/x/launcher"

import (
	"github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"skyboat.io/x/launcher/api"
	"skyboat.io/x/restokit"
	"skyboat.io/x/util/k8sutil"
)

func injectK8S(h fasthttp.RequestHandler) fasthttp.RequestHandler {

	client, err := k8sutil.InClusterClient()
	if err != nil {
		logrus.Warnln("couldn't create kubernetes in-cluster client")
	}

	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetUserValue("k8s", client)
		h(ctx)
	}
}

func main() {
	resto := restokit.NewRestokit(":2390")
	api.FetchAPIRoutes(resto.Router)
	resto.AddGlobalMiddleware(injectK8S)

	logrus.Info("starting launcher :2390")
	resto.Start()
}
