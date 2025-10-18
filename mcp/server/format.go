package server

import (
	"common-tasks-mcp/pkg/graph_manager"
	"common-tasks-mcp/pkg/graph_manager/types"
	"fmt"
	"sort"
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

	// Get all registered relationships from the manager
	allRelationships := tm.GetAllRelationships()

	// Group relationships by direction for this specific node
	backwardRels := []string{}
	forwardRels := []string{}
	noneRels := []string{}

	// Only include relationships that this node actually uses
	for relName := range node.EdgeIDs {
		rel, exists := allRelationships[relName]
		if !exists {
			// If relationship not registered, treat as DirectionNone
			noneRels = append(noneRels, relName)
			continue
		}

		switch rel.Direction {
		case types.DirectionBackward:
			backwardRels = append(backwardRels, relName)
		case types.DirectionForward:
			forwardRels = append(forwardRels, relName)
		case types.DirectionNone:
			noneRels = append(noneRels, relName)
		}
	}

	// Display backward relationships first (things that come before)
	for _, relName := range backwardRels {
		rel := allRelationships[relName]
		edgeIDs := node.GetEdgeIDs(relName)
		if len(edgeIDs) > 0 {
			sb.WriteString(fmt.Sprintf("**%s:**\n\n", capitalizeFirst(rel.Description)))
			for _, id := range edgeIDs {
				relatedNode, err := tm.GetNode(id)
				if err == nil {
					sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", id, relatedNode.Description))
				} else {
					sb.WriteString(fmt.Sprintf("`%s` (not found)\n\n", id))
				}
			}
		}
	}

	// Main node
	sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", node.ID, node.Description))

	// Display forward relationships (things that come after)
	for _, relName := range forwardRels {
		rel := allRelationships[relName]
		edgeIDs := node.GetEdgeIDs(relName)
		if len(edgeIDs) > 0 {
			sb.WriteString(fmt.Sprintf("**%s:**\n\n", capitalizeFirst(rel.Description)))
			for _, id := range edgeIDs {
				relatedNode, err := tm.GetNode(id)
				if err == nil {
					sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", id, relatedNode.Description))
				} else {
					sb.WriteString(fmt.Sprintf("`%s` (not found)\n\n", id))
				}
			}
		}
	}

	// Display relationships with no temporal direction
	for _, relName := range noneRels {
		edgeIDs := node.GetEdgeIDs(relName)
		if len(edgeIDs) > 0 {
			// Try to get description from registered relationship, otherwise use name
			label := relName
			if rel, exists := allRelationships[relName]; exists {
				label = rel.Description
			}
			sb.WriteString(fmt.Sprintf("**%s:**\n\n", capitalizeFirst(label)))
			for _, id := range edgeIDs {
				relatedNode, err := tm.GetNode(id)
				if err == nil {
					sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", id, relatedNode.Description))
				} else {
					sb.WriteString(fmt.Sprintf("`%s` (not found)\n\n", id))
				}
			}
		}
	}

	return strings.TrimSpace(sb.String())
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// formatTagsAsMarkdown formats a map of tags with their counts as markdown
func formatTagsAsMarkdown(tags map[string]int) string {
	if len(tags) == 0 {
		return "No tags found."
	}

	// Sort tags alphabetically for consistent output
	tagNames := make([]string, 0, len(tags))
	for tag := range tags {
		tagNames = append(tagNames, tag)
	}
	sort.Strings(tagNames)

	var sb strings.Builder
	for _, tag := range tagNames {
		count := tags[tag]
		sb.WriteString(fmt.Sprintf("%s (%d)\n", tag, count))
	}

	return strings.TrimSpace(sb.String())
}
