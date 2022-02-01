package main

import (
	"log"

	http_server "github.com/msm-k8s-svc-helper/pkg/api/http"
)

// Define server interface which can be implemented by http or grpc
type Server interface {
	StartAPI(listenAddr string)
}

func main() {
	listenAddr := ":9898"
	log.Printf("Starting server on %s", listenAddr)

	// start http server
	http_server.StartAPI(listenAddr)
}
