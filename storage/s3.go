package storage

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	localcfg "github.com/korzepadawid/qr-codes-analyzer/config"
)

var (
	ErrBucketDoesNotExist = errors.New("given bucket doesn't exist")
	ErrFailedToPutFile    = errors.New("failed to put file into bucket")
)

type AWSS3FileStorage struct {
	s3Client     *s3.Client
	addrCDN      string
	bucketName   string
	bucketRegion string
}

func NewAWSS3FileStorageService(c *localcfg.Config) *AWSS3FileStorage {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)

	if c.Env == localcfg.EnvDev || c.Env == localcfg.EnvTest {
		cfg, err = setUpTestOrDevEnv(ctx, cfg, err, c)
	}

	if err != nil {
		log.Fatalf("failed to load aws configuration, %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	err = bucketExists(ctx, s3Client, &c.AWSBucketName)

	if err != nil {
		if errors.Is(err, ErrBucketDoesNotExist) {
			bErr := createNewBucket(ctx, c, s3Client)
			if bErr != nil {
				log.Fatalf("failed to create a new bucket %v", bErr)
			}
		} else {
			log.Fatal(err)
		}
	}

	return &AWSS3FileStorage{
		s3Client:     s3Client,
		bucketName:   c.AWSBucketName,
		bucketRegion: c.AWSBucketRegion,
		addrCDN:      c.CDNAddress,
	}
}

func createNewBucket(ctx context.Context, c *localcfg.Config, s3Client *s3.Client) error {
	createBucketInput := s3.CreateBucketInput{
		Bucket:                    aws.String(c.AWSBucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{LocationConstraint: types.BucketLocationConstraint(c.AWSBucketRegion)},
	}

	_, err := s3Client.CreateBucket(ctx, &createBucketInput)

	return err
}

func bucketExists(ctx context.Context, s3Client *s3.Client, bucketName *string) error {
	out, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})

	if err != nil {
		return fmt.Errorf("failed to list aws buckets, %v", err)
	}

	for _, bucket := range out.Buckets {
		if bucketName == bucket.Name {
			return nil
		}
	}

	return ErrBucketDoesNotExist
}

func setUpTestOrDevEnv(ctx context.Context, cfg aws.Config, err error, lcfg *localcfg.Config) (aws.Config, error) {
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:               lcfg.LocalstackURL, // localstack
				HostnameImmutable: true,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err = config.LoadDefaultConfig(ctx, config.WithEndpointResolver(customResolver))
	return cfg, err
}
