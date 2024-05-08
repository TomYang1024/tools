package s3

import "context"

type MultipartUploadResult struct {
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	UploadID string `json:"uploadID"`
}

type CompleteMultipartUploadResult struct {
	Location string `json:"location"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	ETag     string `json:"etag"`
}

type Part struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"etag"`
}

type FileUploader interface {
	Engine() string
	// PutFromFile uploads an object from a file.
	PutFromFile(ctx context.Context, name string, file string) error
	DeleteObject(ctx context.Context, name string) error
	GetToFile(ctx context.Context, name string, file string) error
}
