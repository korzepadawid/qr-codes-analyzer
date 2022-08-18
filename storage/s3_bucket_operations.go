package storage

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (fs AWSS3FileStorage) PutFile(ctx context.Context, params PutFileParams) error {
	putObjectParams := s3.PutObjectInput{
		Bucket:      aws.String(fs.bucketName),
		Key:         aws.String(params.StorageKey),
		Body:        bytes.NewReader(*params.Object),
		ContentType: aws.String(params.ContentType),
	}

	_, err := fs.s3Client.PutObject(ctx, &putObjectParams, additionalOptions(fs))

	return err
}

func additionalOptions(fs AWSS3FileStorage) func(opt *s3.Options) {
	return func(opt *s3.Options) {
		opt.Region = fs.bucketRegion
	}
}

func (fs AWSS3FileStorage) DeleteFile(ctx context.Context, storageKey string) error {
	//TODO implement me
	panic("implement me")
}
