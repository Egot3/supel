package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageService struct {
	client          *s3.Client
	presignedClient *s3.PresignClient
	bucket          string
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

func NewStorageService(client *s3.Client, presignedClient *s3.PresignClient, bucket string) *StorageService {
	return &StorageService{
		client:          client,
		presignedClient: presignedClient,
		bucket:          bucket,
	}
}

func (s *StorageService) GETurl(ctx context.Context, key string) (string, error) {
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

func (s *StorageService) PUTurl(ctx context.Context, key, mime string) (string, error) {
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
