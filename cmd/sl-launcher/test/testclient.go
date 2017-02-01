package main

import (
	rpc "github.com/kayteh/saving-light/cmd/sl-launcher/rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func main() {

	conn, err := grpc.Dial("localhost:7254", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	c := rpc.NewLauncherClient(conn)

	r, err := c.Launch(context.Background(), &rpc.LaunchRequest{ Image: "quay.io/taranis/taranis:latest" })
	if err != nil {
		log.Fatalf("launch request failed: %v", err)
	}

	log.Printf("launch request done, server at %s:%d", r.NodeIp, r.Port)

}
