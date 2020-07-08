package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/arshabbir/grpcerror/protopb"
	"google.golang.org/grpc/status"
)

type server struct{}

type Server interface {
	Max(context.Context, *protopb.MaxRequest) (*protopb.MaxResponse, error)
}

func (*server) Max(ctx context.Context, req *protopb.MaxRequest) (*protopb.MaxResponse, error) {

	time.Sleep(1 * time.Second)

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.DeadlineExceeded, "Client cancled ")
	}

	num1 := req.GetNum1()
	num2 := req.GetNum2()

	if num1 == 0 || num2 == 0 {

		err := status.Error(codes.InvalidArgument, "Error string")

		return nil, err

	}

	if num1 > num2 {
		return &protopb.MaxResponse{Max: num1}, nil
	} else {
		return &protopb.MaxResponse{Max: num2}, nil
	}

}

var (
	ServerService Server = &server{}
)

func main() {

	lis, err := net.Listen("tcp", ":1001")

	if err != nil {
		log.Fatal("Error listening")
	}

	certFile := "../ssl/server.crt"
	keyFile := "../ssl/server.pem"
	creds, certerr := credentials.NewServerTLSFromFile(certFile, keyFile)

	if certerr != nil {
		log.Fatal("Error loading certificate : ", certerr)
		return
	}

	s := grpc.NewServer(grpc.Creds(creds))

	reflection.Register(s)
	protopb.RegisterMaxServiceServer(s, ServerService)

	log.Printf("Starting gRPC server....")
	if err := s.Serve(lis); err != nil {
		log.Fatal("Error serving ......")
	}

}
