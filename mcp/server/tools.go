package server

import (
	"common-tasks-mcp/pkg/task_manager/types"
	"context"
	"encoding/json"
	"fmt"
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

	s.logger.Info("Listing tasks", zap.Strings("tags", args.Tags))

	var tasks []*types.Task

	if len(args.Tags) > 0 {
		// Get tasks for each tag and merge (union)
		taskMap := make(map[string]*types.Task)
		for _, tag := range args.Tags {
			tagTasks, _ := s.taskManager.GetTasksByTag(tag)
			s.logger.Debug("Retrieved tasks by tag", zap.String("tag", tag), zap.Int("count", len(tagTasks)))
			for _, task := range tagTasks {
				taskMap[task.ID] = task
			}
		}
		// Convert map to slice
		for _, task := range taskMap {
			tasks = append(tasks, task)
		}
	} else {
		tasks = s.taskManager.ListAllTasks()
		s.logger.Debug("Retrieved all tasks", zap.Int("count", len(tasks)))
	}

	s.logger.Info("Successfully listed tasks", zap.Int("task_count", len(tasks)))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: formatTasksAsMarkdown(tasks),
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

	s.logger.Info("Getting task", zap.String("task_id", args.ID))

	task, err := s.taskManager.GetTask(args.ID)
	if err != nil {
		s.logger.Error("Failed to get task", zap.String("task_id", args.ID), zap.Error(err))
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to get task: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Successfully retrieved task", zap.String("task_id", args.ID), zap.String("task_name", task.Name))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: formatTaskAsMarkdown(task, s.taskManager),
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

	s.logger.Info("Adding new task",
		zap.String("task_id", args.ID),
		zap.String("task_name", args.Name),
		zap.Strings("tags", args.Tags),
	)

	now := time.Now()
	task := &types.Task{
		ID:                     args.ID,
		Name:                   args.Name,
		Summary:                args.Summary,
		Description:            args.Description,
		Tags:                   args.Tags,
		PrerequisiteIDs:        args.PrerequisiteIDs,
		DownstreamRequiredIDs:  args.DownstreamRequiredIDs,
		DownstreamSuggestedIDs: args.DownstreamSuggestedIDs,
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	if err := s.taskManager.AddTask(task); err != nil {
		s.logger.Error("Failed to add task",
			zap.String("task_id", args.ID),
			zap.String("task_name", args.Name),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to add task: %v", err),
				},
			},
		}, nil
	}

	// Persist changes to disk
	if err := s.taskManager.PersistToDir(s.config.Directory); err != nil {
		s.logger.Error("Failed to persist task to disk",
			zap.String("task_id", args.ID),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("task added but failed to persist to disk: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Successfully added task", zap.String("task_id", args.ID), zap.String("task_name", args.Name))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("✓ Task `%s` created successfully\n\n%s", task.ID, formatTaskAsMarkdown(task, s.taskManager)),
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

	s.logger.Info("Updating task",
		zap.String("task_id", args.ID),
		zap.String("task_name", args.Name),
		zap.Strings("tags", args.Tags),
	)

	// Get existing task to preserve CreatedAt timestamp
	existingTask, err := s.taskManager.GetTask(args.ID)
	if err != nil {
		s.logger.Error("Failed to get existing task for update",
			zap.String("task_id", args.ID),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to get existing task: %v", err),
				},
			},
		}, nil
	}

	task := &types.Task{
		ID:                     args.ID,
		Name:                   args.Name,
		Summary:                args.Summary,
		Description:            args.Description,
		Tags:                   args.Tags,
		PrerequisiteIDs:        args.PrerequisiteIDs,
		DownstreamRequiredIDs:  args.DownstreamRequiredIDs,
		DownstreamSuggestedIDs: args.DownstreamSuggestedIDs,
		CreatedAt:              existingTask.CreatedAt,
		UpdatedAt:              time.Now(),
	}

	if err := s.taskManager.UpdateTask(task); err != nil {
		s.logger.Error("Failed to update task",
			zap.String("task_id", args.ID),
			zap.String("task_name", args.Name),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to update task: %v", err),
				},
			},
		}, nil
	}

	// Persist changes to disk
	if err := s.taskManager.PersistToDir(s.config.Directory); err != nil {
		s.logger.Error("Failed to persist task to disk",
			zap.String("task_id", args.ID),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("task updated but failed to persist to disk: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Successfully updated task", zap.String("task_id", args.ID), zap.String("task_name", args.Name))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("✓ Task `%s` updated successfully\n\n%s", task.ID, formatTaskAsMarkdown(task, s.taskManager)),
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

	s.logger.Info("Deleting task", zap.String("task_id", args.ID))

	if err := s.taskManager.DeleteTask(args.ID); err != nil {
		s.logger.Error("Failed to delete task", zap.String("task_id", args.ID), zap.Error(err))
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to delete task: %v", err),
				},
			},
		}, nil
	}

	// Persist changes to disk
	if err := s.taskManager.PersistToDir(s.config.Directory); err != nil {
		s.logger.Error("Failed to persist task deletion to disk",
			zap.String("task_id", args.ID),
			zap.Error(err),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("task deleted but failed to persist to disk: %v", err),
				},
			},
		}, nil
	}

	s.logger.Info("Successfully deleted task", zap.String("task_id", args.ID))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Task %s deleted successfully", args.ID),
			},
		},
	}, nil
}
