package aws

import (
	"context"
	"errors"
	"testing"
)

type mockPaginatedClient struct {
	pages  [][]LiveResource
	callCount int
	err    error
}

func (m *mockPaginatedClient) FetchPage(ctx context.Context, token string) ([]LiveResource, string, error) {
	if m.err != nil {
		return nil, "", m.err
	}
	if m.callCount >= len(m.pages) {
		return nil, "", nil
	}
	page := m.pages[m.callCount]
	m.callCount++
	nextToken := ""
	if m.callCount < len(m.pages) {
		nextToken = "next-token"
	}
	return page, nextToken, nil
}

func TestPaginatedFetcher_SinglePage(t *testing.T) {
	client := &mockPaginatedClient{
		pages: [][]LiveResource{
			{{ID: "i-001", Type: "aws_instance", Attributes: map[string]string{"state": "running"}}},
		},
	}

	fetcher := NewPaginatedFetcher(client)
	results, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestPaginatedFetcher_MultiplePages(t *testing.T) {
	client := &mockPaginatedClient{
		pages: [][]LiveResource{
			{{ID: "i-001", Type: "aws_instance"}},
			{{ID: "i-002", Type: "aws_instance"}},
			{{ID: "i-003", Type: "aws_instance"}},
		},
	}

	fetcher := NewPaginatedFetcher(client)
	results, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
}

func TestPaginatedFetcher_APIError(t *testing.T) {
	client := &mockPaginatedClient{
		err: errors.New("api failure"),
	}

	fetcher := NewPaginatedFetcher(client)
	_, err := fetcher.FetchAll(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPaginatedFetcher_EmptyPages(t *testing.T) {
	client := &mockPaginatedClient{
		pages: [][]LiveResource{},
	}

	fetcher := NewPaginatedFetcher(client)
	results, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
