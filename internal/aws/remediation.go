package aws

import (
	"fmt"
	"strings"
)

// RemediationAction represents a suggested fix for a detected drift.
type RemediationAction struct {
	ResourceID   string
	ResourceType string
	Action       string
	Description  string
	TFCommand    string
}

// Remediator generates remediation suggestions for drifted resources.
type Remediator struct{}

// NewRemediator creates a new Remediator instance.
func NewRemediator() *Remediator {
	return &Remediator{}
}

// Suggest returns a list of remediation actions for the given drift results.
// Each drift result is expected to be a map with keys: resource_id, resource_type, attribute.
func (r *Remediator) Suggest(drifts []map[string]string) []RemediationAction {
	actions := make([]RemediationAction, 0, len(drifts))
	for _, d := range drifts {
		id := d["resource_id"]
		rType := d["resource_type"]
		attribute := d["attribute"]

		if id == "" || rType == "" {
			continue
		}

		action := r.buildAction(id, rType, attribute)
		actions = append(actions, action)
	}
	return actions
}

func (r *Remediator) buildAction(id, rType, attribute string) RemediationAction {
	resourceName := sanitizeResourceName(id)
	tfRef := fmt.Sprintf("%s.%s", rType, resourceName)

	var desc string
	if attribute != "" {
		desc = fmt.Sprintf("Attribute '%s' has drifted on %s '%s'", attribute, rType, id)
	} else {
		desc = fmt.Sprintf("Resource %s '%s' is missing from state", rType, id)
	}

	return RemediationAction{
		ResourceID:   id,
		ResourceType: rType,
		Action:       "terraform apply",
		Description:  desc,
		TFCommand:    fmt.Sprintf("terraform apply -target=%s", tfRef),
	}
}

func sanitizeResourceName(id string) string {
	replacer := strings.NewReplacer("-", "_", ":", "_", "/", "_")
	return replacer.Replace(id)
}
