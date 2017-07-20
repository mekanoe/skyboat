package api // import "skyboat.io/x/launcher/api"

import (
	"fmt"

	"github.com/Sirupsen/logrus"
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
	log := ctx.UserValue("log").(*logrus.Entry)
	k8s := ctx.UserValue("k8s").(*kubernetes.Clientset)
	lr := &launchertypes.LaunchRequest{}

	err := httputil.GetJSON(ctx, lr)
	if err != nil {
		httputil.Error("json parse failed", err, ctx, log)
		return
	}

	id, err := ksuid.NewRandom()
	if err != nil {
		httputil.Error("ksuid creation", err, ctx, log)
		return
	}

	p := k8s.CoreV1().Pods("test")

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
		"skyboat.io/launched": "true",
		"skyboat.io/id":       idStr,
		"skyboat.io/ns":       lr.Namespace,
	}

	createdPod, err := p.Create(pod)
	if err != nil {
		httputil.Error("pod creation failed", err, ctx, log)
		return
	}

	httputil.Write(ctx, launchertypes.Response{
		Success: true,
		Payload: createdPod,
	})
}
