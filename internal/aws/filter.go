package aws

import (
	"strings"
)

// ResourceFilter defines criteria for filtering live resources.
type ResourceFilter struct {
	Types  []string
	Region string
	Tags   map[string]string
}

// FilterResult holds resources that passed or were excluded by a filter.
type FilterResult struct {
	Included []LiveResource
	Excluded []LiveResource
}

// Apply filters a slice of LiveResource values based on the given ResourceFilter.
// A resource is included if it matches all specified criteria.
func Apply(resources []LiveResource, f ResourceFilter) FilterResult {
	result := FilterResult{}

	for _, r := range resources {
		if !matchesType(r, f.Types) {
			result.Excluded = append(result.Excluded, r)
			continue
		}
		if f.Region != "" && r.Region != f.Region {
			result.Excluded = append(result.Excluded, r)
			continue
		}
		if !matchesTags(r, f.Tags) {
			result.Excluded = append(result.Excluded, r)
			continue
		}
		result.Included = append(result.Included, r)
	}

	return result
}

func matchesType(r LiveResource, types []string) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if strings.EqualFold(r.Type, t) {
			return true
		}
	}
	return false
}

func matchesTags(r LiveResource, tags map[string]string) bool {
	if len(tags) == 0 {
		return true
	}
	for k, v := range tags {
		got, ok := r.Tags[k]
		if !ok || got != v {
			return false
		}
	}
	return true
}
