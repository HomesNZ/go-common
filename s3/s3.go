package s3

import (
	"bytes"
	"fmt"
	"github.com/HomesNZ/go-common/s3/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type Service interface {
	Upload(key string, b []byte, expiry time.Time, contentType string) (url string, err error)
	Delete(key string) (string, error)
	Download(key string) ([]byte, error)
}

// S3 is a concrete implementation of cdn.Interface backed by S3 and Cloudfront.
type s3 struct {
	client *awsS3.S3
	config *config.Config
}

// UploadAsset uploads a new asset to S3 with the provided key. The URL returned will be the Cloudfront asset url ifcc
// S3.CloudfrontURL is not nil, otherwise a raw S3 URL is returned.
func (s s3) Upload(key string, b []byte, expiry time.Time, contentType string) (url string, err error) {
	reader := bytes.NewReader(b)

	params := &awsS3.PutObjectInput{
		Bucket:        aws.String(s.config.BucketName),
		Key:           aws.String(key),
		ACL:           aws.String(s.config.ACL),
		Body:          reader,
		ContentLength: aws.Int64(int64(reader.Len())),
		ContentType:   &contentType,
	}

	if !expiry.IsZero() {
		params.SetExpires(expiry)
	}

	_, err = s.client.PutObject(params)
	if err != nil {
		return "", errors.Wrap(err, "Failed to upload asset to aws S3 bucket")
	}

	return s.assetURL(key), nil
}

func (s s3) Delete(key string) (string, error) {
	params := &awsS3.DeleteObjectInput{
		Bucket: aws.String(s.config.BucketName),
		Key:    aws.String(key),
	}

	resp, err := s.client.DeleteObject(params)
	if err != nil {
		logrus.Error(err)
		return resp.String(), errors.Wrap(err, "Failed to delete asset from aws S3 bucket")
	}

	return resp.String(), nil
}

func (s s3) Download(key string) ([]byte, error) {
	buff := &aws.WriteAtBuffer{}
	sess, _ := session.NewSession(&aws.Config{Region: aws.String(s.config.Region)})
	downloader := s3manager.NewDownloader(sess)
	_, err := downloader.Download(buff,
		&awsS3.GetObjectInput{
			Bucket: aws.String(s.config.BucketName),
			Key:    aws.String(key),
		})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to download asset from aws S3 bucket")
	}

	return buff.Bytes(), nil
}

func (s s3) assetURL(key string) string {
	if s.config.CloudfrontURL != "" {
		return fmt.Sprintf("%s/%s", s.config.CloudfrontURL, key)
	}

	return fmt.Sprintf("%s/%s", s.client.Endpoint, key)
}
