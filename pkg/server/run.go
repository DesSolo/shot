package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

func (s *Server) runGRPC() error {
	lis, err := net.Listen("tcp", s.options.GRPCEndpoint)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	slog.Info("gRPC server listening", "address", s.options.GRPCEndpoint)
	if err := s.grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("grpcServer.Serve: %w", err)
	}

	return nil
}

func (s *Server) runHTTP() error {
	slog.Info("HTTP server listening", "address", s.options.HTTPEndpoint)
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("httpServer.ListenAndServe: %w", err)
	}
	return nil
}

func (s *Server) runAdminHTTP() error {
	slog.Info("HTTP admin server listening", "address", s.options.HTTPAdminEndpoint)
	if err := s.adminServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("adminServer.ListenAndServe: %w", err)
	}
	return nil
}
