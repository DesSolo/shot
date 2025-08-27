# Shot 📸
Website Screenshot Microservice built in Go.

Shot is a high-performance, self-hosted microservice that captures screenshots of websites. It provides a simple gRPC API for internal services and an HTTP gateway for easy web integration.

---

## 🌟 Features

- **Dual API Interface**: Native gRPC for performance and an HTTP/JSON gateway (REST) for easy consumption.
- **High-Quality Screenshots**: Uses a headless Chrome browser via `chromedp` to render pages accurately.
- **Simple & Efficient**: Single-purpose service focused on doing one thing well.
- **Robust Input Validation**: Validates URL input using protobuf rules (`protovalidate`).

## 🚀 Getting Started

### Prerequisites

- Go 1.21+ (if building from source)
- Docker and Docker Compose (recommended for easy deployment)
- A system capable of running a headless Chrome browser (the Docker image handles this).

### Installation & Running (Docker - Recommended)

The easiest way to run Shot is using the pre-built Docker image.

```bash
# Build the Docker image
docker build -t dessolo/shot .

# Run the service
docker run --rm -p 8080:8080 -p 8081:8081 -p 50051:50051 dessolo/shot
```

This command starts the service with the following ports:
- **`50051`**: gRPC server
- **`8080`**: HTTP gateway (REST API)
- **`8081`**: Swagger UI and pprof (profiling)

### Installation & Running (From Source)

```bash
# Clone the repository
git clone https://github.com/dessolo/shot.git
cd shot

# Install dependencies and generate code
make install-deps generate

# Run the service
make run
```

## 💡 Usage

Shot provides two ways to interact: a gRPC API and an HTTP/JSON gateway.

### Via HTTP Gateway (REST-like API)

You can use any HTTP client (e.g., `curl`, `wget`, Postman) to get screenshots.

**Example Request (Get an image):**
```bash
curl "http://localhost:8080/v1/image?url=https://github.com" | jq -r '.image' | base64 -d > image.png
```

This command will:
1. Send a request to the `/v1/image` endpoint with `https://github.com` as the URL.
2. The `jq -r '.image'` command extracts the base64 encoded image data from the JSON response.
3. `base64 -d > image.png` decodes the base64 string and saves it as `image.png`.

### Via gRPC API

For high-performance or internal service communication, use the gRPC API. The `.proto` file is located at `api/screenshot/screenshot.proto`.

**Example (Go Client):**
```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "shot/api/screenshot" // Replace with your actual proto import path
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("failed to dial: %v
", err)
		return
	}
	defer conn.Close()

	client := pb.NewScreenshotServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	req := &pb.ScreenshotRequest{
		Url:      "https://charm.sh",
	}

	res, err := client.GetScreenshot(ctx, req)
	if err != nil {
		fmt.Printf("could not get screenshot: %v
", err)
		return
	}

	if err := os.WriteFile("charm_screenshot.png", res.Image, 0644); err != nil {
		fmt.Printf("failed to write image: %v
", err)
		return
	}
	fmt.Println("Screenshot saved to charm_screenshot.png")
}
```

## 💻 Development

### Project Structure

```
.
├── api/                  # Protocol Buffer definitions for gRPC/HTTP API
├── cmd/                  # Main application entry point
├── docs/                 # API documentation (Swagger/OpenAPI)
├── examples/             # Configuration examples
├── internal/             # Internal packages and application logic
├── pkg/                  # Publicly usable packages (e.g., interceptors, server setup)
├── scripts/              # Utility scripts
├── easyp_vendor/         # Proto dependencies managed by easyp
├── Dockerfile            # Docker build configuration
├── Makefile              # Build, lint, generate, and run commands
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── README.md             # This file
└── VERSION               # Project version
```

### Building and Testing

Refer to `Makefile` for common development commands:

-   `make install-deps`: Install all necessary Go tools and linting binaries.
-   `make generate`: Generate Go code from `.proto` files.
-   `make lint`: Run `golangci-lint` for code quality checks.
-   `go test ./...`: Run all Go unit tests.
-   `go test -run TestName ./pkg/path`: Run a specific test.
