PROJECT_NAME=shot
VERSION=$(shell cat VERSION)

tidy:
	go mod tidy

MAIN_FILE_PATH=cmd/main.go
CONFIG_FILE_PATH=examples/config.yml

run:
	CONFIG_FILE_PATH=${CONFIG_FILE_PATH} go run ${MAIN_FILE_PATH}

LOCAL_BIN=$(CURDIR)/bin

bin-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/easyp-tech/easyp/cmd/easyp@v0.7.15
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.27.1
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.2.1
	GOBIN=$(LOCAL_BIN) go install github.com/dessolo/protoc-gen-stub@latest


LINT_VERSION := 2.4.0

.install-lint:
	curl -Ls https://github.com/golangci/golangci-lint/releases/download/v${LINT_VERSION}/golangci-lint-${LINT_VERSION}-linux-amd64.tar.gz | tar xvz --strip-components=1 -C ${LOCAL_BIN} golangci-lint-${LINT_VERSION}-linux-amd64/golangci-lint

install-deps: \
	bin-deps \
	.install-lint

lint: $(LINT_BIN)
	$(LOCAL_BIN)/golangci-lint run

clean:
	rm -rf $(LOCAL_BIN)

generate:
	$(LOCAL_BIN)/easyp mod update
	PATH=$(LOCAL_BIN):$(PATH) $(LOCAL_BIN)/easyp generate
	rm pkg/*.json
	./scripts/fix_openapi.sh docs/apidocs.swagger.json "${PROJECT_NAME}" "${VERSION}"