package snapshot

import (
	"fmt"
	"sort"
)

// FieldDiff describes a changed field between two snapshots.
type FieldDiff struct {
	Field    string `json:"field"`
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

// ResourceDiff holds the per-resource differences between two snapshots.
type ResourceDiff struct {
	ResourceID string      `json:"resource_id"`
	Changes    []FieldDiff `json:"changes"`
}

// DiffResult is the outcome of comparing two snapshots.
type DiffResult struct {
	Added   []string       `json:"added"`
	Removed []string       `json:"removed"`
	Changed []ResourceDiff `json:"changed"`
}

// Compare returns the differences between a previous and current snapshot.
// Results map keys are expected to be resource IDs mapping to attribute maps.
func Compare(prev, curr *Snapshot) (*DiffResult, error) {
	if prev == nil || curr == nil {
		return nil, fmt.Errorf("snapshot diff: both snapshots must be non-nil")
	}
	result := &DiffResult{}

	prevResources := flattenResults(prev.Results)
	currResources := flattenResults(curr.Results)

	for id := range currResources {
		if _, ok := prevResources[id]; !ok {
			result.Added = append(result.Added, id)
		}
	}
	for id := range prevResources {
		if _, ok := currResources[id]; !ok {
			result.Removed = append(result.Removed, id)
		}
	}
	for id, currAttrs := range currResources {
		prevAttrs, ok := prevResources[id]
		if !ok {
			continue
		}
		var diffs []FieldDiff
		keys := mergeStringKeys(prevAttrs, currAttrs)
		for _, k := range keys {
			p := fmt.Sprintf("%v", prevAttrs[k])
			c := fmt.Sprintf("%v", currAttrs[k])
			if p != c {
				diffs = append(diffs, FieldDiff{Field: k, Previous: p, Current: c})
			}
		}
		if len(diffs) > 0 {
			result.Changed = append(result.Changed, ResourceDiff{ResourceID: id, Changes: diffs})
		}
	}
	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	return result, nil
}

func flattenResults(r map[string]interface{}) map[string]map[string]interface{} {
	out := make(map[string]map[string]interface{})
	for k, v := range r {
		if attrs, ok := v.(map[string]interface{}); ok {
			out[k] = attrs
		}
	}
	return out
}

func mergeStringKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
