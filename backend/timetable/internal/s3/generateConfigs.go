package storage

import (
	"fmt"
	"os"

	"github.com/samber/do/v2"
)

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
}

func GenerateUnsignedS3Config(i do.Injector) (Config, error) {
	return Config{
		Endpoint: fmt.Sprintf("http://%v:%v", os.Getenv("STORAGE_HOST"),
			os.Getenv("STORAGE_PORT")),
		AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
		SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
		Bucket:    os.Getenv("STORAGE_TIMETABLE_BUCKET"),
	}, nil
}

func GeneratePresignedS3Config(i do.Injector) (Config, error) {
	return Config{
		Endpoint: fmt.Sprintf("http://%v:%v", os.Getenv("STORAGE_PUBLIC_HOST"),
			os.Getenv("STORAGE_PUBLIC_PORT")),
		AccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
		SecretKey: os.Getenv("STORAGE_SECRET_KEY"),
		Bucket:    os.Getenv("STORAGE_TIMETABLE_BUCKET"),
	}, nil
}
