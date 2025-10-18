package server

import (
	"common-tasks-mcp/pkg/graph_manager/types"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// registerTools registers all MCP tools with the server
func (s *Server) registerTools() {
	naming := s.config.MCP.Naming.Node

	// Read-only tools (always registered)

	// List tasks tool
	s.mcp.AddTool(&mcp.Tool{
		Name:        fmt.Sprintf("list_%s", naming.Plural),
		Description: fmt.Sprintf("Browse available %s, optionally filtered by tags (e.g., 'backend', 'database', 'deployment'). Returns %s summaries with ID, name, and a brief description. Use this to discover relevant workflows when starting work in a new area or looking for standard procedures. If you provide multiple tags, you'll get %s that match any of them.", naming.Plural, naming.Singular, naming.Plural),
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"tags": map[string]interface{}{
					"type":        "array",
					"items":       map[string]string{"type": "string"},
					"description": "Optional array of tags to filter by",
				},
			},
		},
	}, s.handleListTasks)

	// Get task tool
	s.mcp.AddTool(&mcp.Tool{
		Name:        fmt.Sprintf("get_%s", naming.Singular),
		Description: fmt.Sprintf("Get the complete workflow for a specific %s by its ID. Returns the full %s description plus related %s: what must be done first (prerequisites), what must follow (required), and what's recommended (suggested). Use this before starting any %s to understand the complete workflow, not just the immediate action. This helps you avoid missing critical steps.", naming.Singular, naming.Singular, naming.Plural, naming.Singular),
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": fmt.Sprintf("%s ID", naming.DisplaySingular),
				},
			},
			"required": []string{"id"},
		},
	}, s.handleGetTask)

	// List tags tool
	s.mcp.AddTool(&mcp.Tool{
		Name:        "list_tags",
		Description: fmt.Sprintf("Get all unique tags used across all %s. Returns each tag with the count of %s that use it. Use this to discover available tags for filtering and categorization.", naming.Plural, naming.Plural),
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}, s.handleListTags)

	// List prompts tool
	s.mcp.AddTool(&mcp.Tool{
		Name:        "list_prompts",
		Description: "Get all available prompts that can be used with this MCP server. Returns prompt names with their descriptions. Prompts are loaded from the prompts/ directory and can be customized per deployment.",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}, s.handleListPrompts)

	// Write tools (only registered when not in read-only mode)
	if !s.config.ReadOnly {
		// Add task tool
		s.mcp.AddTool(&mcp.Tool{
			Name:        fmt.Sprintf("add_%s", naming.Singular),
			Description: fmt.Sprintf("Create a new %s with its complete workflow. Include what needs to happen before this %s (prerequisites), what must happen after (required follow-ups), and what's recommended after (suggested follow-ups). Use this to document repeatable workflows so future work can follow the same process. The system ensures workflows stay consistent by preventing circular dependencies.", naming.Singular, naming.Singular),
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("Unique %s identifier", naming.Singular),
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("%s name", naming.DisplaySingular),
					},
					"summary": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("Brief summary of the %s", naming.Singular),
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("Detailed description of the %s", naming.Singular),
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"items":       map[string]string{"type": "string"},
						"description": "Array of tags for categorization",
					},
					"prerequisiteIDs": map[string]interface{}{
						"type":        "array",
						"items":       map[string]string{"type": "string"},
						"description": fmt.Sprintf("Array of prerequisite %s IDs that must be completed first", naming.Singular),
					},
					"downstreamRequiredIDs": map[string]interface{}{
						"type":        "array",
						"items":       map[string]string{"type": "string"},
						"description": fmt.Sprintf("Array of required downstream %s IDs that must follow", naming.Singular),
					},
					"downstreamSuggestedIDs": map[string]interface{}{
						"type":        "array",
						"items":       map[string]string{"type": "string"},
						"description": fmt.Sprintf("Array of suggested downstream %s IDs", naming.Singular),
					},
				},
				"required": []string{"id", "name"},
			},
		}, s.handleAddTask)

		// Update task tool
		s.mcp.AddTool(&mcp.Tool{
			Name:        fmt.Sprintf("update_%s", naming.Singular),
			Description: fmt.Sprintf("Modify an existing %s's description or workflow relationships. Use this when a process changes and you need to update the documented workflow - for example, adding a new required step, removing an outdated prerequisite, or refining the %s description. The %s ID must already exist.", naming.Singular, naming.Singular, naming.Singular),
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("%s ID (must exist)", naming.DisplaySingular),
					},
					"name": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("%s name", naming.DisplaySingular),
					},
					"summary": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("Brief summary of the %s", naming.Singular),
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("Detailed description of the %s", naming.Singular),
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"items":       map[string]string{"type": "string"},
						"description": "Array of tags for categorization",
					},
					"prerequisiteIDs": map[string]interface{}{
						"type":        "array",
						"items":       map[string]string{"type": "string"},
						"description": fmt.Sprintf("Array of prerequisite %s IDs that must be completed first", naming.Singular),
					},
					"downstreamRequiredIDs": map[string]interface{}{
						"type":        "array",
						"items":       map[string]string{"type": "string"},
						"description": fmt.Sprintf("Array of required downstream %s IDs that must follow", naming.Singular),
					},
					"downstreamSuggestedIDs": map[string]interface{}{
						"type":        "array",
						"items":       map[string]string{"type": "string"},
						"description": fmt.Sprintf("Array of suggested downstream %s IDs", naming.Singular),
					},
				},
				"required": []string{"id", "name"},
			},
		}, s.handleUpdateTask)

		// Delete task tool
		s.mcp.AddTool(&mcp.Tool{
			Name:        fmt.Sprintf("delete_%s", naming.Singular),
			Description: fmt.Sprintf("Remove a %s entirely. This automatically cleans up any references to this %s in other %s' workflows. Use this when a %s is no longer relevant or has been superseded by a different workflow. This action cannot be undone.", naming.Singular, naming.Singular, naming.Singular, naming.Singular),
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": fmt.Sprintf("%s ID to delete", naming.DisplaySingular),
					},
				},
				"required": []string{"id"},
			},
		}, s.handleDeleteTask)
	}
}

// handleListTasks handles the list_tasks tool
func (s *Server) handleListTasks(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Handling list_tasks request")

	var args struct {
		Tags []string `json:"tags"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		s.logger.Error("Failed to parse list_tasks arguments", zap.Error(err))
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Listing nodes", zap.Strings("tags", args.Tags))

	var nodes []*types.Node

	if len(args.Tags) > 0 {
		// Get nodes for each tag and merge (union)
		nodeMap := make(map[string]*types.Node)
		for _, tag := range args.Tags {
			tagNodes, _ := s.taskManager.GetNodesByTag(tag)
			s.logger.Debug("Retrieved nodes by tag", zap.String("tag", tag), zap.Int("count", len(tagNodes)))
			for _, node := range tagNodes {
				nodeMap[node.ID] = node
			}
		}
		// Convert map to slice
		for _, node := range nodeMap {
			nodes = append(nodes, node)
		}
	} else {
		nodes = s.taskManager.ListAllNodes()
		s.logger.Debug("Retrieved all nodes", zap.Int("count", len(nodes)))
	}

	s.logger.Info("Successfully listed nodes", zap.Int("node_count", len(nodes)))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: formatNodesAsMarkdown(nodes),
			},
		},
	}, nil
}

// handleGetTask handles the get_task tool
func (s *Server) handleGetTask(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Handling get_task request")

	var args struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		s.logger.Error("Failed to parse get_task arguments", zap.Error(err))
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Getting node", zap.String("node_id", args.ID))

	node, err := s.taskManager.GetNode(args.ID)
	if err != nil {
		s.logger.Error("Failed to get node", zap.String("node_id", args.ID), zap.Error(err))
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to get node: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Successfully retrieved node", zap.String("node_id", args.ID), zap.String("node_name", node.Name))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: formatNodeAsMarkdown(node, s.taskManager),
			},
		},
	}, nil
}

// handleListTags handles the list_tags tool
func (s *Server) handleListTags(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Handling list_tags request")

	// Get all tags with counts
	tags := s.taskManager.GetAllTags()

	s.logger.Info("Successfully retrieved tags", zap.Int("tag_count", len(tags)))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: formatTagsAsMarkdown(tags),
			},
		},
	}, nil
}

// handleListPrompts handles the list_prompts tool
func (s *Server) handleListPrompts(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Handling list_prompts request")

	// Get all loaded prompts
	prompts := s.prompts

	s.logger.Info("Successfully retrieved prompts", zap.Int("prompt_count", len(prompts)))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: formatPromptsAsMarkdown(prompts),
			},
		},
	}, nil
}

// handleAddTask handles the add_task tool
func (s *Server) handleAddTask(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Handling add_task request")

	var args struct {
		ID                     string   `json:"id"`
		Name                   string   `json:"name"`
		Summary                string   `json:"summary"`
		Description            string   `json:"description"`
		Tags                   []string `json:"tags"`
		PrerequisiteIDs        []string `json:"prerequisiteIDs"`
		DownstreamRequiredIDs  []string `json:"downstreamRequiredIDs"`
		DownstreamSuggestedIDs []string `json:"downstreamSuggestedIDs"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		s.logger.Error("Failed to parse add_task arguments", zap.Error(err))
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Adding new node",
		zap.String("node_id", args.ID),
		zap.String("node_name", args.Name),
		zap.Strings("tags", args.Tags),
	)

	now := time.Now()
	node := &types.Node{
		ID:          args.ID,
		Name:        args.Name,
		Summary:     args.Summary,
		Description: args.Description,
		Tags:        args.Tags,
		EdgeIDs: map[string][]string{
			"prerequisites":        args.PrerequisiteIDs,
			"downstream_required":  args.DownstreamRequiredIDs,
			"downstream_suggested": args.DownstreamSuggestedIDs,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.taskManager.AddNode(node); err != nil {
		s.logger.Error("Failed to add node",
			zap.String("node_id", args.ID),
			zap.String("node_name", args.Name),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to add node: %v", err),
				},
			},
		}, nil
	}

	// Persist changes to disk
	nodesPath := filepath.Join(s.config.Directory, "nodes")
	if err := s.taskManager.PersistToDir(nodesPath); err != nil {
		s.logger.Error("Failed to persist node to disk",
			zap.String("node_id", args.ID),
			zap.String("path", nodesPath),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("node added but failed to persist to disk: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Successfully added node", zap.String("node_id", args.ID), zap.String("node_name", args.Name))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("✓ Node `%s` created successfully\n\n%s", node.ID, formatNodeAsMarkdown(node, s.taskManager)),
			},
		},
	}, nil
}

// handleUpdateTask handles the update_task tool
func (s *Server) handleUpdateTask(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Handling update_task request")

	var args struct {
		ID                     string   `json:"id"`
		Name                   string   `json:"name"`
		Summary                string   `json:"summary"`
		Description            string   `json:"description"`
		Tags                   []string `json:"tags"`
		PrerequisiteIDs        []string `json:"prerequisiteIDs"`
		DownstreamRequiredIDs  []string `json:"downstreamRequiredIDs"`
		DownstreamSuggestedIDs []string `json:"downstreamSuggestedIDs"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		s.logger.Error("Failed to parse update_task arguments", zap.Error(err))
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Updating node",
		zap.String("node_id", args.ID),
		zap.String("node_name", args.Name),
		zap.Strings("tags", args.Tags),
	)

	// Get existing node to preserve CreatedAt timestamp
	existingNode, err := s.taskManager.GetNode(args.ID)
	if err != nil {
		s.logger.Error("Failed to get existing node for update",
			zap.String("node_id", args.ID),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to get existing node: %v", err),
				},
			},
		}, nil
	}

	node := &types.Node{
		ID:          args.ID,
		Name:        args.Name,
		Summary:     args.Summary,
		Description: args.Description,
		Tags:        args.Tags,
		EdgeIDs: map[string][]string{
			"prerequisites":        args.PrerequisiteIDs,
			"downstream_required":  args.DownstreamRequiredIDs,
			"downstream_suggested": args.DownstreamSuggestedIDs,
		},
		CreatedAt: existingNode.CreatedAt,
		UpdatedAt: time.Now(),
	}

	if err := s.taskManager.UpdateNode(node); err != nil {
		s.logger.Error("Failed to update node",
			zap.String("node_id", args.ID),
			zap.String("node_name", args.Name),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to update node: %v", err),
				},
			},
		}, nil
	}

	// Persist changes to disk
	nodesPath := filepath.Join(s.config.Directory, "nodes")
	if err := s.taskManager.PersistToDir(nodesPath); err != nil {
		s.logger.Error("Failed to persist node to disk",
			zap.String("node_id", args.ID),
			zap.String("path", nodesPath),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("node updated but failed to persist to disk: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Successfully updated node", zap.String("node_id", args.ID), zap.String("node_name", args.Name))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("✓ Node `%s` updated successfully\n\n%s", node.ID, formatNodeAsMarkdown(node, s.taskManager)),
			},
		},
	}, nil
}

// handleDeleteTask handles the delete_task tool
func (s *Server) handleDeleteTask(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Handling delete_task request")

	var args struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		s.logger.Error("Failed to parse delete_task arguments", zap.Error(err))
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Deleting node", zap.String("node_id", args.ID))

	if err := s.taskManager.DeleteNode(args.ID); err != nil {
		s.logger.Error("Failed to delete node", zap.String("node_id", args.ID), zap.Error(err))
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to delete node: %v", err),
				},
			},
		}, nil
	}

	// Persist changes to disk
	nodesPath := filepath.Join(s.config.Directory, "nodes")
	if err := s.taskManager.PersistToDir(nodesPath); err != nil {
		s.logger.Error("Failed to persist node deletion to disk",
			zap.String("node_id", args.ID),
			zap.String("path", nodesPath),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("node deleted but failed to persist to disk: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Successfully deleted node", zap.String("node_id", args.ID))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Node %s deleted successfully", args.ID),
			},
		},
	}, nil
}
