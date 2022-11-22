package s3

import (
	"context"
	"github.com/HomesNZ/go-common/s3/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type PresignService interface {
	SignedPutObjectUrl(ctx context.Context, bucket, key, contentType string) (string, error)
}

// presignS3 is a concrete implementation of presigned s3 client
type presignS3 struct {
	presignClient *awsS3.PresignClient
	presignConfig *config.PresignConfig
}

// SignedPutObjectUrl is to get a presigned URL for uploading an object
func (s presignS3) SignedPutObjectUrl(ctx context.Context, bucket, key, contentType string) (string, error) {
	params := &awsS3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	result, err := s.presignClient.PresignPutObject(ctx, params)
	if err != nil {
		return "", errors.Wrap(err, "PresignPutObject")
	}

	return result.URL, nil
}
