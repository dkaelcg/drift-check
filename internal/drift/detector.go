package drift

import (
	"fmt"
	"sort"
)

// DriftResult holds the comparison result for a single resource.
type DriftResult struct {
	ResourceID   string
	ResourceType string
	HasDrift     bool
	Differences  []Difference
}

// Difference describes a single attribute mismatch.
type Difference struct {
	Attribute string
	StateValue string
	LiveValue  string
}

// StateResource represents a resource as defined in Terraform state.
type StateResource struct {
	ID         string
	Type       string
	Attributes map[string]string
}

// LiveResource represents a resource fetched from the cloud provider.
type LiveResource struct {
	ID         string
	Type       string
	Attributes map[string]string
}

// Detect compares a slice of state resources against their live counterparts.
// liveMap is keyed by resource ID.
func Detect(stateResources []StateResource, liveMap map[string]LiveResource) ([]DriftResult, error) {
	if stateResources == nil {
		return nil, fmt.Errorf("stateResources must not be nil")
	}

	results := make([]DriftResult, 0, len(stateResources))

	for _, sr := range stateResources {
		result := DriftResult{
			ResourceID:   sr.ID,
			ResourceType: sr.Type,
		}

		lr, found := liveMap[sr.ID]
		if !found {
			result.HasDrift = true
			result.Differences = append(result.Differences, Difference{
				Attribute:  "<resource>",
				StateValue: "exists",
				LiveValue:  "not found",
			})
			results = append(results, result)
			continue
		}

		// Collect all attribute keys from both sides.
		keys := mergeKeys(sr.Attributes, lr.Attributes)
		for _, key := range keys {
			sv := sr.Attributes[key]
			lv := lr.Attributes[key]
			if sv != lv {
				result.HasDrift = true
				result.Differences = append(result.Differences, Difference{
					Attribute:  key,
					StateValue: sv,
					LiveValue:  lv,
				})
			}
		}

		results = append(results, result)
	}

	return results, nil
}

func mergeKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
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
