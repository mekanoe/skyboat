package main

import (
	"log"

	rpc "github.com/kayteh/spaceplane/cmd/sl-launcher/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:7254", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	c := rpc.NewLauncherClient(conn)

	r, err := c.Launch(context.Background(), &rpc.LaunchRequest{Image: "katie/gtan-rush"})
	if err != nil {
		log.Fatalf("launch request failed: %v", err)
	}

	log.Printf("launch request done, server at %s:%d", r.NodeIp, r.Port)

}
