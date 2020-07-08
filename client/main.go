package main

import (
	"context"
	"log"
	"time"

	"github.com/arshabbir/grpcerror/protopb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	certKey := "../ssl/ca.crt"

	creds, cerr := credentials.NewClientTLSFromFile(certKey, "")

	if cerr != nil {

		log.Fatal("Error loading certificate")
		return
	}

	cc, err := grpc.Dial("localhost:1001", grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Printf("Error dailing")
	}

	client := protopb.NewMaxServiceClient(cc)

	log.Println(client.Max(context.Background(), &protopb.MaxRequest{Num1: 100, Num2: 200}))
	//client.

	log.Println(client.Max(context.Background(), &protopb.MaxRequest{Num1: 0, Num2: 200}))

	ctx2, cancle := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancle()

	resp, err := client.Max(ctx2, &protopb.MaxRequest{Num1: 1000, Num2: 200})

	if err != nil {

		resperr, ok := status.FromError(err)

		if ok {
			log.Println(resperr.Code(), "Error response")
			return

		} else {
			log.Fatal("Error in RPC level")
			return
		}

	}
	log.Println(resp)

}
