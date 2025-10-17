package task_manager

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"common-tasks-mcp/pkg/types"
)

func TestPersistAndLoad(t *testing.T) {
	// Create a temporary directory for test files
	testDir := filepath.Join(t.TempDir(), "tasks")

	// Create a manager with 3 tasks
	manager1 := NewManager()

	now := time.Now().UTC().Truncate(time.Second)

	// Task A - no dependencies
	taskA := &types.Task{
		ID:            "task-a",
		Name:          "Task A",
		Summary:       "First task",
		Description:   "This is the first task with no dependencies",
		Tags:          []string{"api", "backend"},
		DependencyIDs: []string{},
		DependentIDs:  []string{"task-b"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Task B - depends on A, has C as dependent
	taskB := &types.Task{
		ID:            "task-b",
		Name:          "Task B",
		Summary:       "Second task",
		Description:   "This task depends on A and has C as dependent",
		Tags:          []string{"frontend", "api"},
		DependencyIDs: []string{"task-a"},
		DependentIDs:  []string{"task-c"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Task C - depends on B
	taskC := &types.Task{
		ID:            "task-c",
		Name:          "Task C",
		Summary:       "Third task",
		Description:   "This task depends on B",
		Tags:          []string{"testing"},
		DependencyIDs: []string{"task-b"},
		DependentIDs:  []string{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	manager1.tasks["task-a"] = taskA
	manager1.tasks["task-b"] = taskB
	manager1.tasks["task-c"] = taskC

	// Persist tasks to directory
	if err := manager1.Persist(testDir); err != nil {
		t.Fatalf("Failed to persist tasks: %v", err)
	}

	// Verify files were created
	for _, id := range []string{"task-a", "task-b", "task-c"} {
		filename := filepath.Join(testDir, id+".yaml")
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", filename)
		}
	}

	// Create a new manager and load tasks
	manager2 := NewManager()
	if err := manager2.Load(testDir); err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}

	// Compare the managers
	if len(manager1.tasks) != len(manager2.tasks) {
		t.Errorf("Task count mismatch: expected %d, got %d", len(manager1.tasks), len(manager2.tasks))
	}

	// Compare each task
	for id, originalTask := range manager1.tasks {
		loadedTask, exists := manager2.tasks[id]
		if !exists {
			t.Errorf("Task %s not found in loaded manager", id)
			continue
		}

		if originalTask.ID != loadedTask.ID {
			t.Errorf("Task %s ID mismatch: expected %s, got %s", id, originalTask.ID, loadedTask.ID)
		}
		if originalTask.Name != loadedTask.Name {
			t.Errorf("Task %s Name mismatch: expected %s, got %s", id, originalTask.Name, loadedTask.Name)
		}
		if originalTask.Summary != loadedTask.Summary {
			t.Errorf("Task %s Summary mismatch: expected %s, got %s", id, originalTask.Summary, loadedTask.Summary)
		}
		if originalTask.Description != loadedTask.Description {
			t.Errorf("Task %s Description mismatch: expected %s, got %s", id, originalTask.Description, loadedTask.Description)
		}

		// Compare tags
		if len(originalTask.Tags) != len(loadedTask.Tags) {
			t.Errorf("Task %s Tags count mismatch: expected %d, got %d", id, len(originalTask.Tags), len(loadedTask.Tags))
		} else {
			for i, tag := range originalTask.Tags {
				if tag != loadedTask.Tags[i] {
					t.Errorf("Task %s Tag %d mismatch: expected %s, got %s", id, i, tag, loadedTask.Tags[i])
				}
			}
		}

		// Compare dependencies
		if len(originalTask.DependencyIDs) != len(loadedTask.DependencyIDs) {
			t.Errorf("Task %s Dependencies count mismatch: expected %d, got %d", id, len(originalTask.DependencyIDs), len(loadedTask.DependencyIDs))
		} else {
			for i, dep := range originalTask.DependencyIDs {
				if dep != loadedTask.DependencyIDs[i] {
					t.Errorf("Task %s Dependency %d mismatch: expected %s, got %s", id, i, dep, loadedTask.DependencyIDs[i])
				}
			}
		}

		// Compare dependents
		if len(originalTask.DependentIDs) != len(loadedTask.DependentIDs) {
			t.Errorf("Task %s Dependents count mismatch: expected %d, got %d", id, len(originalTask.DependentIDs), len(loadedTask.DependentIDs))
		} else {
			for i, dep := range originalTask.DependentIDs {
				if dep != loadedTask.DependentIDs[i] {
					t.Errorf("Task %s Dependent %d mismatch: expected %s, got %s", id, i, dep, loadedTask.DependentIDs[i])
				}
			}
		}

		// Compare timestamps (truncated to second precision due to YAML serialization)
		if !originalTask.CreatedAt.Equal(loadedTask.CreatedAt) {
			t.Errorf("Task %s CreatedAt mismatch: expected %v, got %v", id, originalTask.CreatedAt, loadedTask.CreatedAt)
		}
		if !originalTask.UpdatedAt.Equal(loadedTask.UpdatedAt) {
			t.Errorf("Task %s UpdatedAt mismatch: expected %v, got %v", id, originalTask.UpdatedAt, loadedTask.UpdatedAt)
		}
	}
}

func TestAddTask(t *testing.T) {
	manager := NewManager()
	now := time.Now().UTC().Truncate(time.Second)

	// Test adding a valid task
	task1 := &types.Task{
		ID:          "task-1",
		Name:        "Test Task 1",
		Summary:     "First test task",
		Description: "A valid task to add",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := manager.AddTask(task1)
	if err != nil {
		t.Fatalf("Expected no error when adding valid task, got: %v", err)
	}

	// Verify task was added
	if len(manager.tasks) != 1 {
		t.Errorf("Expected 1 task in manager, got %d", len(manager.tasks))
	}

	addedTask, exists := manager.tasks["task-1"]
	if !exists {
		t.Error("Task was not added to manager")
	} else if addedTask.ID != task1.ID {
		t.Errorf("Added task ID mismatch: expected %s, got %s", task1.ID, addedTask.ID)
	}

	// Test adding a nil task
	err = manager.AddTask(nil)
	if err == nil {
		t.Error("Expected error when adding nil task, got nil")
	} else if err.Error() != "task cannot be nil" {
		t.Errorf("Expected 'task cannot be nil' error, got: %v", err)
	}

	// Test adding a task with empty ID
	task2 := &types.Task{
		ID:          "",
		Name:        "Task with empty ID",
		Summary:     "Invalid task",
		Description: "Task with empty ID",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = manager.AddTask(task2)
	if err == nil {
		t.Error("Expected error when adding task with empty ID, got nil")
	} else if err.Error() != "task ID cannot be empty" {
		t.Errorf("Expected 'task ID cannot be empty' error, got: %v", err)
	}

	// Test adding a duplicate task
	task3 := &types.Task{
		ID:          "task-1",
		Name:        "Duplicate Task",
		Summary:     "This is a duplicate",
		Description: "Task with duplicate ID",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = manager.AddTask(task3)
	if err == nil {
		t.Error("Expected error when adding duplicate task, got nil")
	} else if err.Error() != "task with ID task-1 already exists" {
		t.Errorf("Expected duplicate task error, got: %v", err)
	}

	// Verify only the first task remains
	if len(manager.tasks) != 1 {
		t.Errorf("Expected 1 task in manager after duplicate attempt, got %d", len(manager.tasks))
	}

	// Test adding another valid task
	task4 := &types.Task{
		ID:          "task-4",
		Name:        "Test Task 4",
		Summary:     "Fourth test task",
		Description: "Another valid task",
		Tags:        []string{"test", "second"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = manager.AddTask(task4)
	if err != nil {
		t.Fatalf("Expected no error when adding second valid task, got: %v", err)
	}

	// Verify both tasks are present
	if len(manager.tasks) != 2 {
		t.Errorf("Expected 2 tasks in manager, got %d", len(manager.tasks))
	}
}

func TestGetTask(t *testing.T) {
	manager := NewManager()
	now := time.Now().UTC().Truncate(time.Second)

	// Add a task to the manager
	task1 := &types.Task{
		ID:          "task-1",
		Name:        "Test Task 1",
		Summary:     "First test task",
		Description: "A task to retrieve",
		Tags:        []string{"test", "retrieval"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := manager.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	// Test retrieving an existing task
	retrievedTask, err := manager.GetTask("task-1")
	if err != nil {
		t.Fatalf("Expected no error when retrieving existing task, got: %v", err)
	}
	if retrievedTask == nil {
		t.Fatal("Retrieved task is nil")
	}
	if retrievedTask.ID != task1.ID {
		t.Errorf("Retrieved task ID mismatch: expected %s, got %s", task1.ID, retrievedTask.ID)
	}
	if retrievedTask.Name != task1.Name {
		t.Errorf("Retrieved task Name mismatch: expected %s, got %s", task1.Name, retrievedTask.Name)
	}
	if retrievedTask.Summary != task1.Summary {
		t.Errorf("Retrieved task Summary mismatch: expected %s, got %s", task1.Summary, retrievedTask.Summary)
	}

	// Test retrieving a non-existent task
	_, err = manager.GetTask("non-existent")
	if err == nil {
		t.Error("Expected error when retrieving non-existent task, got nil")
	} else if err.Error() != "task with ID non-existent not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}

	// Test retrieving with empty ID
	_, err = manager.GetTask("")
	if err == nil {
		t.Error("Expected error when retrieving with empty ID, got nil")
	} else if err.Error() != "task ID cannot be empty" {
		t.Errorf("Expected 'task ID cannot be empty' error, got: %v", err)
	}

	// Add another task and verify both can be retrieved
	task2 := &types.Task{
		ID:          "task-2",
		Name:        "Test Task 2",
		Summary:     "Second test task",
		Description: "Another task to retrieve",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := manager.AddTask(task2); err != nil {
		t.Fatalf("Failed to add second task: %v", err)
	}

	// Verify we can still retrieve the first task
	retrievedTask1, err := manager.GetTask("task-1")
	if err != nil {
		t.Fatalf("Failed to retrieve task-1 after adding task-2: %v", err)
	}
	if retrievedTask1.ID != "task-1" {
		t.Errorf("Retrieved wrong task: expected task-1, got %s", retrievedTask1.ID)
	}

	// Verify we can retrieve the second task
	retrievedTask2, err := manager.GetTask("task-2")
	if err != nil {
		t.Fatalf("Failed to retrieve task-2: %v", err)
	}
	if retrievedTask2.ID != "task-2" {
		t.Errorf("Retrieved wrong task: expected task-2, got %s", retrievedTask2.ID)
	}
}

func TestGetTasks(t *testing.T) {
	manager := NewManager()
	now := time.Now().UTC().Truncate(time.Second)

	// Add multiple tasks
	task1 := &types.Task{
		ID:          "task-1",
		Name:        "Test Task 1",
		Summary:     "First test task",
		Description: "First task",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	task2 := &types.Task{
		ID:          "task-2",
		Name:        "Test Task 2",
		Summary:     "Second test task",
		Description: "Second task",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	task3 := &types.Task{
		ID:          "task-3",
		Name:        "Test Task 3",
		Summary:     "Third test task",
		Description: "Third task",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := manager.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}
	if err := manager.AddTask(task2); err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}
	if err := manager.AddTask(task3); err != nil {
		t.Fatalf("Failed to add task3: %v", err)
	}

	// Test retrieving multiple existing tasks
	tasks, err := manager.getTasks([]string{"task-1", "task-3"})
	if err != nil {
		t.Fatalf("Expected no error when retrieving existing tasks, got: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
	if tasks[0].ID != "task-1" {
		t.Errorf("Expected first task to be task-1, got %s", tasks[0].ID)
	}
	if tasks[1].ID != "task-3" {
		t.Errorf("Expected second task to be task-3, got %s", tasks[1].ID)
	}

	// Test retrieving all tasks
	allTasks, err := manager.getTasks([]string{"task-1", "task-2", "task-3"})
	if err != nil {
		t.Fatalf("Expected no error when retrieving all tasks, got: %v", err)
	}
	if len(allTasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(allTasks))
	}

	// Test retrieving with empty slice
	emptyTasks, err := manager.getTasks([]string{})
	if err != nil {
		t.Fatalf("Expected no error when retrieving with empty slice, got: %v", err)
	}
	if len(emptyTasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(emptyTasks))
	}

	// Test retrieving with some non-existent tasks
	tasks, err = manager.getTasks([]string{"task-1", "non-existent", "task-2"})
	if err == nil {
		t.Error("Expected error when some tasks don't exist, got nil")
	}
	// Should still return the found tasks
	if len(tasks) != 2 {
		t.Errorf("Expected 2 found tasks even with error, got %d", len(tasks))
	}
	if tasks[0].ID != "task-1" {
		t.Errorf("Expected first task to be task-1, got %s", tasks[0].ID)
	}
	if tasks[1].ID != "task-2" {
		t.Errorf("Expected second task to be task-2, got %s", tasks[1].ID)
	}

	// Test retrieving with all non-existent tasks
	tasks, err = manager.getTasks([]string{"non-existent-1", "non-existent-2"})
	if err == nil {
		t.Error("Expected error when all tasks don't exist, got nil")
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks when all are non-existent, got %d", len(tasks))
	}

	// Test retrieving with empty ID in slice
	_, err = manager.getTasks([]string{"task-1", "", "task-2"})
	if err == nil {
		t.Error("Expected error when ID list contains empty string, got nil")
	} else if err.Error() != "task ID cannot be empty" {
		t.Errorf("Expected 'task ID cannot be empty' error, got: %v", err)
	}
}

func TestUpdateTask(t *testing.T) {
	manager := NewManager()
	now := time.Now().UTC().Truncate(time.Second)

	// Add initial task
	task1 := &types.Task{
		ID:          "task-1",
		Name:        "Original Name",
		Summary:     "Original Summary",
		Description: "Original Description",
		Tags:        []string{"original"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := manager.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	// Test updating an existing task
	updatedTask := &types.Task{
		ID:          "task-1",
		Name:        "Updated Name",
		Summary:     "Updated Summary",
		Description: "Updated Description",
		Tags:        []string{"updated", "modified"},
		CreatedAt:   now,
		UpdatedAt:   now.Add(time.Hour),
	}

	err := manager.UpdateTask(updatedTask)
	if err != nil {
		t.Fatalf("Expected no error when updating existing task, got: %v", err)
	}

	// Verify task was updated
	retrievedTask, err := manager.GetTask("task-1")
	if err != nil {
		t.Fatalf("Failed to retrieve updated task: %v", err)
	}

	if retrievedTask.Name != "Updated Name" {
		t.Errorf("Name not updated: expected 'Updated Name', got %s", retrievedTask.Name)
	}
	if retrievedTask.Summary != "Updated Summary" {
		t.Errorf("Summary not updated: expected 'Updated Summary', got %s", retrievedTask.Summary)
	}
	if retrievedTask.Description != "Updated Description" {
		t.Errorf("Description not updated: expected 'Updated Description', got %s", retrievedTask.Description)
	}
	if len(retrievedTask.Tags) != 2 || retrievedTask.Tags[0] != "updated" {
		t.Errorf("Tags not updated: expected ['updated', 'modified'], got %v", retrievedTask.Tags)
	}

	// Test updating a non-existent task
	nonExistentTask := &types.Task{
		ID:          "non-existent",
		Name:        "Non-existent Task",
		Summary:     "This task doesn't exist",
		Description: "Should fail",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = manager.UpdateTask(nonExistentTask)
	if err == nil {
		t.Error("Expected error when updating non-existent task, got nil")
	} else if err.Error() != "task with ID non-existent not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}

	// Test updating with nil task
	err = manager.UpdateTask(nil)
	if err == nil {
		t.Error("Expected error when updating nil task, got nil")
	} else if err.Error() != "task cannot be nil" {
		t.Errorf("Expected 'task cannot be nil' error, got: %v", err)
	}

	// Test updating with empty ID
	emptyIDTask := &types.Task{
		ID:          "",
		Name:        "Task with empty ID",
		Summary:     "Invalid",
		Description: "Should fail",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = manager.UpdateTask(emptyIDTask)
	if err == nil {
		t.Error("Expected error when updating task with empty ID, got nil")
	} else if err.Error() != "task ID cannot be empty" {
		t.Errorf("Expected 'task ID cannot be empty' error, got: %v", err)
	}

	// Verify original task count unchanged after failed updates
	if len(manager.tasks) != 1 {
		t.Errorf("Expected 1 task in manager after failed updates, got %d", len(manager.tasks))
	}
}

func TestDeleteTask(t *testing.T) {
	manager := NewManager()
	now := time.Now().UTC().Truncate(time.Second)

	// Add multiple tasks
	task1 := &types.Task{
		ID:          "task-1",
		Name:        "Task 1",
		Summary:     "First task",
		Description: "To be deleted",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	task2 := &types.Task{
		ID:          "task-2",
		Name:        "Task 2",
		Summary:     "Second task",
		Description: "To be kept",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	task3 := &types.Task{
		ID:          "task-3",
		Name:        "Task 3",
		Summary:     "Third task",
		Description: "To be deleted later",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := manager.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}
	if err := manager.AddTask(task2); err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}
	if err := manager.AddTask(task3); err != nil {
		t.Fatalf("Failed to add task3: %v", err)
	}

	// Verify initial count
	if len(manager.tasks) != 3 {
		t.Errorf("Expected 3 tasks initially, got %d", len(manager.tasks))
	}

	// Test deleting an existing task
	err := manager.DeleteTask("task-1")
	if err != nil {
		t.Fatalf("Expected no error when deleting existing task, got: %v", err)
	}

	// Verify task was deleted
	if len(manager.tasks) != 2 {
		t.Errorf("Expected 2 tasks after deletion, got %d", len(manager.tasks))
	}

	_, err = manager.GetTask("task-1")
	if err == nil {
		t.Error("Expected error when retrieving deleted task, got nil")
	}

	// Verify other tasks still exist
	task2Retrieved, err := manager.GetTask("task-2")
	if err != nil {
		t.Errorf("Failed to retrieve task-2 after deleting task-1: %v", err)
	}
	if task2Retrieved.ID != "task-2" {
		t.Errorf("Retrieved wrong task: expected task-2, got %s", task2Retrieved.ID)
	}

	// Test deleting a non-existent task
	err = manager.DeleteTask("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent task, got nil")
	} else if err.Error() != "task with ID non-existent not found" {
		t.Errorf("Expected 'task not found' error, got: %v", err)
	}

	// Test deleting with empty ID
	err = manager.DeleteTask("")
	if err == nil {
		t.Error("Expected error when deleting with empty ID, got nil")
	} else if err.Error() != "task ID cannot be empty" {
		t.Errorf("Expected 'task ID cannot be empty' error, got: %v", err)
	}

	// Verify count unchanged after failed deletes
	if len(manager.tasks) != 2 {
		t.Errorf("Expected 2 tasks after failed deletes, got %d", len(manager.tasks))
	}

	// Delete another task
	err = manager.DeleteTask("task-3")
	if err != nil {
		t.Fatalf("Expected no error when deleting task-3, got: %v", err)
	}

	// Verify only task-2 remains
	if len(manager.tasks) != 1 {
		t.Errorf("Expected 1 task remaining, got %d", len(manager.tasks))
	}

	remainingTask, err := manager.GetTask("task-2")
	if err != nil {
		t.Error("Failed to retrieve remaining task-2")
	}
	if remainingTask.ID != "task-2" {
		t.Errorf("Wrong task remaining: expected task-2, got %s", remainingTask.ID)
	}
}

func TestListAllTasks(t *testing.T) {
	manager := NewManager()
	now := time.Now().UTC().Truncate(time.Second)

	// Test with empty manager
	tasks := manager.ListAllTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks in empty manager, got %d", len(tasks))
	}
	if tasks == nil {
		t.Error("Expected non-nil slice, got nil")
	}

	// Add single task
	task1 := &types.Task{
		ID:          "task-1",
		Name:        "Task 1",
		Summary:     "First task",
		Description: "Test task 1",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := manager.AddTask(task1); err != nil {
		t.Fatalf("Failed to add task1: %v", err)
	}

	tasks = manager.ListAllTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].ID != "task-1" {
		t.Errorf("Expected task-1, got %s", tasks[0].ID)
	}

	// Add more tasks
	task2 := &types.Task{
		ID:          "task-2",
		Name:        "Task 2",
		Summary:     "Second task",
		Description: "Test task 2",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	task3 := &types.Task{
		ID:          "task-3",
		Name:        "Task 3",
		Summary:     "Third task",
		Description: "Test task 3",
		Tags:        []string{"test"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := manager.AddTask(task2); err != nil {
		t.Fatalf("Failed to add task2: %v", err)
	}
	if err := manager.AddTask(task3); err != nil {
		t.Fatalf("Failed to add task3: %v", err)
	}

	tasks = manager.ListAllTasks()
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}

	// Verify all tasks are present (order doesn't matter with maps)
	taskIDs := make(map[string]bool)
	for _, task := range tasks {
		taskIDs[task.ID] = true
	}

	if !taskIDs["task-1"] {
		t.Error("task-1 not found in ListAllTasks result")
	}
	if !taskIDs["task-2"] {
		t.Error("task-2 not found in ListAllTasks result")
	}
	if !taskIDs["task-3"] {
		t.Error("task-3 not found in ListAllTasks result")
	}

	// Delete a task and verify list updates
	if err := manager.DeleteTask("task-2"); err != nil {
		t.Fatalf("Failed to delete task-2: %v", err)
	}

	tasks = manager.ListAllTasks()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks after deletion, got %d", len(tasks))
	}

	// Verify task-2 is not in the list
	for _, task := range tasks {
		if task.ID == "task-2" {
			t.Error("Deleted task-2 still appears in ListAllTasks result")
		}
	}
}
