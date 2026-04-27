package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// ResourceAttributes holds key-value pairs describing a live cloud resource.
type ResourceAttributes map[string]string

// LiveResource represents a resource fetched from AWS.
type LiveResource struct {
	Type       string
	ID         string
	Attributes ResourceAttributes
}

// Fetcher retrieves live AWS resource state.
type Fetcher struct {
	cfg aws.Config
}

// NewFetcher creates a Fetcher using the default AWS credential chain.
func NewFetcher(ctx context.Context, region string) (*Fetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("loading aws config: %w", err)
	}
	return &Fetcher{cfg: cfg}, nil
}

// FetchResource fetches live attributes for the given resource type and ID.
func (f *Fetcher) FetchResource(ctx context.Context, resourceType, resourceID string) (*LiveResource, error) {
	switch resourceType {
	case "aws_instance":
		return f.fetchEC2Instance(ctx, resourceID)
	case "aws_s3_bucket":
		return f.fetchS3Bucket(ctx, resourceID)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (f *Fetcher) fetchEC2Instance(ctx context.Context, instanceID string) (*LiveResource, error) {
	client := ec2.NewFromConfig(f.cfg)
	out, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("describe ec2 instance %s: %w", instanceID, err)
	}
	if len(out.Reservations) == 0 || len(out.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("ec2 instance %s not found", instanceID)
	}
	inst := out.Reservations[0].Instances[0]
	attrs := ResourceAttributes{
		"instance_type": string(inst.InstanceType),
		"instance_state": string(inst.State.Name),
	}
	if inst.ImageId != nil {
		attrs["ami"] = *inst.ImageId
	}
	return &LiveResource{Type: "aws_instance", ID: instanceID, Attributes: attrs}, nil
}

func (f *Fetcher) fetchS3Bucket(ctx context.Context, bucketName string) (*LiveResource, error) {
	client := s3.NewFromConfig(f.cfg)
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("head s3 bucket %s: %w", bucketName, err)
	}
	return &LiveResource{Type: "aws_s3_bucket", ID: bucketName, Attributes: ResourceAttributes{}}, nil
}
