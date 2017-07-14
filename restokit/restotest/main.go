package main

import (
	"fmt"

	"github.com/kayteh/spaceplane/restokit"
	"github.com/kayteh/spaceplane/restokit/restotest/api"
)

func main() {
	resto := restokit.NewRestokit(":4665")
	api.FetchAPIRoutes(resto.Router)
	fmt.Println("started :4665")
	resto.Start()
}
