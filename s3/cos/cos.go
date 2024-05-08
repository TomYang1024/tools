package cos

import (
	"context"
	"net/http"
	"net/url"
	"sync"

	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"github.com/tomyang1024/tools/s3"
)

type Config struct {
	BucketURL    string
	SecretID     string
	SecretKey    string
	SessionToken string
	PublicRead   bool
	Debug        bool
}

var (
	mu          sync.Mutex
	_           s3.FileUploader = (*fileCos)(nil)
	CosUploader s3.FileUploader
)

type fileCos struct {
	publicRead bool
	copyURL    string
	client     *cos.Client
	credential *cos.Credential
	IsDebug    bool
}

// NewCos 单列模式
func NewCos(conf Config) (s3.FileUploader, error) {
	if CosUploader == nil {
		mu.Lock()
		defer mu.Unlock()
		if CosUploader == nil {
			u, err := url.Parse(conf.BucketURL)
			if err != nil {
				return nil, err
			}
			var (
				transport *debug.DebugRequestTransport
			)
			if conf.Debug {
				transport = &debug.DebugRequestTransport{
					RequestHeader:  true,
					RequestBody:    true,
					ResponseHeader: true,
					ResponseBody:   false,
				}
			}
			client := cos.NewClient(
				&cos.BaseURL{BucketURL: u},
				&http.Client{
					Transport: &cos.AuthorizationTransport{
						SecretID:     conf.SecretID,
						SecretKey:    conf.SecretKey,
						SessionToken: conf.SessionToken,
						Transport:    transport,
					},
				},
			)
			CosUploader = &fileCos{
				client: client,
			}
		}
	}

	return CosUploader, nil
}

func (c *fileCos) Engine() string {
	return "cos"
}

// PutFromFile uploads an object from a file.
func (c *fileCos) PutFromFile(ctx context.Context, name string, file string) error {
	_, err := c.client.Object.PutFromFile(ctx, name, file, nil)
	return err
}

// GetToFile downloads an object to a file.
func (c *fileCos) GetToFile(ctx context.Context, name string, file string) error {
	_, err := c.client.Object.GetToFile(ctx, name, file, nil)
	return err
}

// DeleteObject deletes an object.
func (c *fileCos) DeleteObject(ctx context.Context, name string) error {
	_, err := c.client.Object.Delete(ctx, name)
	return err
}
