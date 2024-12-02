package S3Client

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/s3/actions"
	"io"
)

type S3Client struct {
	*s3.Client
}

func NewS3Client() *S3Client {
	awsConf, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("eu-north-1"))
	if err != nil {
		panic("failed to load configuration, " + err.Error())
	}
	return &S3Client{
		Client: s3.NewFromConfig(awsConf),
	}
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

func (client *S3Client) GetPresignedUploadUrl(bucketName string, objectKey string, expiration int64) (string, error) {
	presignClient := s3.NewPresignClient(client.Client)
	presigner := actions.Presigner{PresignClient: presignClient}
	presignedPutRequest, err := presigner.PutObject(context.Background(), bucketName, objectKey, expiration)
	if err != nil {
		return "", err
	}
	return presignedPutRequest.URL, nil
}

func (client *S3Client) GetPresignedDownloadUrl(bucketName string, objectKey string, expiration int64) (string, error) {
	presignClient := s3.NewPresignClient(client.Client)
	presigner := actions.Presigner{PresignClient: presignClient}
	presignedGetRequest, err := presigner.GetObject(context.Background(), bucketName, objectKey, expiration)
	if err != nil {
		return "", err
	}
	return presignedGetRequest.URL, nil
}
