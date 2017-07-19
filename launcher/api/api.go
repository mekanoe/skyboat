package api // import "skyboat.io/x/launcher/api"

import (
	"fmt"

	"github.com/segmentio/ksuid"
	"github.com/valyala/fasthttp"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"skyboat.io/x/launcher/types"
	"skyboat.io/x/util/httputil"
)

//go:generate go run $PWD/restokit/codegen/codegen.go $CWD

// POST /launch
func postLaunch(ctx *fasthttp.RequestCtx) {
	k8s := ctx.UserValue("k8s").(*kubernetes.Clientset)
	lr := &launchertypes.LaunchRequest{}
	httputil.GetJSON(ctx, lr)

	id, err := ksuid.NewRandom()
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	p := k8s.CoreV1().Pods(lr.Namespace)

	/*
		apiVersion: v1
		kind: Pod
		metadata:
		  name: pod-example
		spec:
		  containers:
		  - image: ubuntu:trusty
		    command: ["echo"]
		    args: ["Hello World"]
	*/
	pod := &v1.Pod{
		Spec: lr.PodSpec,
	}

	idStr := id.String()

	pod.Name = fmt.Sprintf("%s-%s", lr.Namespace, idStr)
	pod.Annotations = map[string]string{
		"spaceplane/launched": "true",
		"spaceplane/id":       idStr,
		"spaceplane/ns":       lr.Namespace,
	}

	p.Create(pod)
}
