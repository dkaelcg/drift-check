package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// ScanOptions controls which resources are scanned.
type ScanOptions struct {
	Region      string
	Types       []string
	TagFilters  map[string]string
	MaxResults  int
}

// ScanResult holds the aggregated output of a region scan.
type ScanResult struct {
	Region    string
	Resources []LiveResource
	Errors    []string
}

// Scanner orchestrates fetching and enriching resources across types.
type Scanner struct {
	cfg      aws.Config
	fetcher  *PaginatedFetcher
	enricher *Enricher
	tagger   *Tagger
}

// NewScanner creates a Scanner backed by the provided AWS config.
func NewScanner(cfg aws.Config) *Scanner {
	return &Scanner{
		cfg:      cfg,
		fetcher:  NewPaginatedFetcher(cfg),
		enricher: NewEnricher(cfg),
		tagger:   NewTagger(cfg),
	}
}

// Scan fetches live resources according to opts and returns a ScanResult.
func (s *Scanner) Scan(ctx context.Context, opts ScanOptions) (*ScanResult, error) {
	if opts.Region == "" {
		return nil, fmt.Errorf("scanner: region must not be empty")
	}

	result := &ScanResult{Region: opts.Region}

	types := opts.Types
	if len(types) == 0 {
		types = SupportedResourceTypes()
	}

	for _, rt := range types {
		if !IsSupported(rt) {
			result.Errors = append(result.Errors, fmt.Sprintf("unsupported resource type: %s", rt))
			continue
		}

		resources, err := s.fetcher.FetchAll(ctx, rt)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("fetch %s: %v", rt, err))
			continue
		}

		for i := range resources {
			if err := s.enricher.Enrich(ctx, &resources[i]); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("enrich %s: %v", resources[i].ID, err))
			}
		}

		filtered := Apply(resources, Filter{
			Types:  []string{rt},
			Region: opts.Region,
			Tags:   opts.TagFilters,
		})

		result.Resources = append(result.Resources, filtered...)
	}

	if opts.MaxResults > 0 && len(result.Resources) > opts.MaxResults {
		result.Resources = result.Resources[:opts.MaxResults]
	}

	return result, nil
}
