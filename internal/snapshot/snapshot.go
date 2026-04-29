package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time capture of drift detection results.
type Snapshot struct {
	ID        string                 `json:"id"`
	CreatedAt time.Time              `json:"created_at"`
	StateFile string                 `json:"state_file"`
	Results   map[string]interface{} `json:"results"`
}

// Manager handles saving and loading snapshots to/from disk.
type Manager struct {
	dir string
}

// NewManager creates a new snapshot Manager writing to dir.
func NewManager(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("snapshot: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Save persists a snapshot to disk and returns the file path.
func (m *Manager) Save(s *Snapshot) (string, error) {
	if s.ID == "" {
		s.ID = fmt.Sprintf("%d", s.CreatedAt.UnixNano())
	}
	path := filepath.Join(m.dir, s.ID+".json")
	f, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(s); err != nil {
		return "", fmt.Errorf("snapshot: encode: %w", err)
	}
	return path, nil
}

// Load reads a snapshot by ID from disk.
func (m *Manager) Load(id string) (*Snapshot, error) {
	path := filepath.Join(m.dir, id+".json")
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: open %q: %w", id, err)
	}
	defer f.Close()
	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot: decode: %w", err)
	}
	return &s, nil
}

// List returns all snapshot IDs stored in the manager directory.
func (m *Manager) List() ([]string, error) {
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read dir: %w", err)
	}
	var ids []string
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".json" {
			ids = append(ids, e.Name()[:len(e.Name())-5])
		}
	}
	return ids, nil
}
