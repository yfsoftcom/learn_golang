// the grpc server
package main

import (
	"fmt"
	"log"
	"net"

	// "strings"

	pb "grpc-foo/foo"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = "5009"
)

type server struct{}

func main() {

	fmt.Println("Startup")
	//起服务
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	s.Serve(lis)

	log.Println("grpc server in: %s", port)
}

func (t *server) SayHello(ctx context.Context, request *pb.HelloRequest) (response *pb.HelloReply, err error) {
	response = &pb.HelloReply{
		Message: request.Name + ":hi:test",
	}
	return response, nil
}
