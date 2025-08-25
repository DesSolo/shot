package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/felixge/fgprof"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"shot/pkg/server/swagger"
)

func (s *Server) initGRPC(service Service) {
	grpcServer := grpc.NewServer(s.options.GRPCServerOptions...)

	reflection.Register(grpcServer)

	service.RegisterGRPC(grpcServer)

	s.grpcServer = grpcServer
}

func (s *Server) initHTTP(service Service) error {
	mux := runtime.NewServeMux()

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if err := service.RegisterGatewayFromEndpoint(context.Background(), mux, s.options.GRPCEndpoint, dialOpts); err != nil {
		return fmt.Errorf("service.RegisterGateway: %w", err)
	}

	withCors := cors.New(s.options.CorsOptions).Handler(mux)

	s.httpServer = &http.Server{
		Addr:              s.options.HTTPEndpoint,
		Handler:           withCors,
		ReadHeaderTimeout: s.options.HTTPReadHeaderTimeout,
	}

	return nil
}

func (s *Server) initAdminHTTP() {
	mux := chi.NewMux()

	// profilers
	mux.Mount("/debug", middleware.Profiler())
	mux.Handle("/debug/fgprof", fgprof.Handler())

	// swagger
	if _, httpEndpointPort, err := net.SplitHostPort(s.options.HTTPEndpoint); err == nil {
		swagger.Init(mux, httpEndpointPort)
	}

	s.adminServer = &http.Server{
		Addr:              s.options.HTTPAdminEndpoint,
		Handler:           mux,
		ReadHeaderTimeout: s.options.HTTPReadHeaderTimeout,
	}
}
