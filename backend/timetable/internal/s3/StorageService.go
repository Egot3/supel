package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/samber/do/v2"
)

type s3StorageService struct {
	client          *s3.Client
	presignedClient *s3.PresignClient
	bucket          string
}
type StorageService interface {
	GETurl(ctx context.Context, key string) (string, error)
	PUTurl(ctx context.Context, key, mime string) (string, error)
	EnsureBuckets(ctx context.Context, buckets []string) error
}

// func NewS3Client(endpoint, accessKey, secretKey string) (*s3.Client, error) {
// 	cfg, err := config.LoadDefaultConfig(context.Background(),
// 		config.WithRegion("ru-middle-1"),
// 		config.WithCredentialsProvider(
// 			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
// 		),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s3.NewFromConfig(cfg, func(o *s3.Options) {
// 		o.BaseEndpoint = aws.String(endpoint)
// 		o.UsePathStyle = true
// 	}), nil
// }

func NewStorageService(i do.Injector) (StorageService, error) {
	cfg, err := do.InvokeNamed[Config](i, "s3config.unsigned")
	if err != nil {
		return nil, err
	}
	client, err := do.Invoke[*s3.Client](i)
	if err != nil {
		return nil, err
	}
	presignedClient, err := do.Invoke[*s3.PresignClient](i)
	if err != nil {
		return nil, err
	}

	return &s3StorageService{
		client:          client,
		presignedClient: presignedClient,
		bucket:          cfg.Bucket,
	}, nil
}

func (s *s3StorageService) GETurl(ctx context.Context, key string) (string, error) {
	presignedUrl, err := s.presignedClient.PresignGetObject(ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(time.Minute*20),
	)
	if err != nil {
		return "", err
	}

	return presignedUrl.URL, nil
}

func (s *s3StorageService) PUTurl(ctx context.Context, key, mime string) (string, error) {
	presignedUrl, err := s.presignedClient.PresignPutObject(ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(s.bucket),
			Key:         aws.String(key),
			ContentType: aws.String(mime),
		}, s3.WithPresignExpires(time.Minute*30))
	if err != nil {
		return "", err
	}

	return presignedUrl.URL, nil
}
