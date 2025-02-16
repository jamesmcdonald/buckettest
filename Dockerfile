FROM cgr.dev/chainguard/go AS builder

COPY . /app
WORKDIR /app

RUN --mount=type=cache,target=/root/go/pkg/mod go mod download
RUN --mount=type=cache,target=/root/go/pkg/mod go generate ./...

RUN --mount=type=cache,target=/root/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOARCH=${GOARCH} GOOS=${GOOS} go build -ldflags '-extldflags "-static"' ./cmd/buckettest
RUN --mount=type=cache,target=/root/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go test ./...


FROM cgr.dev/chainguard/static

WORKDIR /app
COPY --from=builder /app/buckettest .

EXPOSE 8080

ENTRYPOINT ["./buckettest"]
