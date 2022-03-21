package s3

import (
	"bytes"
	"context"
	"fmt"

	"github.com/HomesNZ/go-common/s3/config"

	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type Service interface {
	Upload(ctx context.Context, key string, b []byte, expiry time.Time, contentType string) (url string, err error)
	Delete(ctx context.Context, key string) error
	Download(ctx context.Context, key string) ([]byte, error)
}

// S3 is a concrete implementation of cdn.Interface backed by S3 and Cloudfront.
type s3 struct {
	client *awsS3.Client
	config *config.Config
}

// UploadAsset uploads a new asset to S3 with the provided key. The URL returned will be the Cloudfront asset url ifcc
// S3.CloudfrontURL is not nil, otherwise a raw S3 URL is returned.
func (s s3) Upload(ctx context.Context, key string, b []byte, expiry time.Time, contentType string) (url string, err error) {
	reader := bytes.NewReader(b)

	params := &awsS3.PutObjectInput{
		Key:           aws.String(key),
		Bucket:        &s.config.BucketName,
		ACL:           s.config.ACL,
		Body:          reader,
		ContentLength: int64(reader.Len()),
		ContentType:   &contentType,
	}

	if !expiry.IsZero() {
		params.Expires = &expiry
	}

	uploader := manager.NewUploader(s.client)
	result, err := uploader.Upload(ctx, params)
	if err != nil {
		return "", errors.Wrap(err, "Failed to upload asset to aws S3 bucket")
	}

	return result.Location, nil
}

func (s s3) Delete(ctx context.Context, key string) error {
	params := &awsS3.DeleteObjectInput{
		Bucket: &s.config.BucketName,
		Key:    aws.String(key),
	}

	_, err := s.client.DeleteObject(ctx, params)
	if err != nil {
		return errors.Wrap(err, "Failed to delete asset from aws S3 bucket")
	}

	return nil
}

func (s s3) Download(ctx context.Context, key string) ([]byte, error) {
	file := &manager.WriteAtBuffer{}
	downloader := manager.NewDownloader(s.client)
	_, err := downloader.Download(ctx, file,
		&awsS3.GetObjectInput{
			Bucket: &s.config.BucketName,
			Key:    aws.String(key),
		})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to download asset from aws S3 bucket")
	}

	return file.Bytes(), nil
}

func (s s3) assetURL(key string) string {
	if s.config.CloudfrontURL != "" {
		return fmt.Sprintf("%s/%s", s.config.CloudfrontURL, key)
	}

	return fmt.Sprintf("%s/%s", s.config.Endpoint, key)
}
