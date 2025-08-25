package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Service ...
type Service interface {
	RegisterGRPC(*grpc.Server)
	RegisterGateway(context.Context, *runtime.ServeMux) error
	RegisterGatewayFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
}
