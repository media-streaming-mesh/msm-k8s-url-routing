package gRPC_api

import (
	context "context"
	"log"
	"net"

	url_handler "github.com/msm-k8s-svc-helper/pkg/url-handler"
	grpc "google.golang.org/grpc"
)

type Server struct {
	url_handler.UrlHandler // compose k8s api-server client
	UnimplementedGetEndpointServer
}

func (s *Server) Send(ctx context.Context, req *EndpointRequest) (*EndpointResponse, error) {
	urls := s.GetInternalURLs(req.GetReq())

	log.Printf("request handled: %s -> : %v", req.GetReq(), urls)
	return &EndpointResponse{Res: urls}, nil
}

func StartAPI(listenAddr string) {
	// start grpc server
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := grpc.NewServer()
	grpcServer := &Server{}
	grpcServer.InitializeUrlHandler()

	RegisterGetEndpointServer(gs, grpcServer)

	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
