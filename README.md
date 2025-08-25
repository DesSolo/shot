# Shot
Website Screenshot Microservice

Shot is a high-performance, self-hosted microservice built in Go that captures screenshots of websites. It provides a simple gRPC API for internal services and an HTTP gateway for easy web integration.
## Features
- Dual API Interface: Native gRPC for performance and an HTTP/JSON gateway (REST) for easy consumption.
- High-Quality Screenshots: Uses a headless Chrome browser via chromedp to render pages accurately.
- Simple & Efficient: Single-purpose service focused on doing one thing well.
- Input Validation: Validates URL input using protobuf rules (protovalidate).

## Getting Started
Prerequisites

- Go 1.21+ (if building from source)
- Docker and Docker Compose (recommended for easy deployment)
- A system capable of running a headless Chrome browser (the Docker image handles this).

### Installation & Running
1. Using Docker (Recommended)

The easiest way to run Shot is using the pre-built Docker image.
bash

```shell
# build and run the image
docker build -t dessolo/shot
docker run --rm -p 8080:8080 -p 8081:8081 -p 50051:50051 dessolo/shot
```
This command starts the service:
- gRPC server on port 50051
- HTTP gateway on port 8080
- Swagger and pprof on port 8081

## Usage
Via HTTP Gateway (REST-like API)

You can use any HTTP client like curl, wget, or Postman.

Request:
```shell
curl "http://localhost:8080/v1/image?url=https://github.com" | jq -r '.image' | base64 -d > image.png
```

Response:
A successful response will return HTTP status `200 OK` and the PNG image data in the response body.
