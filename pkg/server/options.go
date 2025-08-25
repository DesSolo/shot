package server

import (
	"time"

	"github.com/rs/cors"
	"google.golang.org/grpc"
)

const (
	defaultGRPCEndpoint      = ":50051"
	defaultHTTPEndpoint      = ":8080"
	defaultHTTPAdminEndpoint = ":8081"

	defaultHTTPReadHeaderTimeout = 2 * time.Second
	defaultHTTPShutdownTimeout   = 5 * time.Second
)

// OptionFunc ...
type OptionFunc func(*Options)

// Options ...
type Options struct {
	GRPCEndpoint          string
	HTTPEndpoint          string
	HTTPAdminEndpoint     string
	HTTPReadHeaderTimeout time.Duration
	HTTPShutdownTimeout   time.Duration

	GRPCServerOptions []grpc.ServerOption
	CorsOptions       cors.Options
}

// NewDefaultOptions ...
func NewDefaultOptions(options ...OptionFunc) Options {
	o := &Options{
		GRPCEndpoint:          defaultGRPCEndpoint,
		HTTPEndpoint:          defaultHTTPEndpoint,
		HTTPAdminEndpoint:     defaultHTTPAdminEndpoint,
		HTTPReadHeaderTimeout: defaultHTTPReadHeaderTimeout,
		HTTPShutdownTimeout:   defaultHTTPShutdownTimeout,
	}

	for _, option := range options {
		option(o)
	}

	return *o
}

// WithGRPCServerOptions ...
func WithGRPCServerOptions(options ...grpc.ServerOption) OptionFunc {
	return func(o *Options) {
		o.GRPCServerOptions = append(o.GRPCServerOptions, options...)
	}
}

// WithCorsOptions ...
func WithCorsOptions(options cors.Options) OptionFunc {
	return func(o *Options) {
		o.CorsOptions = options
	}
}
