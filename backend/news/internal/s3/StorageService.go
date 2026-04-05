package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type StorageService struct {
	client *s3.Client
	bucket string
}

func NewStorageService(client *s3.Client, bucket string) *StorageService {
	return &StorageService{
		client: client,
		bucket: bucket,
	}
}

func (s *StorageService) GETurl(ctx context.Context, key string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	presignedUrl, err := presignClient.PresignGetObject(ctx,
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

func (s *StorageService) PUTurl(ctx context.Context, key string) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	presignedUrl, err := presignClient.PresignPutObject(ctx,
		&s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		}, s3.WithPresignExpires(time.Minute*30))
	if err != nil {
		return "", err
	}

	return presignedUrl.URL, nil
}
