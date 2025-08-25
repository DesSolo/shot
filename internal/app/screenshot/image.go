package screenshot

import (
	"context"
	"fmt"
	"log/slog"

	"shot/pkg/api/screenshot"
)

// Image Get site image
func (i *Implementation) Image(ctx context.Context, req *screenshot.ImageRequest) (*screenshot.ImageResponse, error) {
	data, err := i.provider.Image(ctx, req.GetUrl())
	if err != nil {
		slog.ErrorContext(ctx, "provider.Image", "err", err)
		return nil, fmt.Errorf("provider.Image: %w", err)
	}

	return &screenshot.ImageResponse{
		Image: data,
	}, nil
}
