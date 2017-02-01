// LAUNCHER
// The launcher creates pods, assigns a free port, then tells the master server about it.

package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	launcherRpc "github.com/kayteh/saving-light/cmd/sl-launcher/rpc"
	"github.com/kayteh/saving-light/etc"
)

func main() {

	log.Printf("sl-launcher %s", etc.Version)

	log.Print("starting gRPC server...")
	lis, err := net.Listen("tcp", ":7254")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := newLauncherServer()
	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)
	launcherRpc.RegisterLauncherServer(grpcServer, srv)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
