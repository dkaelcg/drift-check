package aws

import (
	"fmt"
	"strings"
)

// CostEstimate holds a rough monthly cost estimate for a resource.
type CostEstimate struct {
	ResourceID   string
	ResourceType string
	Region       string
	MonthlyCost  float64
	Currency     string
	Note         string
}

// baseMonthlyCosts holds approximate on-demand monthly costs (USD) per resource type.
var baseMonthlyCosts = map[string]float64{
	"aws_instance":        72.00,
	"aws_db_instance":     50.00,
	"aws_s3_bucket":        2.30,
	"aws_lambda_function":  0.20,
	"aws_elb":             18.00,
	"aws_alb":             22.00,
	"aws_nat_gateway":     35.00,
	"aws_eip":              3.65,
	"aws_cloudfront_distribution": 10.00,
	"aws_rds_cluster":    120.00,
}

// CostEstimator estimates monthly cost for live resources.
type CostEstimator struct{}

// NewCostEstimator creates a new CostEstimator.
func NewCostEstimator() *CostEstimator {
	return &CostEstimator{}
}

// Estimate returns a CostEstimate for the given LiveResource.
func (e *CostEstimator) Estimate(r LiveResource) (CostEstimate, error) {
	if r.ID == "" {
		return CostEstimate{}, fmt.Errorf("resource ID must not be empty")
	}

	resType := strings.ToLower(r.Type)
	cost, ok := baseMonthlyCosts[resType]
	note := ""
	if !ok {
		cost = 0.0
		note = "no cost data available for this resource type"
	}

	return CostEstimate{
		ResourceID:   r.ID,
		ResourceType: r.Type,
		Region:       r.Region,
		MonthlyCost:  cost,
		Currency:     "USD",
		Note:         note,
	}, nil
}

// EstimateAll returns cost estimates for a slice of LiveResources.
func (e *CostEstimator) EstimateAll(resources []LiveResource) []CostEstimate {
	results := make([]CostEstimate, 0, len(resources))
	for _, r := range resources {
		est, err := e.Estimate(r)
		if err != nil {
			continue
		}
		results = append(results, est)
	}
	return results
}

// TotalMonthlyCost sums the monthly cost across all estimates.
func TotalMonthlyCost(estimates []CostEstimate) float64 {
	var total float64
	for _, est := range estimates {
		total += est.MonthlyCost
	}
	return total
}
