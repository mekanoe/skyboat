package launchertypes // import "skyboat.io/x/launcher/types"

import "k8s.io/api/core/v1"

type LaunchRequest struct {
	Metadata  map[string]interface{}
	Namespace string
	PodSpec   v1.PodSpec
}

type LaunchedInstance struct {
	IP   string
	Port int
}

type Response struct {
	Success bool
	Payload interface{}
}
