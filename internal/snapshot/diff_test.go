package snapshot_test

import (
	"testing"
	"time"

	"github.com/example/drift-check/internal/snapshot"
)

func makeSnap(results map[string]interface{}) *snapshot.Snapshot {
	return &snapshot.Snapshot{
		ID:        "snap",
		CreatedAt: time.Now(),
		Results:   results,
	}
}

func TestCompare_NoDiff(t *testing.T) {
	attrs := map[string]interface{}{"region": "us-east-1", "size": "t2.micro"}
	prev := makeSnap(map[string]interface{}{"res-1": attrs})
	curr := makeSnap(map[string]interface{}{"res-1": attrs})
	r, err := snapshot.Compare(prev, curr)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if len(r.Added)+len(r.Removed)+len(r.Changed) != 0 {
		t.Errorf("expected no diff, got %+v", r)
	}
}

func TestCompare_Added(t *testing.T) {
	prev := makeSnap(map[string]interface{}{})
	curr := makeSnap(map[string]interface{}{"res-new": map[string]interface{}{"k": "v"}})
	r, err := snapshot.Compare(prev, curr)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if len(r.Added) != 1 || r.Added[0] != "res-new" {
		t.Errorf("expected added res-new, got %v", r.Added)
	}
}

func TestCompare_Removed(t *testing.T) {
	prev := makeSnap(map[string]interface{}{"res-old": map[string]interface{}{"k": "v"}})
	curr := makeSnap(map[string]interface{}{})
	r, err := snapshot.Compare(prev, curr)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if len(r.Removed) != 1 || r.Removed[0] != "res-old" {
		t.Errorf("expected removed res-old, got %v", r.Removed)
	}
}

func TestCompare_Changed(t *testing.T) {
	prev := makeSnap(map[string]interface{}{"res-1": map[string]interface{}{"size": "t2.micro"}})
	curr := makeSnap(map[string]interface{}{"res-1": map[string]interface{}{"size": "t3.large"}})
	r, err := snapshot.Compare(prev, curr)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if len(r.Changed) != 1 {
		t.Fatalf("expected 1 changed resource, got %d", len(r.Changed))
	}
	if r.Changed[0].Changes[0].Previous != "t2.micro" {
		t.Errorf("unexpected previous value: %s", r.Changed[0].Changes[0].Previous)
	}
}

func TestCompare_NilInputs(t *testing.T) {
	_, err := snapshot.Compare(nil, nil)
	if err == nil {
		t.Error("expected error for nil snapshots")
	}
}

func TestCompare_NonMapResultsIgnored(t *testing.T) {
	prev := makeSnap(map[string]interface{}{"meta": "string-value"})
	curr := makeSnap(map[string]interface{}{"meta": "string-value"})
	r, err := snapshot.Compare(prev, curr)
	if err != nil {
		t.Fatalf("Compare: %v", err)
	}
	if len(r.Added)+len(r.Removed)+len(r.Changed) != 0 {
		t.Errorf("non-map results should be ignored, got %+v", r)
	}
}
