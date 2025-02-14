package bucket

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type Bucket struct {
	bucket *storage.BucketHandle
}

func New(ctx context.Context, bucketName string) (*Bucket, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bucket := client.Bucket(bucketName)
	return &Bucket{bucket: bucket}, nil
}

func (a *Bucket) ListObjects(ctx context.Context) ([]string, error) {
	var objects []string
	it := a.bucket.Objects(ctx, nil)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		objects = append(objects, obj.Name)
	}
	return objects, nil
}

func (a *Bucket) GetObject(ctx context.Context, object string) ([]byte, error) {
	r, err := a.bucket.Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

func (a *Bucket) PutObject(ctx context.Context, object string, data []byte) error {
	w := a.bucket.Object(object).NewWriter(ctx)
	defer w.Close()
	_, err := w.Write(data)
	return err
}

func (a *Bucket) DeleteObject(ctx context.Context, object string) error {
	return a.bucket.Object(object).Delete(ctx)
}
