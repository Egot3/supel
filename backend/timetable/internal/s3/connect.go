package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
}

func NewClient(ctx context.Context, cfg Config) (*s3.Client, error) {
	if cfg.Region == "" {
		cfg.Region = "ru-west-1"
	}

	sdkConfig, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKey,
				cfg.SecretKey,
				"",
			)),
		config.WithRegion(cfg.Region),
		config.WithBaseEndpoint(cfg.Endpoint),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return client, nil
}
