FROM golang:1.24 AS builder

WORKDIR /build
ADD . .
RUN make build-docker && \
    cd bin/release && \
    mv *_docker shot

FROM alpine:3.15.0

RUN apk --update add chromium

WORKDIR /shot
COPY examples/config.yml config.yml
COPY docs/ docs/
COPY --from=builder /build/bin/release .
CMD ["./shot"]