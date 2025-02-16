FROM cgr.dev/chainguard/go AS builder

COPY . /app
WORKDIR /app

RUN go get ./...
RUN go generate ./...
RUN go test ./...

RUN CGO_ENABLED=0 GOARCH=${GOARCH} GOOS=${GOOS} go build -ldflags '-extldflags "-static"' ./cmd/buckettest

FROM cgr.dev/chainguard/static

WORKDIR /app
COPY --from=builder /app/buckettest .

EXPOSE 8080

ENTRYPOINT ["./buckettest"]
