package screenshot

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	desc "shot/pkg/api/screenshot"
)

// Image Get site image
func (i *Implementation) Image(ctx context.Context, req *desc.ImageRequest) (*desc.ImageResponse, error) {
	data, err := i.provider.Image(ctx, req.GetUrl())
	if err != nil {
		slog.ErrorContext(ctx, "provider.Image", "err", err)
		return nil, status.Errorf(codes.Internal, "provider.Image: %s", err.Error())
	}

	return &desc.ImageResponse{
		Image: data,
	}, nil
}
