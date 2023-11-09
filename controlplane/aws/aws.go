package aws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
	Env []string
)

func CreateBucket(ctx context.Context, name string) error {
	slog.Info("CreateBucket", "name", name)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		slog.Error("failed to create aws config", "error", err)
	}

	var bucketAlreadyOwnedByYou *types.BucketAlreadyOwnedByYou

	client := s3.NewFromConfig(cfg)
	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintEuWest2,
		},
	})
	if err != nil {
		if errors.As(err, &bucketAlreadyOwnedByYou) {
			slog.Info("You already own this bucket")
		} else {
			slog.Error("CreateBucket failed", "error", err)
			return err
		}
	}

	return nil
}

func init() {
	envKeys := []string{
		"AWS_PROFILE",
		"AWS_REGION",
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_SESSION_TOKEN",
	}

	for _, key := range envKeys {
		value, ok := os.LookupEnv(key)
		if ok && len(value) > 0 {
			Env = append(Env, fmt.Sprintf("%s=%s", key, value))
		}
	}
}
