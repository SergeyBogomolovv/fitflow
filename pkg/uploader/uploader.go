package uploader

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type uploader struct {
	manager *manager.Uploader
	client  *s3.Client
	bucket  string
}

func New(AccessKey, SecretKey, Region, Endpoint, Bucket string) *uploader {
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

	return &uploader{manager, client, Bucket}
}

func (u *uploader) Upload(ctx context.Context, key string, body io.Reader) (string, error) {
	result, err := u.manager.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
		Body:   body,
	})

	return result.Location, err
}

func (u *uploader) Delete(ctx context.Context, key string) error {
	_, err := u.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(key),
	})
	return err
}
