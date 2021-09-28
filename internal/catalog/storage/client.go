package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/c0llinn/ebook-store/config/log"
	"io"
	"time"
)

type Bucket string

type S3Client struct {
	service *s3.S3
	bucket  Bucket
}

func NewS3Client(service *s3.S3, bucket Bucket) S3Client {
	return S3Client{service: service, bucket: bucket}
}

func (c S3Client) GeneratePreSignedUrl(key string) (string, error) {
	request, _ := c.service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(string(c.bucket)),
		Key:    aws.String(key),
	})

	url, err := request.Presign(time.Hour)
	if err != nil {
		log.Logger.Errorf("Error generating get presignUrl for key %s: %v", key, err)
	}

	return url, err
}

func (c S3Client) SaveFile(key string, content io.ReadSeeker) error {
	_, err := c.service.PutObject(&s3.PutObjectInput{
		Key: aws.String(key),
		Bucket: aws.String(string(c.bucket)),
		Body: content,
	})

	if err != nil {
		return err
	}

	return nil
}