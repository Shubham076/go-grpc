package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc/proto"
	"grpc/routes"
	"log"
	"net"
)

var add = "localhost:50051"

func main() {
	lis, err := net.Listen("tcp", add)
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("server listening at %s", add)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	proto.RegisterTestServer(grpcServer, &routes.Server{})
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}
