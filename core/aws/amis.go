package aws

import (
	"context"
	"errors"
	"fmt"
	"sort"
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
	Tags    map[string]string
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
		return nil, fmt.Errorf("failed to create ec2 client: %w", err)
	}

	resp, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Owners: []string{
			"self",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to describe images: %w", err)
	}

	if len(resp.Images) == 0 {
		return nil, errors.New("no images found")
	}

	amis := make([]AMI, len(resp.Images))

	for index, image := range resp.Images {
		tags := make(map[string]string)
		for _, tag := range image.Tags {
			tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
		}

		amis[index] = AMI{
			Name:    aws.ToString(image.Name),
			ID:      aws.ToString(image.ImageId),
			Created: aws.ToString(image.CreationDate),
			Tags:    tags,
		}
	}

	sort.Slice(amis, func(i, j int) bool {
		return amis[i].Created > amis[j].Created
	})

	return amis, nil
}

func (m EC2) Delete(ctx context.Context, amiID string) error {
	client, err := m.ec2Client(ctx)
	if err != nil {
		return fmt.Errorf("failed to create ec2 client: %w", err)
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
		return fmt.Errorf("failed to describe image: %w", err)
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
		return fmt.Errorf("failed to describe snapshots: %w", err)
	}

	if _, err := client.DeregisterImage(ctx, &ec2.DeregisterImageInput{ImageId: imageID}); err != nil {
		return fmt.Errorf("failed to deregister image: %w", err)
	}

	for _, snapshot := range describeSnapshots.Snapshots {
		if _, err := client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{SnapshotId: snapshot.SnapshotId}); err != nil {
			return fmt.Errorf("failed to delete snapshot: %w", err)
		}
	}

	return nil
}

func (m EC2) ec2Client(ctx context.Context) (*ec2.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create EC2 client: %w", err)
	}

	return ec2.NewFromConfig(cfg), nil
}
