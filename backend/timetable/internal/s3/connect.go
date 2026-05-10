package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/samber/do/v2"
)

func NewUnsignedClient(i do.Injector) (*s3.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cfg, err := do.InvokeNamed[Config](i, "s3config.unsigned")
	if err != nil {
		return nil, err
	}
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

func NewPresignedClient(i do.Injector) (*s3.PresignClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cfg, err := do.InvokeNamed[Config](i, "s3config.presigned")
	if err != nil {
		return nil, err
	}
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

	return s3.NewPresignClient(client), nil
}
