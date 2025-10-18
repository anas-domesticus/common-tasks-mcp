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
	// List tasks tool
	s.mcp.AddTool(&mcp.Tool{
		Name:        "list_tasks",
		Description: "Browse available tasks, optionally filtered by tags (e.g., 'backend', 'database', 'deployment'). Returns task summaries with ID, name, and a brief description. Use this to discover relevant workflows when starting work in a new area or looking for standard procedures. If you provide multiple tags, you'll get tasks that match any of them.",
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
		Name:        "get_task",
		Description: "Get the complete workflow for a specific task by its ID. Returns the full task description plus related tasks: what must be done first (prerequisites), what must follow (required), and what's recommended (suggested). Use this before starting any task to understand the complete workflow, not just the immediate action. This helps you avoid missing critical steps.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Task ID",
				},
			},
			"required": []string{"id"},
		},
	}, s.handleGetTask)

	// Add task tool
	s.mcp.AddTool(&mcp.Tool{
		Name:        "add_task",
		Description: "Create a new task with its complete workflow. Include what needs to happen before this task (prerequisites), what must happen after (required follow-ups), and what's recommended after (suggested follow-ups). Use this to document repeatable workflows so future work can follow the same process. The system ensures workflows stay consistent by preventing circular dependencies.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Unique task identifier",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Task name",
				},
				"summary": map[string]interface{}{
					"type":        "string",
					"description": "Brief summary of the task",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Detailed description of the task",
				},
				"tags": map[string]interface{}{
					"type":        "array",
					"items":       map[string]string{"type": "string"},
					"description": "Array of tags for categorization",
				},
				"prerequisiteIDs": map[string]interface{}{
					"type":        "array",
					"items":       map[string]string{"type": "string"},
					"description": "Array of prerequisite task IDs that must be completed first",
				},
				"downstreamRequiredIDs": map[string]interface{}{
					"type":        "array",
					"items":       map[string]string{"type": "string"},
					"description": "Array of required downstream task IDs that must follow",
				},
				"downstreamSuggestedIDs": map[string]interface{}{
					"type":        "array",
					"items":       map[string]string{"type": "string"},
					"description": "Array of suggested downstream task IDs",
				},
			},
			"required": []string{"id", "name"},
		},
	}, s.handleAddTask)

	// Update task tool
	s.mcp.AddTool(&mcp.Tool{
		Name:        "update_task",
		Description: "Modify an existing task's description or workflow relationships. Use this when a process changes and you need to update the documented workflow - for example, adding a new required step, removing an outdated prerequisite, or refining the task description. The task ID must already exist.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Task ID (must exist)",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Task name",
				},
				"summary": map[string]interface{}{
					"type":        "string",
					"description": "Brief summary of the task",
				},
				"description": map[string]interface{}{
					"type":        "string",
					"description": "Detailed description of the task",
				},
				"tags": map[string]interface{}{
					"type":        "array",
					"items":       map[string]string{"type": "string"},
					"description": "Array of tags for categorization",
				},
				"prerequisiteIDs": map[string]interface{}{
					"type":        "array",
					"items":       map[string]string{"type": "string"},
					"description": "Array of prerequisite task IDs that must be completed first",
				},
				"downstreamRequiredIDs": map[string]interface{}{
					"type":        "array",
					"items":       map[string]string{"type": "string"},
					"description": "Array of required downstream task IDs that must follow",
				},
				"downstreamSuggestedIDs": map[string]interface{}{
					"type":        "array",
					"items":       map[string]string{"type": "string"},
					"description": "Array of suggested downstream task IDs",
				},
			},
			"required": []string{"id", "name"},
		},
	}, s.handleUpdateTask)

	// Delete task tool
	s.mcp.AddTool(&mcp.Tool{
		Name:        "delete_task",
		Description: "Remove a task entirely. This automatically cleans up any references to this task in other tasks' workflows. Use this when a task is no longer relevant or has been superseded by a different workflow. This action cannot be undone.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "string",
					"description": "Task ID to delete",
				},
			},
			"required": []string{"id"},
		},
	}, s.handleDeleteTask)
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
