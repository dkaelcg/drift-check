package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

type mockTaggingClient struct {
	pages  []*resourcegroupstaggingapi.GetResourcesOutput
	callNo int
	err    error
}

func (m *mockTaggingClient) GetResources(_ context.Context, _ *resourcegroupstaggingapi.GetResourcesInput, _ ...func(*resourcegroupstaggingapi.Options)) (*resourcegroupstaggingapi.GetResourcesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	page := m.pages[m.callNo]
	m.callNo++
	return page, nil
}

func TestFetchByTags_ReturnsTaggedResources(t *testing.T) {
	client := &mockTaggingClient{
		pages: []*resourcegroupstaggingapi.GetResourcesOutput{
			{
				ResourceTagMappingList: []resourcegroupstaggingapi.ResourceTagMapping{
					{
						ResourceARN: aws.String("arn:aws:ec2:us-east-1:123456789012:instance/i-abc123"),
						Tags: []resourcegroupstaggingapi.Tag{
							{Key: aws.String("env"), Value: aws.String("prod")},
						},
					},
				},
			},
		},
	}

	tagger := NewTagger(client)
	results, err := tagger.FetchByTags(context.Background(), []TagFilter{{Key: "env", Values: []string{"prod"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Tags["env"] != "prod" {
		t.Errorf("expected tag env=prod, got %q", results[0].Tags["env"])
	}
}

func TestFetchByTags_Pagination(t *testing.T) {
	client := &mockTaggingClient{
		pages: []*resourcegroupstaggingapi.GetResourcesOutput{
			{
				PaginationToken: aws.String("token1"),
				ResourceTagMappingList: []resourcegroupstaggingapi.ResourceTagMapping{
					{ResourceARN: aws.String("arn:aws:s3:::bucket-a"), Tags: nil},
				},
			},
			{
				ResourceTagMappingList: []resourcegroupstaggingapi.ResourceTagMapping{
					{ResourceARN: aws.String("arn:aws:s3:::bucket-b"), Tags: nil},
				},
			},
		},
	}

	tagger := NewTagger(client)
	results, err := tagger.FetchByTags(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results across pages, got %d", len(results))
	}
}

func TestFetchByTags_APIError(t *testing.T) {
	client := &mockTaggingClient{err: errors.New("api failure")}
	tagger := NewTagger(client)
	_, err := tagger.FetchByTags(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetchByTags_NilARNSkipped(t *testing.T) {
	client := &mockTaggingClient{
		pages: []*resourcegroupstaggingapi.GetResourcesOutput{
			{
				ResourceTagMappingList: []resourcegroupstaggingapi.ResourceTagMapping{
					{ResourceARN: nil},
					{ResourceARN: aws.String("arn:aws:ec2:::valid")},
				},
			},
		},
	}

	tagger := NewTagger(client)
	results, err := tagger.FetchByTags(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result (nil ARN skipped), got %d", len(results))
	}
}
