package aws

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type AMIManager interface {
	List(ctx context.Context) ([]AMI, error)
	Delete(ctx context.Context, id string) error
}

type AMI struct {
	Name    string
	ID      string
	Created string
}

func (a AMI) CreationDate() string {
	if createTime, err := time.Parse(time.RFC3339, a.Created); err == nil {
		return createTime.Format(time.RFC822)
	}

	return a.Created
}

type EC2 struct{}

func (m EC2) List(ctx context.Context) ([]AMI, error) {
	client, err := m.ec2Client(ctx)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create ec2 client"), err)
	}

	resp, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Owners: []string{
			"self",
		},
	})
	if err != nil {
		return nil, errors.Join(errors.New("failed to list images"), err)
	}

	amis := make([]AMI, 0)
	for _, image := range resp.Images {
		amis = append(amis, AMI{
			Name:    aws.ToString(image.Name),
			ID:      aws.ToString(image.ImageId),
			Created: aws.ToString(image.CreationDate),
		})
	}

	return amis, nil
}

func (m EC2) Delete(ctx context.Context, amiID string) error {
	client, err := m.ec2Client(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to create ec2 client"), err)
	}

	describeImages, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		ImageIds: []string{
			amiID,
		},
		Owners: []string{
			"self",
		},
	})
	if err != nil {
		return errors.Join(errors.New("failed to describe image"), err)
	}

	if len(describeImages.Images) > 1 {
		return errors.New("too many images retrieved")
	}

	imageID := describeImages.Images[0].ImageId
	imageName := describeImages.Images[0].Name

	describeSnapshots, err := client.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		OwnerIds: []string{
			"self",
		},
		Filters: []types.Filter{
			{
				Name: aws.String("tag:AMI_Name"),
				Values: []string{
					aws.ToString(imageName),
				},
			},
		},
	})
	if err != nil {
		return errors.Join(errors.New("failed to describe snapshots"), err)
	}

	if _, err = client.DeregisterImage(ctx, &ec2.DeregisterImageInput{ImageId: imageID}); err != nil {
		return errors.Join(errors.New("failed to deregister image"), err)
	}

	for _, snapshot := range describeSnapshots.Snapshots {
		if _, err := client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{SnapshotId: snapshot.SnapshotId}); err != nil {
			return errors.Join(errors.New("failed to delete snapshot"), err)
		}
	}

	return nil
}

func (m EC2) ec2Client(ctx context.Context) (*ec2.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create client"), err)
	}

	return ec2.NewFromConfig(cfg), nil
}
