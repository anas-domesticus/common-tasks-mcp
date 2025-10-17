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
	for _, task := range tasks {
		sb.WriteString(fmt.Sprintf("%s - %s\n", task.ID, task.Summary))
	}

	return sb.String()
}

// formatTaskAsMarkdown formats a single task with full details as markdown
func formatTaskAsMarkdown(task *types.Task, tm *task_manager.Manager) string {
	var sb strings.Builder

	// Prerequisites first
	if len(task.PrerequisiteIDs) > 0 {
		sb.WriteString("**Prerequisites:**\n\n")
		for _, id := range task.PrerequisiteIDs {
			prereq, err := tm.GetTask(id)
			if err == nil {
				sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", id, prereq.Description))
			} else {
				sb.WriteString(fmt.Sprintf("`%s` (not found)\n\n", id))
			}
		}
	}

	// Main task
	sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", task.ID, task.Description))

	// Required downstream
	if len(task.DownstreamRequiredIDs) > 0 {
		sb.WriteString("**Required next:**\n\n")
		for _, id := range task.DownstreamRequiredIDs {
			next, err := tm.GetTask(id)
			if err == nil {
				sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", id, next.Description))
			} else {
				sb.WriteString(fmt.Sprintf("`%s` (not found)\n\n", id))
			}
		}
	}

	// Suggested downstream
	if len(task.DownstreamSuggestedIDs) > 0 {
		sb.WriteString("**Suggested next:**\n\n")
		for _, id := range task.DownstreamSuggestedIDs {
			next, err := tm.GetTask(id)
			if err == nil {
				sb.WriteString(fmt.Sprintf("`%s`\n\n%s\n\n", id, next.Description))
			} else {
				sb.WriteString(fmt.Sprintf("`%s` (not found)\n\n", id))
			}
		}
	}

	return strings.TrimSpace(sb.String())
}
