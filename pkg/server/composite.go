package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type compositeService struct {
	services []Service
}

func newCompositeService(services ...Service) *compositeService {
	return &compositeService{
		services: services,
	}
}

// RegisterGRPC ...
func (c *compositeService) RegisterGRPC(server *grpc.Server) {
	for _, service := range c.services {
		service.RegisterGRPC(server)
	}
}

// RegisterGateway ...
func (c *compositeService) RegisterGateway(ctx context.Context, mux *runtime.ServeMux) error {
	for _, service := range c.services {
		if err := service.RegisterGateway(ctx, mux); err != nil {
			return err // nolint:wrapcheck
		}
	}

	return nil
}

func (c *compositeService) RegisterGatewayFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	for _, service := range c.services {
		if err := service.RegisterGatewayFromEndpoint(ctx, mux, endpoint, opts); err != nil {
			return err // nolint:wrapcheck
		}
	}

	return nil
}
