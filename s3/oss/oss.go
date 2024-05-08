package oss

import (
	"context"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/tomyang1024/tools/s3"
)

type Config struct {
	Endpoint        string
	Bucket          string
	BucketURL       string
	AccessKeyID     string
	AccessKeySecret string
	SessionToken    string
}

var (
	_           s3.FileUploader = (*fileOss)(nil)
	mu          sync.Mutex
	OssUploader s3.FileUploader
)

type fileOss struct {
	bucket *oss.Bucket
}

func NewOss(conf Config) (s3.FileUploader, error) {
	if OssUploader == nil {
		mu.Lock()
		defer mu.Unlock()
		if OssUploader == nil {
			client, err := oss.New(
				conf.Endpoint,
				conf.AccessKeyID,
				conf.AccessKeySecret,
				oss.SecurityToken(conf.SessionToken),
			)
			if err != nil {
				return nil, err
			}
			// 判断bucket是否存在，不存在则创建
			if exists, err := client.IsBucketExist(conf.Bucket); err == nil && !exists {
				if err = client.CreateBucket(conf.Bucket); err != nil {
					return nil, err
				}
			}
			bucket, err := client.Bucket(conf.Bucket)
			if err != nil {
				return nil, err
			}
			OssUploader = &fileOss{
				bucket: bucket,
			}
		}
	}
	return OssUploader, nil
}

// Engine returns the engine name.
func (f *fileOss) Engine() string {
	return "oss"
}

// PutFromFile uploads an object from a file.
func (f *fileOss) PutFromFile(ctx context.Context, name string, file string) error {
	return f.bucket.PutObjectFromFile(name, file)
}

// DeleteObject deletes an object.
func (f *fileOss) DeleteObject(ctx context.Context, name string) error {
	return f.bucket.DeleteObject(name)
}

// GetToFile downloads an object to a file.
func (f *fileOss) GetToFile(ctx context.Context, name string, file string) error {
	return f.bucket.GetObjectToFile(name, file)
}
