package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate() error
}

// Validation ...
func Validation(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if validatable, ok := req.(validator); ok {
		if err := validatable.Validate(); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "Validation failed: %s", parseValidationErrors(err))
		}
	}

	return handler(ctx, req)
}

type validationFieldError interface {
	Field() string
	Reason() string
	Cause() error
}

func parseValidationErrors(err error) string {
	if validationFieldErr, ok := err.(validationFieldError); ok {
		if cause := validationFieldErr.Cause(); cause != nil {
			return parseValidationErrors(cause)
		}

		return validationFieldErr.Field() + " " + validationFieldErr.Reason()
	}

	return err.Error()
}
