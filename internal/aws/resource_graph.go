package aws

import (
	"fmt"
	"strings"
)

// Edge represents a directional dependency between two resources.
type Edge struct {
	From string
	To   string
	Kind string
}

// ResourceGraph holds resources and their dependency relationships.
type ResourceGraph struct {
	nodes map[string]LiveResource
	edges []Edge
}

// NewResourceGraph builds a dependency graph from a slice of LiveResources.
func NewResourceGraph(resources []LiveResource) *ResourceGraph {
	g := &ResourceGraph{
		nodes: make(map[string]LiveResource),
	}
	for _, r := range resources {
		g.nodes[r.ID] = r
	}
	g.buildEdges(resources)
	return g
}

// buildEdges infers relationships from well-known attribute keys.
func (g *ResourceGraph) buildEdges(resources []LiveResource) {
	refKeys := []string{"vpc_id", "subnet_id", "security_group_id", "instance_id", "bucket"}
	for _, r := range resources {
		for _, key := range refKeys {
			val, ok := r.Attributes[key]
			if !ok || val == "" {
				continue
			}
			if _, exists := g.nodes[val]; exists {
				g.edges = append(g.edges, Edge{From: r.ID, To: val, Kind: key})
			}
		}
	}
}

// Neighbors returns all resource IDs that the given ID depends on.
func (g *ResourceGraph) Neighbors(id string) []string {
	var out []string
	for _, e := range g.edges {
		if e.From == id {
			out = append(out, e.To)
		}
	}
	return out
}

// Dot renders the graph in Graphviz DOT format.
func (g *ResourceGraph) Dot() string {
	var sb strings.Builder
	sb.WriteString("digraph drift {\n")
	for id := range g.nodes {
		label := fmt.Sprintf("%s\\n(%s)", id, g.nodes[id].Type)
		sb.WriteString(fmt.Sprintf("  %q [label=%q];\n", id, label))
	}
	for _, e := range g.edges {
		sb.WriteString(fmt.Sprintf("  %q -> %q [label=%q];\n", e.From, e.To, e.Kind))
	}
	sb.WriteString("}\n")
	return sb.String()
}

// NodeCount returns the number of nodes in the graph.
func (g *ResourceGraph) NodeCount() int { return len(g.nodes) }

// EdgeCount returns the number of edges in the graph.
func (g *ResourceGraph) EdgeCount() int { return len(g.edges) }
