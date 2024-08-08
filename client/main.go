package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc/proto"
	"io"
	"log"
	"sync"
	"time"
)

var add = "localhost:50051"

func main() {
	conn, err := grpc.NewClient(add, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := proto.NewTestClient(conn)
	ctx := context.Background()

	// unary
	log.Println("------request-response--------")
	res, err := client.HealthCheck(ctx, &proto.Request{Id: 1})
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(res)

	// server streaming
	log.Println("----------Server streaming------------")
	serverStream, err := client.ServerStream(ctx, &proto.Request{Id: 10})
	for {
		res, err := serverStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res)
	}

	// client serverStream
	log.Println("----------Client streaming------------")
	reqs := []proto.Request{{Id: 10}, {Id: 100}}
	clientStream, err := client.ClientStream(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, req := range reqs {
		if err := clientStream.Send(&req); err != nil {
			log.Println(err)
		}
		log.Println("Sent", req.Id)
		time.Sleep(1 * time.Second)
	}
	res, err = clientStream.CloseAndRecv()
	if err != nil {
		log.Println(err)
	}
	log.Println(res)

	// Bi directional Stream
	log.Println("----------Bi-directional streaming------------")
	stream, err := client.BiDirectional(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	processBidrectional(stream)
	defer conn.Close()
}

func processBidrectional(stream proto.Test_BiDirectionalClient) {
	reqs := []*proto.Request{{Id: 20}, {Id: 30}}
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for _, req := range reqs {
			if err := stream.Send(req); err != nil {
				log.Println(err)
			}
			log.Println("Sent: ", req)
			time.Sleep(500 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				wg.Done()
				break
			}
			if err != nil {
				log.Println(err)
			}
			log.Println("Received: ", res)
		}
	}()
	wg.Wait()
}
