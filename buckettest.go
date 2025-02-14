package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/jamesmcdonald/buckettest/internal/bucket"
	"github.com/jamesmcdonald/buckettest/internal/web"
)

// buckettest is a little program to test access to GCS buckets. It doesn't do
// any authentication so you can test that Workload Identity is providing the
// proper access.

func main() {
	bucketName := flag.String("bucket", os.Getenv("BUCKET"), "Bucket to use")
	flag.Parse()

	if *bucketName == "" {
		slog.Error("No bucket specified")
		os.Exit(1)
	}

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	bucket, err := bucket.New(ctx, *bucketName)
	if err != nil {
		logger.ErrorContext(ctx, "Error setting up bucket", "err", err)
		os.Exit(1)
	}

	logger.DebugContext(ctx, "Bucket setup", "bucket", *bucketName)

	logger.InfoContext(ctx, "Starting up", "bucket", *bucketName)

	webapp, err := web.New(bucket, logger)
	if err != nil {
		logger.ErrorContext(ctx, "Error setting up webapp", "err", err)
		os.Exit(1)
	}
	logger.InfoContext(ctx, "Webapp starting up")
	http.ListenAndServe(":8080", webapp)

}
