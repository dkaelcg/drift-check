package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/drift-check/internal/snapshot"
)

func newTempManager(t *testing.T) *snapshot.Manager {
	t.Helper()
	dir := t.TempDir()
	m, err := snapshot.NewManager(dir)
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	return m
}

func sampleSnapshot() *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        "test-snap-1",
		CreatedAt: time.Now().UTC(),
		StateFile: "terraform.tfstate",
		Results:   map[string]interface{}{"drift": true, "count": 3},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	m := newTempManager(t)
	s := sampleSnapshot()
	path, err := m.Save(s)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file at %s, got: %v", path, err)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	m := newTempManager(t)
	s := sampleSnapshot()
	if _, err := m.Save(s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := m.Load(s.ID)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.ID != s.ID {
		t.Errorf("ID mismatch: got %q want %q", loaded.ID, s.ID)
	}
	if loaded.StateFile != s.StateFile {
		t.Errorf("StateFile mismatch: got %q want %q", loaded.StateFile, s.StateFile)
	}
}

func TestLoad_NotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Load("nonexistent")
	if err == nil {
		t.Error("expected error for missing snapshot, got nil")
	}
}

func TestList_ReturnsIDs(t *testing.T) {
	m := newTempManager(t)
	for _, id := range []string{"snap-a", "snap-b", "snap-c"} {
		s := &snapshot.Snapshot{ID: id, CreatedAt: time.Now()}
		if _, err := m.Save(s); err != nil {
			t.Fatalf("Save %s: %v", id, err)
		}
	}
	ids, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(ids) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(ids))
	}
}

func TestList_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	// add a non-json file to ensure it is ignored
	_ = os.WriteFile(filepath.Join(dir, "ignore.txt"), []byte("x"), 0644)
	m, _ := snapshot.NewManager(dir)
	ids, err := m.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(ids) != 0 {
		t.Errorf("expected 0 snapshots, got %d", len(ids))
	}
}

func TestSave_AutoID(t *testing.T) {
	m := newTempManager(t)
	s := &snapshot.Snapshot{CreatedAt: time.Now()}
	path, err := m.Save(s)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if s.ID == "" {
		t.Error("expected auto-generated ID")
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not found: %v", err)
	}
}
