package s3

import (
	"context"
	"io"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/vasconcellos/financial-control/internal/domain/port"
)

type Storage struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
}

var _ port.ObjectStorage = (*Storage)(nil)

func NewStorage(cfg aws.Config, bucket string) *Storage {
	client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(client)
	return &Storage{
		client:    client,
		presigner: presigner,
		bucket:    bucket,
	}
}

func (s *Storage) Upload(ctx context.Context, key string, body io.Reader, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPrivate,
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	})
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *Storage) GetPresignedURL(ctx context.Context, key string) (string, error) {
	request, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		return "", err
	}
	return request.URL, nil
}
