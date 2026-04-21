package audit

import (
	"fmt"
	"sort"
	"strings"
)

// GraphNode represents a secret path node in a dependency graph.
type GraphNode struct {
	Path     string   `json:"path"`
	Env      string   `json:"env"`
	Drifted  bool     `json:"drifted"`
	DependsOn []string `json:"depends_on,omitempty"`
}

// GraphResult holds the full dependency graph output.
type GraphResult struct {
	Nodes []GraphNode            `json:"nodes"`
	Edges map[string][]string    `json:"edges"`
}

// BuildGraph constructs a dependency graph from scored reports and an
// optional prefix-based adjacency map (parent path -> child paths).
func BuildGraph(reports []ScoredReport, deps map[string][]string) GraphResult {
	nodes := make([]GraphNode, 0, len(reports))
	edges := make(map[string][]string)

	for _, r := range reports {
		node := GraphNode{
			Path:    r.Path,
			Env:     firstEnvFromScored(r),
			Drifted: r.RiskLevel != "low" && r.RiskLevel != "",
		}
		if children, ok := deps[r.Path]; ok {
			node.DependsOn = children
			edges[r.Path] = children
		}
		nodes = append(nodes, node)
	}

	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Path < nodes[j].Path
	})

	return GraphResult{Nodes: nodes, Edges: edges}
}

// PrintGraph renders the graph as a simple text tree to a strings.Builder.
func PrintGraph(g GraphResult) string {
	var sb strings.Builder
	visited := map[string]bool{}

	// Print root nodes (those not appearing as a child of anything)
	childSet := map[string]bool{}
	for _, children := range g.Edges {
		for _, c := range children {
			childSet[c] = true
		}
	}

	for _, node := range g.Nodes {
		if !childSet[node.Path] {
			printNode(&sb, g, node, 0, visited)
		}
	}
	return sb.String()
}

func printNode(sb *strings.Builder, g GraphResult, node GraphNode, depth int, visited map[string]bool) {
	if visited[node.Path] {
		return
	}
	visited[node.Path] = true
	indent := strings.Repeat("  ", depth)
	driftMark := ""
	if node.Drifted {
		driftMark = " [DRIFTED]"
	}
	fmt.Fprintf(sb, "%s- %s%s\n", indent, node.Path, driftMark)
	for _, childPath := range g.Edges[node.Path] {
		for _, n := range g.Nodes {
			if n.Path == childPath {
				printNode(sb, g, n, depth+1, visited)
				break
			}
		}
	}
}

func firstEnvFromScored(r ScoredReport) string {
	if len(r.Keys) > 0 {
		return ""
	}
	return ""
}
