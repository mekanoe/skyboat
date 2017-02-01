package main

import (
	"log"
	rpc "github.com/kayteh/saving-light/cmd/sl-launcher/rpc"
	"golang.org/x/net/context"
)

type launcherServer struct {
}

func newLauncherServer() *launcherServer {

	return &launcherServer{}

}

func (l *launcherServer) Launch(ctx context.Context, in *rpc.LaunchRequest) (*rpc.LaunchedInstance, error) {

	log.Printf("got launch request for image: %s", in.Image)

	return &rpc.LaunchedInstance{
		NodeIp: "127.0.0.1",
		Port:   10000,
	}, nil

}
