package simulator

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// CreateRemoteStateBucket initialises a remote-state bucket
func CreateRemoteStateBucket(logger *zap.SugaredLogger, bucket string) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	svc := s3.New(sess)

	_, err = svc.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucket)})
	if err != nil {
		return errors.Wrapf(err, "Unable to create bucket %q", bucket)
	}

	logger.Infof("Waiting for bucket %q to be created...", bucket)
	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return errors.Wrapf(err, "Error occurred while waiting for bucket to be created, %v", bucket)
	}

	logger.Infof("Bucket %q successfully created", bucket)
	return nil
}
