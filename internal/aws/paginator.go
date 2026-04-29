package aws

import (
	"context"
	"fmt"
)

// PageClient defines the interface for fetching a single page of live resources.
type PageClient interface {
	FetchPage(ctx context.Context, token string) ([]LiveResource, string, error)
}

// PaginatedFetcher iterates over all pages returned by a PageClient.
type PaginatedFetcher struct {
	client   PageClient
	maxPages int
}

// NewPaginatedFetcher creates a new PaginatedFetcher with the given client.
// maxPages defaults to 100 to prevent runaway pagination.
func NewPaginatedFetcher(client PageClient) *PaginatedFetcher {
	return &PaginatedFetcher{
		client:   client,
		maxPages: 100,
	}
}

// FetchAll retrieves all resources across all pages.
func (p *PaginatedFetcher) FetchAll(ctx context.Context) ([]LiveResource, error) {
	var all []LiveResource
	token := ""
	page := 0

	for {
		if page >= p.maxPages {
			return nil, fmt.Errorf("pagination exceeded max page limit of %d", p.maxPages)
		}

		resources, nextToken, err := p.client.FetchPage(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("fetching page %d: %w", page+1, err)
		}

		all = append(all, resources...)
		page++

		if nextToken == "" {
			break
		}
		token = nextToken
	}

	return all, nil
}

// FetchAllWithLimit retrieves resources up to a maximum count.
func (p *PaginatedFetcher) FetchAllWithLimit(ctx context.Context, limit int) ([]LiveResource, error) {
	var all []LiveResource
	token := ""
	page := 0

	for {
		if page >= p.maxPages {
			return nil, fmt.Errorf("pagination exceeded max page limit of %d", p.maxPages)
		}

		resources, nextToken, err := p.client.FetchPage(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("fetching page %d: %w", page+1, err)
		}

		all = append(all, resources...)
		page++

		if limit > 0 && len(all) >= limit {
			return all[:limit], nil
		}

		if nextToken == "" {
			break
		}
		token = nextToken
	}

	return all, nil
}
