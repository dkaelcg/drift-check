package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a cached resource snapshot.
type Entry struct {
	ResourceID string            `json:"resource_id"`
	Attributes map[string]string `json:"attributes"`
	FetchedAt  time.Time         `json:"fetched_at"`
}

// Cache provides simple file-based caching for live resource snapshots.
type Cache struct {
	dir string
	ttl time.Duration
}

// New creates a Cache that stores entries under dir with the given TTL.
func New(dir string, ttl time.Duration) (*Cache, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cache: create dir: %w", err)
	}
	return &Cache{dir: dir, ttl: ttl}, nil
}

// Get retrieves a cached entry for the given resource ID.
// Returns (nil, nil) when the entry is absent or expired.
func (c *Cache) Get(resourceID string) (*Entry, error) {
	path := c.entryPath(resourceID)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("cache: read: %w", err)
	}
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("cache: unmarshal: %w", err)
	}
	if time.Since(entry.FetchedAt) > c.ttl {
		_ = os.Remove(path)
		return nil, nil
	}
	return &entry, nil
}

// Set writes an entry to the cache.
func (c *Cache) Set(entry Entry) error {
	entry.FetchedAt = time.Now()
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	if err := os.WriteFile(c.entryPath(entry.ResourceID), data, 0o644); err != nil {
		return fmt.Errorf("cache: write: %w", err)
	}
	return nil
}

// Invalidate removes a cached entry by resource ID.
func (c *Cache) Invalidate(resourceID string) error {
	err := os.Remove(c.entryPath(resourceID))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (c *Cache) entryPath(resourceID string) string {
	safe := filepath.Base(resourceID)
	return filepath.Join(c.dir, safe+".json")
}
