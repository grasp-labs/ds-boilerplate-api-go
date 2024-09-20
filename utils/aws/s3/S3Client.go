package S3Client

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grasp-labs/ds-boilerplate-api-go/config"
	"io"
)

type S3Client struct {
	*s3.Client
	cfg *config.Config
}

func NewS3Client(cfg *config.Config) *S3Client {
	return &S3Client{cfg: cfg}
}

func (client *S3Client) ReadBytes(bucketName string, objectKey string) ([]byte, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}
	obj, err := client.GetObject(context.Background(), getObjectInput)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(obj.Body)
	return io.ReadAll(obj.Body)
}

func (client *S3Client) WriteBytes(bucketName string, objectKey string, data []byte) error {
	_, err := client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(data),
	})
	return err
}

func (client *S3Client) RemoveObject(bucketName string, objectKey string) error {
	_, err := client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	return err
}
