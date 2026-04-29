package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// PageSize is the default number of results to request per API page.
const PageSize = 100

// PaginatedFetcher extends Fetcher with support for listing all resources of a
// given type across multiple API pages.
type PaginatedFetcher struct {
	fetcher *Fetcher
}

// NewPaginatedFetcher creates a PaginatedFetcher wrapping the provided Fetcher.
func NewPaginatedFetcher(f *Fetcher) *PaginatedFetcher {
	return &PaginatedFetcher{fetcher: f}
}

// ListAll returns all live resource IDs for the given resource type by
// iterating through every available API page.
func (p *PaginatedFetcher) ListAll(ctx context.Context, resourceType string) ([]string, error) {
	switch resourceType {
	case "aws_instance":
		return p.listEC2Instances(ctx)
	case "aws_s3_bucket":
		return p.listS3Buckets(ctx)
	default:
		return nil, fmt.Errorf("paginator: unsupported resource type %q", resourceType)
	}
}

// listEC2Instances pages through DescribeInstances and collects all instance IDs.
func (p *PaginatedFetcher) listEC2Instances(ctx context.Context) ([]string, error) {
	client, ok := p.fetcher.ec2Client.(*ec2.Client)
	if !ok {
		return nil, fmt.Errorf("paginator: ec2 client unavailable")
	}

	var ids []string
	paginator := ec2.NewDescribeInstancesPaginator(client, &ec2.DescribeInstancesInput{
		MaxResults: aws.Int32(int32(PageSize)),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("paginator: describing ec2 instances: %w", err)
		}
		for _, reservation := range page.Reservations {
			for _, instance := range reservation.Instances {
				if instance.InstanceId != nil {
					ids = append(ids, *instance.InstanceId)
				}
			}
		}
	}

	return ids, nil
}

// listS3Buckets lists all S3 buckets. The S3 ListBuckets API is not paginated,
// so all buckets are returned in a single call.
func (p *PaginatedFetcher) listS3Buckets(ctx context.Context) ([]string, error) {
	client, ok := p.fetcher.s3Client.(*s3.Client)
	if !ok {
		return nil, fmt.Errorf("paginator: s3 client unavailable")
	}

	out, err := client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("paginator: listing s3 buckets: %w", err)
	}

	ids := make([]string, 0, len(out.Buckets))
	for _, b := range out.Buckets {
		if b.Name != nil {
			ids = append(ids, *b.Name)
		}
	}

	return ids, nil
}
