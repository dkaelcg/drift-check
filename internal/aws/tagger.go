package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

// TagFilter represents a key/value filter for AWS resource tags.
type TagFilter struct {
	Key    string
	Values []string
}

// TaggedResource holds an ARN and its associated tags.
type TaggedResource struct {
	ARN  string
	Tags map[string]string
}

// TaggingClient is the interface for the AWS Resource Groups Tagging API.
type TaggingClient interface {
	GetResources(ctx context.Context, params *resourcegroupstaggingapi.GetResourcesInput, optFns ...func(*resourcegroupstaggingapi.Options)) (*resourcegroupstaggingapi.GetResourcesOutput, error)
}

// Tagger fetches resources by tag filters using the AWS tagging API.
type Tagger struct {
	client TaggingClient
}

// NewTagger creates a new Tagger with the provided tagging client.
func NewTagger(client TaggingClient) *Tagger {
	return &Tagger{client: client}
}

// FetchByTags returns all resources matching the given tag filters.
func (t *Tagger) FetchByTags(ctx context.Context, filters []TagFilter) ([]TaggedResource, error) {
	input := &resourcegroupstaggingapi.GetResourcesInput{}

	for _, f := range filters {
		input.TagFilters = append(input.TagFilters, resourcegroupstaggingapi.TagFilter{
			Key:    aws.String(f.Key),
			Values: f.Values,
		})
	}

	var results []TaggedResource

	for {
		resp, err := t.client.GetResources(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("tagger: GetResources failed: %w", err)
		}

		for _, rm := range resp.ResourceTagMappingList {
			if rm.ResourceARN == nil {
				continue
			}
			tr := TaggedResource{
				ARN:  *rm.ResourceARN,
				Tags: make(map[string]string, len(rm.Tags)),
			}
			for _, tag := range rm.Tags {
				if tag.Key != nil && tag.Value != nil {
					tr.Tags[*tag.Key] = *tag.Value
				}
			}
			results = append(results, tr)
		}

		if resp.PaginationToken == nil || *resp.PaginationToken == "" {
			break
		}
		input.PaginationToken = resp.PaginationToken
	}

	return results, nil
}
