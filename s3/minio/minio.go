package minio

import (
	"context"
	"sync"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tomyang1024/tools/s3"
)

var (
	maxTry int = 3
)

type Config struct {
	BucketName      string
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Secure          bool
	Location        string
}

var (
	_             s3.FileUploader = (*fileMinio)(nil)
	mu            sync.Mutex
	MinioUploader s3.FileUploader
)

type fileMinio struct {
	minioClient *minio.Client
	bucketName  string
}

func NewMinio(ctx context.Context, conf Config) (s3.FileUploader, error) {
	if MinioUploader == nil {
		mu.Lock()
		defer mu.Unlock()
		if MinioUploader == nil {
			client, err := minio.New(
				conf.Endpoint,
				&minio.Options{
					Creds:  credentials.NewStaticV4(conf.AccessKeyID, conf.SecretAccessKey, ""),
					Secure: conf.Secure,
				},
			)
			if err != nil {
				return nil, err
			}
			if err := client.MakeBucket(ctx, conf.BucketName,
				minio.MakeBucketOptions{Region: conf.Location}); err != nil {
				// 进行重试
				for i := 0; i < maxTry; i++ {
					if err := client.MakeBucket(ctx, conf.BucketName,
						minio.MakeBucketOptions{Region: conf.Location}); err == nil {
						break
					}
				}
				if _, err = client.BucketExists(ctx, conf.BucketName); err != nil {
					panic(err)
				}
			}
			MinioUploader = &fileMinio{
				minioClient: client,
				bucketName:  conf.BucketName,
			}

		}
	}
	return MinioUploader, nil
}

func (f *fileMinio) Engine() string {
	return "minio"
}

func (f *fileMinio) PutFromFile(ctx context.Context, name string, file string) error {
	_, err := f.minioClient.FPutObject(ctx, f.bucketName, name, file, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

func (f *fileMinio) DeleteObject(ctx context.Context, name string) error {
	return f.minioClient.RemoveObject(ctx, f.bucketName, name, minio.RemoveObjectOptions{})
}

// GetToFile 下载对象到本地
func (f *fileMinio) GetToFile(ctx context.Context, name string, file string) error {
	return f.minioClient.FGetObject(ctx, f.bucketName, name, file, minio.GetObjectOptions{})
}
