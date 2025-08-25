package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	grpcServer  *grpc.Server
	httpServer  *http.Server
	adminServer *http.Server

	options Options
}

// New ...
func New(options Options) *Server {
	return &Server{
		options: options,
	}
}

// Run ...
func (s *Server) Run(ctx context.Context, services ...Service) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s.initAdminHTTP()

	service := newCompositeService(services...)

	s.initGRPC(service)

	// nolint:contextcheck
	if err := s.initHTTP(service); err != nil {
		return fmt.Errorf("initHTTP: %w", err)
	}

	errChan := make(chan error)

	go func() {
		if err := s.runAdminHTTP(); err != nil {
			errChan <- fmt.Errorf("admin http: %w", err)
		}
	}()

	go func() {
		if err := s.runGRPC(); err != nil {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	go func() {
		if err := s.runHTTP(); err != nil {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	select {
	case err := <-errChan:
		cancel()
		s.shutdown() // nolint:contextcheck
		return err
	case <-ctx.Done():
		slog.Info("Shutting down servers...")
		s.shutdown() // nolint:contextcheck
		return nil
	}
}

func (s *Server) shutdown() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.grpcServer.GracefulStop()
	}()

	go func() {
		defer wg.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.options.HTTPShutdownTimeout)
		defer cancel()

		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("HTTP server shutdown error", "error", err)
		}

		if err := s.adminServer.Shutdown(shutdownCtx); err != nil {
			slog.Error("HTTP admin server shutdown error", "error", err)
		}
	}()

	wg.Wait()
}
