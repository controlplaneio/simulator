package aws

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type BucketCreator interface {
	Create(ctx context.Context, name string) error
}

func NewS3Client(ctx context.Context) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Join(errors.New("failed to load default config"), err)
	}

	return &S3Client{
		client: s3.NewFromConfig(cfg),
	}, nil
}

type S3Client struct {
	client *s3.Client
}

func (c S3Client) Create(ctx context.Context, name string) error {
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		return errors.New("failed to create bucket, aws region not set")
	}

	var bucketAlreadyOwnedByYou *types.BucketAlreadyOwnedByYou

	_, err := c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	if err != nil && !errors.As(err, &bucketAlreadyOwnedByYou) {
		return errors.Join(errors.New("failed to create bucket"), err)
	}

	return nil
}

func (c S3Client) Delete(ctx context.Context, name string) error {
	objects, err := c.client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return errors.Join(errors.New("failed to list bucket objects"), err)
	}

	objectIDs := make([]types.ObjectIdentifier, len(objects.Contents))
	for idx, object := range objects.Contents {
		objectIDs[idx] = types.ObjectIdentifier{Key: object.Key}
	}

	_, err = c.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(name),
		Delete: &types.Delete{
			Objects: objectIDs,
		},
	})
	if err != nil {
		return errors.Join(errors.New("failed to delete bucket objects"), err)
	}

	_, err = c.client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	})
	if err != nil {
		return errors.Join(errors.New("failed to delete bucket"), err)
	}

	return nil
}