package uploader

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Uploader interface {
	Upload(ctx context.Context, key string, body io.Reader) (string, error)
	Delete(ctx context.Context, key string) error
}

type uploader struct {
	manager  *manager.Uploader
	client   *s3.Client
	bucket   string
	endpoint string
}

func MustNew(AccessKey, SecretKey, Region, Endpoint, Bucket string) Uploader {
	config, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AccessKey, SecretKey, "")),
		config.WithRegion(Region),
		config.WithBaseEndpoint(Endpoint),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	client := s3.NewFromConfig(config)
	manager := manager.NewUploader(client)

	endpoint, found := strings.CutPrefix(Endpoint, "https://")
	if !found {
		log.Fatalf("failed to parse endpoint: %s", Endpoint)
	}
	return &uploader{manager, client, Bucket, endpoint}
}

func (u *uploader) Upload(ctx context.Context, key string, body io.Reader) (string, error) {
	result, err := u.manager.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
		Body:   body,
	})

	return result.Location, err
}

func (u *uploader) Delete(ctx context.Context, url string) error {
	key, found := strings.CutPrefix(url, fmt.Sprintf("https://%s.%s", u.bucket, u.endpoint))
	if !found {
		return fmt.Errorf("failed to parse key: %s", url)
	}
	_, err := u.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
	})
	return err
}
