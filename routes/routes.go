package routes

import (
	"context"
	"fmt"
	"grpc/proto"
	"io"
	"log"
)

type Server struct {
	proto.TestServer
}

func (s *Server) HealthCheck(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	fmt.Println("Got health check request")
	return &proto.Response{Message: "ok", Id: req.Id}, nil
}

func (s *Server) ServerStream(in *proto.Request, stream proto.Test_ServerStreamServer) error {
	fmt.Println("Got server stream request with Id:", in.Id)
	for i := in.Id; i < in.Id+5; i++ {
		err := stream.Send(&proto.Response{Message: "ok", Id: i})
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func (s *Server) ClientStream(stream proto.Test_ClientStreamServer) error {
	fmt.Println("Got client stream request")
	var sum uint32
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		sum += msg.Id
	}

	err := stream.SendAndClose(&proto.Response{Message: "ok", Id: sum})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) BiDirectional(stream proto.Test_BiDirectionalServer) error {
	log.Println("Bi directional handler invoked")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = stream.Send(&proto.Response{Message: "ok", Id: req.Id})
		if err != nil {
			return err
		}
	}
	return nil
}
