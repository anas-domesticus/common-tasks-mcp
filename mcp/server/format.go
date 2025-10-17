package server

import (
	"common-tasks-mcp/pkg/task_manager"
	"common-tasks-mcp/pkg/task_manager/types"
	"fmt"
	"strings"
)

// formatTasksAsMarkdown formats a list of tasks as markdown
func formatTasksAsMarkdown(tasks []*types.Task) string {
	if len(tasks) == 0 {
		return "No tasks found."
	}

	var sb strings.Builder
	sb.WriteString("## Available Tasks\n\n")

	for _, task := range tasks {
		sb.WriteString(fmt.Sprintf("### %s\n", task.ID))
		sb.WriteString(fmt.Sprintf("**%s**\n\n", task.Name))

		if task.Summary != "" {
			sb.WriteString(fmt.Sprintf("%s\n\n", task.Summary))
		}

		sb.WriteString("---\n\n")
	}

	return sb.String()
}

// formatTaskAsMarkdown formats a single task with full details as markdown
func formatTaskAsMarkdown(task *types.Task, tm *task_manager.Manager) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s\n\n", task.Name))
	sb.WriteString(fmt.Sprintf("**ID:** `%s`\n\n", task.ID))

	if task.Summary != "" {
		sb.WriteString(fmt.Sprintf("**Summary:** %s\n\n", task.Summary))
	}

	if task.Description != "" {
		sb.WriteString("## Description\n\n")
		sb.WriteString(fmt.Sprintf("%s\n\n", task.Description))
	}

	// Workflow section
	hasWorkflow := len(task.PrerequisiteIDs) > 0 || len(task.DownstreamRequiredIDs) > 0 || len(task.DownstreamSuggestedIDs) > 0
	if hasWorkflow {
		sb.WriteString("## Workflow\n\n")

		if len(task.PrerequisiteIDs) > 0 {
			sb.WriteString("**Prerequisites** (must be done first):\n")
			for _, id := range task.PrerequisiteIDs {
				prereq, err := tm.GetTask(id)
				if err == nil {
					sb.WriteString(fmt.Sprintf("- `%s`: %s\n", id, prereq.Name))
				} else {
					sb.WriteString(fmt.Sprintf("- `%s` (task not found)\n", id))
				}
			}
			sb.WriteString("\n")
		}

		if len(task.DownstreamRequiredIDs) > 0 {
			sb.WriteString("**Required Next Steps**:\n")
			for _, id := range task.DownstreamRequiredIDs {
				next, err := tm.GetTask(id)
				if err == nil {
					sb.WriteString(fmt.Sprintf("- `%s`: %s\n", id, next.Name))
				} else {
					sb.WriteString(fmt.Sprintf("- `%s` (task not found)\n", id))
				}
			}
			sb.WriteString("\n")
		}

		if len(task.DownstreamSuggestedIDs) > 0 {
			sb.WriteString("**Suggested Next Steps**:\n")
			for _, id := range task.DownstreamSuggestedIDs {
				next, err := tm.GetTask(id)
				if err == nil {
					sb.WriteString(fmt.Sprintf("- `%s`: %s\n", id, next.Name))
				} else {
					sb.WriteString(fmt.Sprintf("- `%s` (task not found)\n", id))
				}
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
