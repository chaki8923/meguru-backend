
package r2

import (
	"context"
	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Service struct {
	s3Client   *s3.Client
	bucketName string
	bucketURL  string
}

func NewR2Service() (*R2Service, error) {
	r2Endpoint := os.Getenv("R2_ENDPOINT")
	r2AccessKey := os.Getenv("R2_ACCESS_KEY")
	r2SecretKey := os.Getenv("R2_SECRET_KEY")
	r2BucketName := "meguru" // Your bucket name
	r2BucketURL := os.Getenv("R2_BUCKET_URL")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r2AccessKey, r2SecretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: r2Endpoint}, nil
			},
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	return &R2Service{
		s3Client:   s3Client,
		bucketName: r2BucketName,
		bucketURL:  r2BucketURL,
	}, nil
}

func (s *R2Service) UploadImage(ctx context.Context, fileKey string, imageData []byte) (string, error) {
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucketName,
		Key:    &fileKey,
		Body:   bytes.NewReader(imageData),
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to R2: %w", err)
	}

	imageURL := fmt.Sprintf("%s/%s", s.bucketURL, fileKey)
	return imageURL, nil
}
