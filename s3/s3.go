package s3

import (
	"bytes"
	"fmt"
	"time"

	"github.com/HomesNZ/go-common/env"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

// ACL is policy for S3 assets.
// Region is the default region that assets are uploaded to.
// Endpoint is the default endpoint to be used when uploading assets.
// BucketName is aws S3 bucket Name
// CloudfrontURL is CDN url
type Config struct {
	BucketName      string
	ACL             string
	Region          string
	Endpoint        string
	CloudfrontURL   string
	AccessKeyID     string
	SecretAccessKey string
}

func ConfigFromEnv() Config {
	return Config{
		ACL:             env.GetString("AWS_S3_ACL", "private"),
		Region:          env.GetString("AWS_S3_REGION", "ap-southeast-2"),
		Endpoint:        env.GetString("AWS_S3_ENDPOINT", "s3-ap-southeast-2.amazonaws.com"),
		BucketName:      env.GetString("AWS_S3_BUCKET", ""),
		CloudfrontURL:   env.GetString("CDN_URL", ""),
		AccessKeyID:     env.GetString("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: env.GetString("AWS_SECRET_ACCESS_KEY", ""),
	}
}

// S3 is a concrete implementation of cdn.Interface backed by S3 and Cloudfront.
type S3 struct {
	*s3.S3
	Config *Config
}

// New initializes a new S3. If cloudfrontURL is not nil, URLs returned from UploadAsset will return the assets URL on
// Cloudfront distibution, otherwise the raw S3 URL will be returned.
func New(cfg *Config) (*S3, error) {
	if cfg.Region == "" {
		return nil, errors.New("aws S3 region was not specified")
	}
	if cfg.Endpoint == "" {
		return nil, errors.New("aws S3 endpoint was not specified")
	}
	if cfg.AccessKeyID == "" {
		return nil, errors.New("empty aws access key id")
	}
	if cfg.SecretAccessKey == "" {
		return nil, errors.New("empty aws secret access key")
	}

	awsConfig := &aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     cfg.AccessKeyID,
			SecretAccessKey: cfg.SecretAccessKey,
		}}),
	}

	return &S3{
		S3:     s3.New(session.New(), awsConfig),
		Config: cfg,
	}, nil
}

// UploadAsset uploads a new asset to S3 with the provided key. The URL returned will be the Cloudfront asset url ifcc
// S3.CloudfrontURL is not nil, otherwise a raw S3 URL is returned.
func (s S3) UploadAsset(key string, b []byte, expiry time.Time, contentType string) (url string, err error) {
	reader := bytes.NewReader(b)

	params := &s3.PutObjectInput{
		Bucket:        aws.String(s.Config.BucketName),
		Key:           aws.String(key),
		ACL:           aws.String(s.Config.ACL),
		Body:          reader,
		ContentLength: aws.Int64(int64(reader.Len())),
		ContentType:   &contentType,
	}

	if !expiry.IsZero() {
		params.SetExpires(expiry)
	}

	_, err = s.PutObject(params)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return s.assetURL(key), nil
}

func (s S3) DeleteAsset(key string) (string, error) {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(s.Config.BucketName),
		Key:    aws.String(key),
	}

	resp, err := s.DeleteObject(params)
	if err != nil {
		logrus.Error(err)
		return resp.String(), err
	}

	return resp.String(), nil
}

func (s S3) assetURL(key string) string {
	if s.Config.CloudfrontURL != "" {
		return fmt.Sprintf("%s/%s", s.Config.CloudfrontURL, key)
	}

	return fmt.Sprintf("%s/%s", s.Endpoint, key)
}
