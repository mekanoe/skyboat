// LAUNCHER
// The launcher creates pods, assigns a free port, then tells the master server about it.

package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	launcherRpc "github.com/kayteh/spaceplane/cmd/launcher/rpc"
	"github.com/kayteh/spaceplane/etc"
)

func main() {

	log.Printf("spaceplane launcher %s", etc.Version)

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
