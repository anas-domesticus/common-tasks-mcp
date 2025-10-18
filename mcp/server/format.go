package server

import (
	"common-tasks-mcp/pkg/graph_manager"
	"common-tasks-mcp/pkg/graph_manager/types"
	"fmt"
	"strings"
)

// formatNodesAsMarkdown formats a list of nodes as markdown
func formatNodesAsMarkdown(nodes []*types.Node) string {
	if len(nodes) == 0 {
		return "No nodes found."
	}

	var sb strings.Builder
	for _, node := range nodes {
		sb.WriteString(fmt.Sprintf("%s - %s\n", node.ID, node.Summary))
	}

	return sb.String()
}

// formatNodeAsMarkdown formats a single node with full details as markdown
func formatNodeAsMarkdown(node *types.Node, tm *graph_manager.Manager) string {
	var sb strings.Builder

	// Get prerequisites from EdgeIDs
	prereqIDs := node.GetEdgeIDs("prerequisites")
	if len(prereqIDs) > 0 {
		sb.WriteString("**Prerequisites:**\n\n")
		for _, id := range prereqIDs {
			prereq, err := tm.GetNode(id)
			if err == nil {
				sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", id, prereq.Description))
			} else {
				sb.WriteString(fmt.Sprintf("`%s` (not found)\n\n", id))
			}
		}
	}

	// Main node
	sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", node.ID, node.Description))

	// Required downstream
	downstreamReqIDs := node.GetEdgeIDs("downstream_required")
	if len(downstreamReqIDs) > 0 {
		sb.WriteString("**Required next:**\n\n")
		for _, id := range downstreamReqIDs {
			next, err := tm.GetNode(id)
			if err == nil {
				sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", id, next.Description))
			} else {
				sb.WriteString(fmt.Sprintf("`%s` (not found)\n\n", id))
			}
		}
	}

	// Suggested downstream
	downstreamSugIDs := node.GetEdgeIDs("downstream_suggested")
	if len(downstreamSugIDs) > 0 {
		sb.WriteString("**Suggested next:**\n\n")
		for _, id := range downstreamSugIDs {
			next, err := tm.GetNode(id)
			if err == nil {
				sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", id, next.Description))
			} else {
				sb.WriteString(fmt.Sprintf("`%s` (not found)\n\n", id))
			}
		}
	}

	return strings.TrimSpace(sb.String())
}
