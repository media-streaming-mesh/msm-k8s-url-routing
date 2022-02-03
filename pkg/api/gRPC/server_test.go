package gRPC_api

import (
	"context"
	"log"
	"testing"

	grpc "google.golang.org/grpc"
)

const listenAddr = ":9898"
const grpcHost = "localhost:9898"

func TestGetInternalURLs(t *testing.T) {

	// Set up a connection to the server.
	conn, err := grpc.Dial(grpcHost, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	log.Print("connected...")
	defer conn.Close()

	client := NewGetEndpointClient(conn)
	resp, err := client.Send(context.TODO(), &EndpointRequest{Req: "http://kubernetes:443"})

	if err != nil {
		t.Fatalf("error requesting gRPC server: %v", err)
	}

	log.Println(resp.Res)
}
