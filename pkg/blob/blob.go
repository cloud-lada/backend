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
	// The Bucket type represents a blob storage bucket that blobs can be read or written to.
	Bucket struct {
		bucket *blob.Bucket
	}
)

// Open a connection with the blob storage bucket described in the url string.
func Open(ctx context.Context, url string) (*Bucket, error) {
	bucket, err := blob.OpenBucket(ctx, url)
	if err != nil {
		return nil, err
	}

	return &Bucket{bucket: bucket}, nil
}

// Close the connection to the blob storage bucket.
func (b *Bucket) Close() error {
	return b.bucket.Close()
}

// NewWriter returns an io.WriteCloser implementation that will write binary data as a blob within the Bucket under
// the specified name.
func (b *Bucket) NewWriter(ctx context.Context, name string) (io.WriteCloser, error) {
	writer, err := b.bucket.NewWriter(ctx, name, nil)
	if err != nil {
		return nil, err
	}

	return writer, nil
}
