package main

import (
	"flag"
	"log"

	grpc_server "github.com/msm-k8s-svc-helper/pkg/api/gRPC"
	http_server "github.com/msm-k8s-svc-helper/pkg/api/http"
)

const listenAddr = ":9898"

func main() {
	log.Printf("Starting server on %s", listenAddr)

	protocol := flag.String("transport", "http", "a string")
	flag.Parse()

	log.Println(*protocol)

	switch *protocol {
	case "http":
		http_server.StartAPI(listenAddr)
	case "grpc":
		grpc_server.StartAPI(listenAddr)
	default:
		log.Fatalf("protocol %s not supported", *protocol)
	}
}
