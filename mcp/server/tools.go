package server

import (
	"common-tasks-mcp/pkg/task_manager/types"
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// registerTools registers all MCP tools with the server
func (s *Server) registerTools() {
	// List tasks tool
	s.mcp.AddTool(&mcp.Tool{
		Name:        "list_tasks",
		Description: "List all tasks or filter by tags",
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
		Description: "Get a specific task by ID with full details including prerequisites and downstream tasks",
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
		Description: "Create a new task",
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
		Description: "Update an existing task",
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
		Description: "Delete a task and clean up all references to it from other tasks",
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
	var args struct {
		Tags []string `json:"tags"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
				},
			},
		}, nil
	}

	var tasks []*types.Task

	if len(args.Tags) > 0 {
		// Get tasks for each tag and merge (union)
		taskMap := make(map[string]*types.Task)
		for _, tag := range args.Tags {
			tagTasks, _ := s.taskManager.GetTasksByTag(tag)
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
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to marshal tasks: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(data),
			},
		},
	}, nil
}

// handleGetTask handles the get_task tool
func (s *Server) handleGetTask(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
				},
			},
		}, nil
	}

	task, err := s.taskManager.GetTask(args.ID)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to get task: %v", err),
				},
			},
		}, nil
	}

	data, err := json.Marshal(task)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to marshal task: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(data),
			},
		},
	}, nil
}

// handleAddTask handles the add_task tool
func (s *Server) handleAddTask(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
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
	}

	if err := s.taskManager.AddTask(task); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to add task: %v", err),
				},
			},
		}, nil
	}

	data, err := json.Marshal(task)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to marshal task: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(data),
			},
		},
	}, nil
}

// handleUpdateTask handles the update_task tool
func (s *Server) handleUpdateTask(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
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
	}

	if err := s.taskManager.UpdateTask(task); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to update task: %v", err),
				},
			},
		}, nil
	}

	data, err := json.Marshal(task)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to marshal task: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(data),
			},
		},
	}, nil
}

// handleDeleteTask handles the delete_task tool
func (s *Server) handleDeleteTask(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to parse arguments: %v", err),
				},
			},
		}, nil
	}

	if err := s.taskManager.DeleteTask(args.ID); err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("failed to delete task: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Task %s deleted successfully", args.ID),
			},
		},
	}, nil
}
