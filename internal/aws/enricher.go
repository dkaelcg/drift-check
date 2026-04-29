package aws

import (
	"fmt"
	"strings"
)

// ResourceTag represents a key-value tag on a cloud resource.
type ResourceTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// EnrichedResource wraps a LiveResource with additional metadata.
type EnrichedResource struct {
	LiveResource
	Tags       []ResourceTag     `json:"tags"`
	Region     string            `json:"region"`
	AccountID  string            `json:"account_id"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Enricher adds contextual metadata to fetched live resources.
type Enricher struct {
	region    string
	accountID string
}

// NewEnricher creates an Enricher with the given region and account ID.
func NewEnricher(region, accountID string) *Enricher {
	return &Enricher{
		region:    region,
		accountID: accountID,
	}
}

// Enrich takes a LiveResource and returns an EnrichedResource with
// region, account, and normalised tag metadata attached.
func (e *Enricher) Enrich(r LiveResource, rawTags map[string]string) (*EnrichedResource, error) {
	if r.ID == "" {
		return nil, fmt.Errorf("enricher: resource ID must not be empty")
	}

	tags := make([]ResourceTag, 0, len(rawTags))
	for k, v := range rawTags {
		tags = append(tags, ResourceTag{
			Key:   strings.TrimSpace(k),
			Value: strings.TrimSpace(v),
		})
	}

	meta := map[string]string{
		"source": "aws-live",
	}

	return &EnrichedResource{
		LiveResource: r,
		Tags:         tags,
		Region:       e.region,
		AccountID:    e.accountID,
		Metadata:     meta,
	}, nil
}

// TagMap converts the tag slice back to a plain map for downstream use.
func (er *EnrichedResource) TagMap() map[string]string {
	out := make(map[string]string, len(er.Tags))
	for _, t := range er.Tags {
		out[t.Key] = t.Value
	}
	return out
}
