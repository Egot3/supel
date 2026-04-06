package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (s *StorageService) EnsureBuckets(ctx context.Context, buckets []string) error {
	for _, bucket := range buckets {
		_, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})

		if err != nil {
			log.Println("Bucket already created: ", err)
			continue
		}

		var owned *types.BucketAlreadyOwnedByYou
		var exists *types.BucketAlreadyExists
		if errors.As(err, &owned) || errors.As(err, &exists) {
			log.Printf("Bucket %q already exists, skipping", bucket)
			continue
		}

		return fmt.Errorf("failed to create bucket %q: %w", bucket, err)
	}

	return nil
}
