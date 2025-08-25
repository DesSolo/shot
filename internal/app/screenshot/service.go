package screenshot

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	desc "shot/pkg/api/screenshot"
)

// Provider ...
type Provider interface {
	Image(ctx context.Context, url string) ([]byte, error)
}

// Implementation ...
type Implementation struct {
	desc.UnimplementedScreenshotServer

	provider Provider
}

// New return new instance of Implementation.
func New(provider Provider) *Implementation {
	return &Implementation{
		provider: provider,
	}
}

// RegisterGRPC register gRPC service.
func (i *Implementation) RegisterGRPC(server *grpc.Server) {
	desc.RegisterScreenshotServer(server, i)
}

// RegisterGateway register HTTP handlers.
func (i *Implementation) RegisterGateway(ctx context.Context, mux *runtime.ServeMux) error {
	return desc.RegisterScreenshotHandlerServer(ctx, mux, i) // nolint:wrapcheck
}

// RegisterGatewayFromEndpoint register HTTP handlers for specific endpoint.
func (i *Implementation) RegisterGatewayFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return desc.RegisterScreenshotHandlerFromEndpoint(ctx, mux, endpoint, opts) // nolint:wrapcheck
}
