package aws

import (
	"strings"
	"testing"
)

func sampleGraphResources() []LiveResource {
	return []LiveResource{
		{
			ID:   "vpc-123",
			Type: "aws_vpc",
			Region: "us-east-1",
			Attributes: map[string]string{},
		},
		{
			ID:   "subnet-456",
			Type: "aws_subnet",
			Region: "us-east-1",
			Attributes: map[string]string{"vpc_id": "vpc-123"},
		},
		{
			ID:   "sg-789",
			Type: "aws_security_group",
			Region: "us-east-1",
			Attributes: map[string]string{"vpc_id": "vpc-123"},
		},
	}
}

func TestNewResourceGraph_NodeCount(t *testing.T) {
	g := NewResourceGraph(sampleGraphResources())
	if g.NodeCount() != 3 {
		t.Errorf("expected 3 nodes, got %d", g.NodeCount())
	}
}

func TestNewResourceGraph_EdgeCount(t *testing.T) {
	g := NewResourceGraph(sampleGraphResources())
	if g.EdgeCount() != 2 {
		t.Errorf("expected 2 edges, got %d", g.EdgeCount())
	}
}

func TestNeighbors_ReturnsLinkedIDs(t *testing.T) {
	g := NewResourceGraph(sampleGraphResources())
	neighbors := g.Neighbors("subnet-456")
	if len(neighbors) != 1 || neighbors[0] != "vpc-123" {
		t.Errorf("expected [vpc-123], got %v", neighbors)
	}
}

func TestNeighbors_NoEdges(t *testing.T) {
	g := NewResourceGraph(sampleGraphResources())
	neighbors := g.Neighbors("vpc-123")
	if len(neighbors) != 0 {
		t.Errorf("expected no neighbors for vpc-123, got %v", neighbors)
	}
}

func TestDot_ContainsNodes(t *testing.T) {
	g := NewResourceGraph(sampleGraphResources())
	dot := g.Dot()
	if !strings.Contains(dot, "digraph drift") {
		t.Error("DOT output missing digraph header")
	}
	if !strings.Contains(dot, "vpc-123") {
		t.Error("DOT output missing node vpc-123")
	}
}

func TestDot_ContainsEdges(t *testing.T) {
	g := NewResourceGraph(sampleGraphResources())
	dot := g.Dot()
	if !strings.Contains(dot, "->") {
		t.Error("DOT output missing edges")
	}
	if !strings.Contains(dot, "vpc_id") {
		t.Error("DOT output missing edge label vpc_id")
	}
}

func TestNewResourceGraph_Empty(t *testing.T) {
	g := NewResourceGraph([]LiveResource{})
	if g.NodeCount() != 0 || g.EdgeCount() != 0 {
		t.Error("expected empty graph for empty input")
	}
}
