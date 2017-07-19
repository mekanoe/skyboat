package main

import "skyboat.io/x/restokit"

func main() {
	resto := restokit.NewRestokit(":2391")

	api.FetchAPIMethods(resto.Router)
	// resto.AddGlobalMiddleware(injectK8S)

	resto.Start()
}
