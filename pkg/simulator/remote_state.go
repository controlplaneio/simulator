package simulator

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CreateRemoteStateBucket initialises a remote-state bucket
func CreateRemoteStateBucket(logger *logrus.Logger, bucket string) error {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	if _, err := svc.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(bucket)}); err != nil {
		return errors.Wrapf(err, "Unable to create bucket %q", bucket)
	}

	logger.WithFields(logrus.Fields{
		"BucketName": bucket,
	}).Info("Waiting for bucket to be created")
	if err := svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return errors.Wrapf(err, "Error occurred while waiting for bucket to be created, %v", bucket)
	}

	logger.WithFields(logrus.Fields{
		"BucketName": bucket,
	}).Infof("Bucket successfully created")
	return nil
}
