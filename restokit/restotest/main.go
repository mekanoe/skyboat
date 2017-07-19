package main // import "skyboat.io/x/restokit/restotest"

import (
	"fmt"

	"skyboat.io/x/restokit"
	"skyboat.io/x/restokit/restotest/api"
)

func main() {
	resto := restokit.NewRestokit(":4665")
	api.FetchAPIRoutes(resto.Router)
	fmt.Println("started :4665")
	resto.Start()
}
