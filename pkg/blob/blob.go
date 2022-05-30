// Package blob provides types for interacting with blob storage providers.
package blob

import (
	"context"
	"io"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/memblob"
	_ "gocloud.dev/blob/s3blob"
)

type (
	Bucket struct {
		bucket *blob.Bucket
	}
)

func Open(ctx context.Context, url string) (*Bucket, error) {
	bucket, err := blob.OpenBucket(ctx, url)
	if err != nil {
		return nil, err
	}

	return &Bucket{bucket: bucket}, nil
}

func (b *Bucket) Close() error {
	return b.bucket.Close()
}

func (b *Bucket) NewWriter(ctx context.Context, name string) (io.WriteCloser, error) {
	writer, err := b.bucket.NewWriter(ctx, name, &blob.WriterOptions{})
	if err != nil {
		return nil, err
	}

	return writer, nil
}
