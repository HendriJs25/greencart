package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"user-service/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/transfermanager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Storage struct {
	client        *s3.Client
	uploader      *transfermanager.Client
	bucket        string
	publicBaseURL string
}

func NewS3Storage(ctx context.Context, cfg *config.AWS) (Storage, error) {
	opts := []func(*awsConfig.LoadOptions) error{
		awsConfig.WithRegion(cfg.Region),
	}

	if cfg.S3Endpoint != "" {
		if cfg.SecretAccessKey == "" || cfg.AccessKeyID == "" {
			return nil, fmt.Errorf("missing AWS credentials")
		}
		opts = append(opts, awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)))
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, opts...)

	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.S3Endpoint)
			o.UsePathStyle = true
		}
	})

	return &s3Storage{
		client:        client,
		uploader:      transfermanager.New(client),
		bucket:        cfg.S3Bucket,
		publicBaseURL: cfg.S3PublicBaseURL,
	}, nil

}

func (s *s3Storage) UploadFile(ctx context.Context, file *multipart.FileHeader, path string, contentType string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open source file %s: %w", file.Filename, err)
	}
	defer src.Close()

	input := &transfermanager.UploadObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        src,
		ContentType: aws.String(contentType),
	}

	result, err := s.uploader.UploadObject(ctx, input)

	if err != nil {
		return "", fmt.Errorf("failed to upload object %s: %w", file.Filename, err)
	}

	return *result.Key, nil
}

func (s *s3Storage) DeleteFile(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(strings.TrimPrefix(path, "/")),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", path, err)
	}
	return nil
}

func (s *s3Storage) GetFileURL(key string) *string {
	URL := fmt.Sprintf("%s/%s", s.publicBaseURL, key)
	return &URL
}
