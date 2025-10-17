package task_manager

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"common-tasks-mcp/pkg/task_manager/types"
)

func TestPersistAndLoad(t *testing.T) {
	// Create a temporary directory for test files
	testDir := filepath.Join(t.TempDir(), "tasks")

	// Create a manager with 3 tasks
	manager1 := NewManager()

	now := time.Now().UTC().Truncate(time.Second)

	// Task A - no dependencies
	taskA := &types.Task{
		ID:                    "task-a",
		Name:                  "Task A",
		Summary:               "First task",
		Description:           "This is the first task with no dependencies",
		Tags:                  []string{"api", "backend"},
		PrerequisiteIDs:       []string{},
		DownstreamRequiredIDs: []string{"task-b"},
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// Task B - depends on A, has C as dependent
	taskB := &types.Task{
		ID:                    "task-b",
		Name:                  "Task B",
		Summary:               "Second task",
		Description:           "This task depends on A and has C as dependent",
		Tags:                  []string{"frontend", "api"},
		PrerequisiteIDs:       []string{"task-a"},
		DownstreamRequiredIDs: []string{"task-c"},
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// Task C - depends on B
	taskC := &types.Task{
		ID:                    "task-c",
		Name:                  "Task C",
		Summary:               "Third task",
		Description:           "This task depends on B",
		Tags:                  []string{"testing"},
		PrerequisiteIDs:       []string{"task-b"},
		DownstreamRequiredIDs: []string{},
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	manager1.tasks["task-a"] = taskA
	manager1.tasks["task-b"] = taskB
	manager1.tasks["task-c"] = taskC

	// PersistToDir tasks to directory
	if err := manager1.PersistToDir(testDir); err != nil {
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
	if err := manager2.LoadFromDir(testDir); err != nil {
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

		if !originalTask.Equals(loadedTask) {
			t.Errorf("Task %s does not match after persist/load cycle", id)
		}
	}
}

func TestAddTask(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name           string
		setupTasks     []*types.Task // Tasks to add before the test task
		taskToAdd      *types.Task
		wantError      bool
		expectedError  string
		expectedCount  int // Expected number of tasks after operation
		validateResult func(t *testing.T, manager *Manager)
	}{
		{
			name:       "add valid task to empty manager",
			setupTasks: []*types.Task{},
			taskToAdd: &types.Task{
				ID:          "task-1",
				Name:        "Test Task 1",
				Summary:     "First test task",
				Description: "A valid task to add",
				Tags:        []string{"test"},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantError:     false,
			expectedCount: 1,
			validateResult: func(t *testing.T, manager *Manager) {
				addedTask, exists := manager.tasks["task-1"]
				if !exists {
					t.Error("Task was not added to manager")
				} else if addedTask.ID != "task-1" {
					t.Errorf("Added task ID mismatch: expected task-1, got %s", addedTask.ID)
				}
			},
		},
		{
			name:          "add nil task",
			setupTasks:    []*types.Task{},
			taskToAdd:     nil,
			wantError:     true,
			expectedError: "task cannot be nil",
			expectedCount: 0,
		},
		{
			name:       "add task with empty ID",
			setupTasks: []*types.Task{},
			taskToAdd: &types.Task{
				ID:          "",
				Name:        "Task with empty ID",
				Summary:     "Invalid task",
				Description: "Task with empty ID",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantError:     true,
			expectedError: "task ID cannot be empty",
			expectedCount: 0,
		},
		{
			name: "add duplicate task",
			setupTasks: []*types.Task{
				{
					ID:          "task-1",
					Name:        "Test Task 1",
					Summary:     "First test task",
					Description: "A valid task to add",
					Tags:        []string{"test"},
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			taskToAdd: &types.Task{
				ID:          "task-1",
				Name:        "Duplicate Task",
				Summary:     "This is a duplicate",
				Description: "Task with duplicate ID",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantError:     true,
			expectedError: "task with ID task-1 already exists",
			expectedCount: 1,
		},
		{
			name: "add second valid task",
			setupTasks: []*types.Task{
				{
					ID:          "task-1",
					Name:        "Test Task 1",
					Summary:     "First test task",
					Description: "A valid task to add",
					Tags:        []string{"test"},
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			taskToAdd: &types.Task{
				ID:          "task-4",
				Name:        "Test Task 4",
				Summary:     "Fourth test task",
				Description: "Another valid task",
				Tags:        []string{"test", "second"},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantError:     false,
			expectedCount: 2,
			validateResult: func(t *testing.T, manager *Manager) {
				task1, exists1 := manager.tasks["task-1"]
				if !exists1 {
					t.Error("task-1 should still exist")
				} else if task1.Name != "Test Task 1" {
					t.Errorf("task-1 name mismatch: expected 'Test Task 1', got %q", task1.Name)
				}

				task4, exists4 := manager.tasks["task-4"]
				if !exists4 {
					t.Error("task-4 was not added")
				} else if task4.Name != "Test Task 4" {
					t.Errorf("task-4 name mismatch: expected 'Test Task 4', got %q", task4.Name)
				}
			},
		},
		{
			name: "add task that creates cycle - simple two-task cycle",
			setupTasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					Summary:               "First task",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-b"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			taskToAdd: &types.Task{
				ID:                    "task-b",
				Name:                  "Task B",
				Summary:               "Second task that creates a cycle",
				PrerequisiteIDs:       []string{"task-a"},
				DownstreamRequiredIDs: []string{"task-a"}, // This creates a cycle: A -> B -> A
				CreatedAt:             now,
				UpdatedAt:             now,
			},
			wantError:     true,
			expectedError: "cycle",
			expectedCount: 1, // task-b should not be added
		},
		{
			name: "add task that creates cycle - three-task cycle",
			setupTasks: []*types.Task{
				{
					ID:                    "task-x",
					Name:                  "Task X",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-y"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-y",
					Name:                  "Task Y",
					PrerequisiteIDs:       []string{"task-x"},
					DownstreamRequiredIDs: []string{"task-z"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			taskToAdd: &types.Task{
				ID:                    "task-z",
				Name:                  "Task Z",
				Summary:               "Third task that completes a cycle",
				PrerequisiteIDs:       []string{"task-y"},
				DownstreamRequiredIDs: []string{"task-x"}, // This creates a cycle: X -> Y -> Z -> X
				CreatedAt:             now,
				UpdatedAt:             now,
			},
			wantError:     true,
			expectedError: "cycle",
			expectedCount: 2, // task-z should not be added
		},
		{
			name: "add task with self-cycle",
			setupTasks: []*types.Task{
				{
					ID:        "task-1",
					Name:      "Task 1",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			taskToAdd: &types.Task{
				ID:                    "task-self",
				Name:                  "Self-referencing Task",
				PrerequisiteIDs:       []string{"task-self"}, // Self-cycle
				DownstreamRequiredIDs: []string{},
				CreatedAt:             now,
				UpdatedAt:             now,
			},
			wantError:     true,
			expectedError: "cycle",
			expectedCount: 1, // task-self should not be added
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()

			// Add setup tasks
			for _, task := range tt.setupTasks {
				if err := manager.AddTask(task); err != nil {
					t.Fatalf("Failed to add setup task %s: %v", task.ID, err)
				}
			}

			// Perform the test operation
			err := manager.AddTask(tt.taskToAdd)

			// Check error expectations
			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.expectedError)
				} else if tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing %q, got %q", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
			}

			// Verify task count
			if len(manager.tasks) != tt.expectedCount {
				t.Errorf("Expected %d tasks in manager, got %d", tt.expectedCount, len(manager.tasks))
			}

			// Run custom validation if provided and no error expected
			if tt.validateResult != nil && !tt.wantError {
				tt.validateResult(t, manager)
			}
		})
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
	if !retrievedTask.Equals(task1) {
		t.Error("Retrieved task does not match original task")
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
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name            string
		setupTasks      []*types.Task
		updateTask      *types.Task
		wantError       bool
		expectedError   string
		validatePointer func(t *testing.T, manager *Manager)
	}{
		{
			name: "update existing task",
			setupTasks: []*types.Task{
				{
					ID:          "task-1",
					Name:        "Original Name",
					Summary:     "Original Summary",
					Description: "Original Description",
					Tags:        []string{"original"},
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			updateTask: &types.Task{
				ID:          "task-1",
				Name:        "Updated Name",
				Summary:     "Updated Summary",
				Description: "Updated Description",
				Tags:        []string{"updated", "modified"},
				CreatedAt:   now,
				UpdatedAt:   now.Add(time.Hour),
			},
			wantError: false,
			validatePointer: func(t *testing.T, manager *Manager) {
				retrievedTask, err := manager.GetTask("task-1")
				if err != nil {
					t.Fatalf("Failed to retrieve updated task: %v", err)
				}
				if retrievedTask.Name != "Updated Name" {
					t.Errorf("Expected name 'Updated Name', got %q", retrievedTask.Name)
				}
				if len(retrievedTask.Tags) != 2 || retrievedTask.Tags[0] != "updated" {
					t.Errorf("Expected tags [updated, modified], got %v", retrievedTask.Tags)
				}
			},
		},
		{
			name: "update non-existent task",
			setupTasks: []*types.Task{
				{
					ID:        "task-1",
					Name:      "Existing Task",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			updateTask: &types.Task{
				ID:          "non-existent",
				Name:        "Non-existent Task",
				Summary:     "This task doesn't exist",
				Description: "Should fail",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantError:     true,
			expectedError: "task with ID non-existent not found",
		},
		{
			name: "update with nil task",
			setupTasks: []*types.Task{
				{
					ID:        "task-1",
					Name:      "Existing Task",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			updateTask:    nil,
			wantError:     true,
			expectedError: "task cannot be nil",
		},
		{
			name: "update with empty ID",
			setupTasks: []*types.Task{
				{
					ID:        "task-1",
					Name:      "Existing Task",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			updateTask: &types.Task{
				ID:          "",
				Name:        "Task with empty ID",
				Summary:     "Invalid",
				Description: "Should fail",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantError:     true,
			expectedError: "task ID cannot be empty",
		},
		{
			name: "stale pointer: prerequisite pointer becomes stale",
			setupTasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A - Version 1",
					Summary:               "First version",
					Description:           "Original task A",
					Tags:                  []string{"v1"},
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-b"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-b",
					Name:                  "Task B",
					Summary:               "Middle task",
					Description:           "Depends on A",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			updateTask: &types.Task{
				ID:                    "task-a",
				Name:                  "Task A - Version 2",
				Summary:               "Second version",
				Description:           "Updated task A",
				Tags:                  []string{"v2", "updated"},
				PrerequisiteIDs:       []string{},
				DownstreamRequiredIDs: []string{"task-b"},
				CreatedAt:             now,
				UpdatedAt:             now.Add(time.Hour),
			},
			wantError: false,
			validatePointer: func(t *testing.T, manager *Manager) {
				retrievedB, err := manager.GetTask("task-b")
				if err != nil {
					t.Fatalf("Failed to retrieve task-b: %v", err)
				}
				if len(retrievedB.Prerequisites) != 1 {
					t.Fatalf("Expected task-b to have 1 prerequisite, got %d", len(retrievedB.Prerequisites))
				}
				prereqA := retrievedB.Prerequisites[0]
				if prereqA.Name == "Task A - Version 1" {
					t.Errorf("BUG: task-b.Prerequisites[0] points to OLD task-a (Version 1)")
					t.Errorf("  Expected Name='Task A - Version 2', got %q", prereqA.Name)
				}
			},
		},
		{
			name: "stale pointer: downstream required pointer becomes stale",
			setupTasks: []*types.Task{
				{
					ID:                    "task-x",
					Name:                  "Task X",
					Summary:               "First task",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-y"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-y",
					Name:                  "Task Y - Version 1",
					Summary:               "Original version",
					Tags:                  []string{"old"},
					PrerequisiteIDs:       []string{"task-x"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			updateTask: &types.Task{
				ID:                    "task-y",
				Name:                  "Task Y - Version 2",
				Summary:               "Updated version",
				Tags:                  []string{"new"},
				PrerequisiteIDs:       []string{"task-x"},
				DownstreamRequiredIDs: []string{},
				CreatedAt:             now,
				UpdatedAt:             now.Add(time.Hour),
			},
			wantError: false,
			validatePointer: func(t *testing.T, manager *Manager) {
				retrievedX, err := manager.GetTask("task-x")
				if err != nil {
					t.Fatalf("Failed to retrieve task-x: %v", err)
				}
				if len(retrievedX.DownstreamRequired) != 1 {
					t.Fatalf("Expected task-x to have 1 downstream required, got %d", len(retrievedX.DownstreamRequired))
				}
				downstreamY := retrievedX.DownstreamRequired[0]
				if downstreamY.Name == "Task Y - Version 1" {
					t.Errorf("BUG: task-x.DownstreamRequired[0] points to OLD task-y (Version 1)")
					t.Errorf("  Expected Name='Task Y - Version 2', got %q", downstreamY.Name)
				}
			},
		},
		{
			name: "stale pointer: downstream suggested pointer becomes stale",
			setupTasks: []*types.Task{
				{
					ID:                     "task-p",
					Name:                   "Task P",
					Summary:                "First task",
					PrerequisiteIDs:        []string{},
					DownstreamRequiredIDs:  []string{},
					DownstreamSuggestedIDs: []string{"task-q"},
					CreatedAt:              now,
					UpdatedAt:              now,
				},
				{
					ID:                    "task-q",
					Name:                  "Task Q - Version 1",
					Summary:               "Original version",
					Tags:                  []string{"old"},
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			updateTask: &types.Task{
				ID:                    "task-q",
				Name:                  "Task Q - Version 2",
				Summary:               "Updated version",
				Tags:                  []string{"new"},
				PrerequisiteIDs:       []string{},
				DownstreamRequiredIDs: []string{},
				CreatedAt:             now,
				UpdatedAt:             now.Add(time.Hour),
			},
			wantError: false,
			validatePointer: func(t *testing.T, manager *Manager) {
				retrievedP, err := manager.GetTask("task-p")
				if err != nil {
					t.Fatalf("Failed to retrieve task-p: %v", err)
				}
				if len(retrievedP.DownstreamSuggested) != 1 {
					t.Fatalf("Expected task-p to have 1 downstream suggested, got %d", len(retrievedP.DownstreamSuggested))
				}
				suggestedQ := retrievedP.DownstreamSuggested[0]
				if suggestedQ.Name == "Task Q - Version 1" {
					t.Errorf("BUG: task-p.DownstreamSuggested[0] points to OLD task-q (Version 1)")
					t.Errorf("  Expected Name='Task Q - Version 2', got %q", suggestedQ.Name)
				}
			},
		},
		{
			name: "stale pointer: nested pointer becomes stale (multi-hop)",
			setupTasks: []*types.Task{
				{
					ID:                    "task-1",
					Name:                  "Task 1 - Version 1",
					Summary:               "Original",
					Tags:                  []string{"v1"},
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-2"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-2",
					Name:                  "Task 2",
					Summary:               "Middle task",
					PrerequisiteIDs:       []string{"task-1"},
					DownstreamRequiredIDs: []string{"task-3"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-3",
					Name:                  "Task 3",
					Summary:               "End task",
					PrerequisiteIDs:       []string{"task-2"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			updateTask: &types.Task{
				ID:                    "task-1",
				Name:                  "Task 1 - Version 2",
				Summary:               "Updated",
				Tags:                  []string{"v2"},
				PrerequisiteIDs:       []string{},
				DownstreamRequiredIDs: []string{"task-2"},
				CreatedAt:             now,
				UpdatedAt:             now.Add(time.Hour),
			},
			wantError: false,
			validatePointer: func(t *testing.T, manager *Manager) {
				// Check direct pointer: task-2 -> task-1
				retrieved2, err := manager.GetTask("task-2")
				if err != nil {
					t.Fatalf("Failed to retrieve task-2: %v", err)
				}
				if len(retrieved2.Prerequisites) == 1 {
					prereq1 := retrieved2.Prerequisites[0]
					if prereq1.Name == "Task 1 - Version 1" {
						t.Errorf("BUG: task-2.Prerequisites[0] points to OLD task-1 (Version 1)")
					}
				}

				// Check nested pointer: task-3 -> task-2 -> task-1
				retrieved3, err := manager.GetTask("task-3")
				if err != nil {
					t.Fatalf("Failed to retrieve task-3: %v", err)
				}
				if len(retrieved3.Prerequisites) == 1 {
					task2FromTask3 := retrieved3.Prerequisites[0]
					if len(task2FromTask3.Prerequisites) == 1 {
						task1FromTask2 := task2FromTask3.Prerequisites[0]
						if task1FromTask2.Name == "Task 1 - Version 1" {
							t.Errorf("BUG: Nested pointer task-3 -> task-2 -> task-1 is stale")
							t.Errorf("  Expected Name='Task 1 - Version 2', got %q", task1FromTask2.Name)
						}
					}
				}
			},
		},
		{
			name: "stale pointer: task referenced in multiple relationships",
			setupTasks: []*types.Task{
				{
					ID:                    "task-hub",
					Name:                  "Hub Task - Version 1",
					Summary:               "Central task",
					Tags:                  []string{"v1"},
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-dep1"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:              "task-dep1",
					Name:            "Dependent 1",
					PrerequisiteIDs: []string{"task-hub"},
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:              "task-dep2",
					Name:            "Dependent 2",
					PrerequisiteIDs: []string{"task-hub"},
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:                     "task-dep3",
					Name:                   "Dependent 3",
					DownstreamSuggestedIDs: []string{"task-hub"},
					CreatedAt:              now,
					UpdatedAt:              now,
				},
			},
			updateTask: &types.Task{
				ID:                    "task-hub",
				Name:                  "Hub Task - Version 2",
				Summary:               "Updated central task",
				Tags:                  []string{"v2"},
				PrerequisiteIDs:       []string{},
				DownstreamRequiredIDs: []string{"task-dep1"},
				CreatedAt:             now,
				UpdatedAt:             now.Add(time.Hour),
			},
			wantError: false,
			validatePointer: func(t *testing.T, manager *Manager) {
				// Check task-dep1's prerequisite pointer
				dep1, _ := manager.GetTask("task-dep1")
				if len(dep1.Prerequisites) == 1 && dep1.Prerequisites[0].Name == "Hub Task - Version 1" {
					t.Errorf("BUG: task-dep1 has stale pointer to Hub Task (Version 1)")
				}

				// Check task-dep2's prerequisite pointer
				dep2, _ := manager.GetTask("task-dep2")
				if len(dep2.Prerequisites) == 1 && dep2.Prerequisites[0].Name == "Hub Task - Version 1" {
					t.Errorf("BUG: task-dep2 has stale pointer to Hub Task (Version 1)")
				}

				// Check task-dep3's suggested downstream pointer
				dep3, _ := manager.GetTask("task-dep3")
				if len(dep3.DownstreamSuggested) == 1 && dep3.DownstreamSuggested[0].Name == "Hub Task - Version 1" {
					t.Errorf("BUG: task-dep3 has stale pointer to Hub Task (Version 1)")
				}
			},
		},
		{
			name: "update that would introduce cycle is rejected - prerequisites",
			setupTasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-b"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-b",
					Name:                  "Task B",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			updateTask: &types.Task{
				ID:                    "task-a",
				Name:                  "Task A - Updated",
				PrerequisiteIDs:       []string{"task-b"}, // This would create a cycle: A -> B -> A
				DownstreamRequiredIDs: []string{"task-b"},
				CreatedAt:             now,
				UpdatedAt:             now.Add(time.Hour),
			},
			wantError:     true,
			expectedError: "update would introduce cycle",
		},
		{
			name: "update that would introduce cycle is rejected - downstream required",
			setupTasks: []*types.Task{
				{
					ID:                    "task-x",
					Name:                  "Task X",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-y"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-y",
					Name:                  "Task Y",
					PrerequisiteIDs:       []string{"task-x"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			updateTask: &types.Task{
				ID:                    "task-y",
				Name:                  "Task Y - Updated",
				PrerequisiteIDs:       []string{"task-x"},
				DownstreamRequiredIDs: []string{"task-x"}, // This would create a cycle: X -> Y -> X
				CreatedAt:             now,
				UpdatedAt:             now.Add(time.Hour),
			},
			wantError:     true,
			expectedError: "update would introduce cycle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()

			// Add all setup tasks
			for _, task := range tt.setupTasks {
				if err := manager.AddTask(task); err != nil {
					t.Fatalf("Failed to add task %s: %v", task.ID, err)
				}
			}

			// For stale pointer tests, resolve pointers before update
			if tt.validatePointer != nil && !tt.wantError {
				if err := manager.ResolveTaskPointers(); err != nil {
					t.Fatalf("Failed to resolve task pointers: %v", err)
				}
			}

			// Perform the update
			err := manager.UpdateTask(tt.updateTask)

			// Check error expectations
			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error containing %q, got nil", tt.expectedError)
				} else if tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing %q, got %q", tt.expectedError, err.Error())
				}
				return // Don't run validations if we expected an error
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Verify the manager has the updated task
			if tt.updateTask != nil {
				retrievedTask, err := manager.GetTask(tt.updateTask.ID)
				if err != nil {
					t.Fatalf("Failed to retrieve updated task %s: %v", tt.updateTask.ID, err)
				}
				if !retrievedTask.Equals(tt.updateTask) {
					t.Errorf("Retrieved task does not match updated task")
				}
			}

			// Run custom validation if provided
			if tt.validatePointer != nil {
				tt.validatePointer(t, manager)
			}
		})
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

func TestDeleteTaskWithReferences(t *testing.T) {
	manager := NewManager()
	now := time.Now().UTC().Truncate(time.Second)

	// Create a task graph: task-a -> task-b -> task-c
	// task-b also has task-d as a suggested downstream task
	taskA := &types.Task{
		ID:                    "task-a",
		Name:                  "Task A",
		Summary:               "First task",
		Description:           "Root task",
		PrerequisiteIDs:       []string{},
		DownstreamRequiredIDs: []string{"task-b"},
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	taskB := &types.Task{
		ID:                     "task-b",
		Name:                   "Task B",
		Summary:                "Second task",
		Description:            "Middle task",
		PrerequisiteIDs:        []string{"task-a"},
		DownstreamRequiredIDs:  []string{"task-c"},
		DownstreamSuggestedIDs: []string{"task-d"},
		CreatedAt:              now,
		UpdatedAt:              now,
	}

	taskC := &types.Task{
		ID:                    "task-c",
		Name:                  "Task C",
		Summary:               "Third task",
		Description:           "End task",
		PrerequisiteIDs:       []string{"task-b"},
		DownstreamRequiredIDs: []string{},
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	taskD := &types.Task{
		ID:                    "task-d",
		Name:                  "Task D",
		Summary:               "Fourth task",
		Description:           "Suggested task",
		PrerequisiteIDs:       []string{},
		DownstreamRequiredIDs: []string{},
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	// Add all tasks
	for _, task := range []*types.Task{taskA, taskB, taskC, taskD} {
		if err := manager.AddTask(task); err != nil {
			t.Fatalf("Failed to add task %s: %v", task.ID, err)
		}
	}

	// Resolve pointers so we can test pointer cleanup too
	if err := manager.ResolveTaskPointers(); err != nil {
		t.Fatalf("Failed to resolve task pointers: %v", err)
	}

	// Verify initial state
	if len(manager.tasks) != 4 {
		t.Fatalf("Expected 4 tasks initially, got %d", len(manager.tasks))
	}

	// Delete task-b (which is referenced by task-a and task-c)
	err := manager.DeleteTask("task-b")
	if err != nil {
		t.Fatalf("Failed to delete task-b: %v", err)
	}

	// Verify task-b is deleted
	if len(manager.tasks) != 3 {
		t.Errorf("Expected 3 tasks after deletion, got %d", len(manager.tasks))
	}

	_, err = manager.GetTask("task-b")
	if err == nil {
		t.Error("task-b should have been deleted")
	}

	// Verify task-a's DownstreamRequiredIDs no longer contains task-b
	retrievedA, err := manager.GetTask("task-a")
	if err != nil {
		t.Fatalf("Failed to retrieve task-a: %v", err)
	}

	if len(retrievedA.DownstreamRequiredIDs) != 0 {
		t.Errorf("task-a.DownstreamRequiredIDs should be empty, got %v", retrievedA.DownstreamRequiredIDs)
	}
	for _, id := range retrievedA.DownstreamRequiredIDs {
		if id == "task-b" {
			t.Error("task-a.DownstreamRequiredIDs should not contain task-b")
		}
	}

	// Verify task-a's DownstreamRequired pointer slice is also cleaned
	if len(retrievedA.DownstreamRequired) != 0 {
		t.Errorf("task-a.DownstreamRequired should be empty, got %v", retrievedA.DownstreamRequired)
	}
	for _, task := range retrievedA.DownstreamRequired {
		if task != nil && task.ID == "task-b" {
			t.Error("task-a.DownstreamRequired should not contain task-b")
		}
	}

	// Verify task-c's PrerequisiteIDs no longer contains task-b
	retrievedC, err := manager.GetTask("task-c")
	if err != nil {
		t.Fatalf("Failed to retrieve task-c: %v", err)
	}

	if len(retrievedC.PrerequisiteIDs) != 0 {
		t.Errorf("task-c.PrerequisiteIDs should be empty, got %v", retrievedC.PrerequisiteIDs)
	}
	for _, id := range retrievedC.PrerequisiteIDs {
		if id == "task-b" {
			t.Error("task-c.PrerequisiteIDs should not contain task-b")
		}
	}

	// Verify task-c's Prerequisites pointer slice is also cleaned
	if len(retrievedC.Prerequisites) != 0 {
		t.Errorf("task-c.Prerequisites should be empty, got %v", retrievedC.Prerequisites)
	}
	for _, task := range retrievedC.Prerequisites {
		if task != nil && task.ID == "task-b" {
			t.Error("task-c.Prerequisites should not contain task-b")
		}
	}

	// Verify task-d is unaffected (it was only in task-b's suggested list)
	retrievedD, err := manager.GetTask("task-d")
	if err != nil {
		t.Fatalf("Failed to retrieve task-d: %v", err)
	}
	if retrievedD.ID != "task-d" {
		t.Error("task-d should still exist and be unchanged")
	}

	// Now delete task-d and verify task-b's references would have been cleaned
	// (but task-b is already deleted, so this just verifies the operation doesn't fail)
	err = manager.DeleteTask("task-d")
	if err != nil {
		t.Fatalf("Failed to delete task-d: %v", err)
	}

	if len(manager.tasks) != 2 {
		t.Errorf("Expected 2 tasks after deleting task-d, got %d", len(manager.tasks))
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

func TestDetectCycles(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	tests := []struct {
		name      string
		tasks     []*types.Task
		wantError bool
		errorMsg  string
	}{
		{
			name: "no cycles - valid DAG",
			tasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-b"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-b",
					Name:                  "Task B",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{"task-c"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-c",
					Name:                  "Task C",
					PrerequisiteIDs:       []string{"task-b"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantError: false,
		},
		{
			name: "self-cycle in prerequisites",
			tasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantError: true,
			errorMsg:  "Prerequisites DAG: task-a -> task-a",
		},
		{
			name: "self-cycle in downstream required",
			tasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-a"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantError: true,
			errorMsg:  "Downstream Required DAG: task-a -> task-a",
		},
		{
			name: "self-cycle in downstream suggested",
			tasks: []*types.Task{
				{
					ID:                     "task-a",
					Name:                   "Task A",
					PrerequisiteIDs:        []string{},
					DownstreamRequiredIDs:  []string{},
					DownstreamSuggestedIDs: []string{"task-a"},
					CreatedAt:              now,
					UpdatedAt:              now,
				},
			},
			wantError: true,
			errorMsg:  "Downstream Suggested DAG: task-a -> task-a",
		},
		{
			name: "two-task cycle in prerequisites",
			tasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					PrerequisiteIDs:       []string{"task-b"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-b",
					Name:                  "Task B",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantError: true,
			errorMsg:  "Prerequisites DAG:",
		},
		{
			name: "two-task cycle in downstream required",
			tasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-b"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-b",
					Name:                  "Task B",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-a"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantError: true,
			errorMsg:  "Downstream Required DAG:",
		},
		{
			name: "three-task cycle in prerequisites",
			tasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					PrerequisiteIDs:       []string{"task-c"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-b",
					Name:                  "Task B",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-c",
					Name:                  "Task C",
					PrerequisiteIDs:       []string{"task-b"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantError: true,
			errorMsg:  "Prerequisites DAG:",
		},
		{
			name: "longer cycle - 4 tasks in prerequisites",
			tasks: []*types.Task{
				{
					ID:              "task-a",
					Name:            "Task A",
					PrerequisiteIDs: []string{"task-d"},
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:              "task-b",
					Name:            "Task B",
					PrerequisiteIDs: []string{"task-a"},
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:              "task-c",
					Name:            "Task C",
					PrerequisiteIDs: []string{"task-b"},
					CreatedAt:       now,
					UpdatedAt:       now,
				},
				{
					ID:              "task-d",
					Name:            "Task D",
					PrerequisiteIDs: []string{"task-c"},
					CreatedAt:       now,
					UpdatedAt:       now,
				},
			},
			wantError: true,
			errorMsg:  "Prerequisites DAG:",
		},
		{
			name: "diamond pattern - not a cycle, valid",
			tasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-b", "task-c"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-b",
					Name:                  "Task B",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{"task-d"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-c",
					Name:                  "Task C",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{"task-d"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-d",
					Name:                  "Task D",
					PrerequisiteIDs:       []string{"task-b", "task-c"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantError: false,
		},
		{
			name: "cycles in multiple DAGs",
			tasks: []*types.Task{
				{
					ID:                     "task-a",
					Name:                   "Task A",
					PrerequisiteIDs:        []string{"task-b"},
					DownstreamRequiredIDs:  []string{},
					DownstreamSuggestedIDs: []string{"task-c"},
					CreatedAt:              now,
					UpdatedAt:              now,
				},
				{
					ID:                     "task-b",
					Name:                   "Task B",
					PrerequisiteIDs:        []string{"task-a"},
					DownstreamRequiredIDs:  []string{},
					DownstreamSuggestedIDs: []string{},
					CreatedAt:              now,
					UpdatedAt:              now,
				},
				{
					ID:                     "task-c",
					Name:                   "Task C",
					PrerequisiteIDs:        []string{},
					DownstreamRequiredIDs:  []string{},
					DownstreamSuggestedIDs: []string{"task-a"},
					CreatedAt:              now,
					UpdatedAt:              now,
				},
			},
			wantError: true,
			errorMsg:  "detected 2 cycle(s)",
		},
		{
			name: "complex valid graph with multiple paths",
			tasks: []*types.Task{
				{
					ID:                    "task-a",
					Name:                  "Task A",
					PrerequisiteIDs:       []string{},
					DownstreamRequiredIDs: []string{"task-b", "task-c"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-b",
					Name:                  "Task B",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{"task-d"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-c",
					Name:                  "Task C",
					PrerequisiteIDs:       []string{"task-a"},
					DownstreamRequiredIDs: []string{"task-e"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-d",
					Name:                  "Task D",
					PrerequisiteIDs:       []string{"task-b"},
					DownstreamRequiredIDs: []string{"task-f"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-e",
					Name:                  "Task E",
					PrerequisiteIDs:       []string{"task-c"},
					DownstreamRequiredIDs: []string{"task-f"},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
				{
					ID:                    "task-f",
					Name:                  "Task F",
					PrerequisiteIDs:       []string{"task-d", "task-e"},
					DownstreamRequiredIDs: []string{},
					CreatedAt:             now,
					UpdatedAt:             now,
				},
			},
			wantError: false,
		},
		{
			name: "suggested downstream cycle only",
			tasks: []*types.Task{
				{
					ID:                     "task-a",
					Name:                   "Task A",
					PrerequisiteIDs:        []string{},
					DownstreamRequiredIDs:  []string{},
					DownstreamSuggestedIDs: []string{"task-b"},
					CreatedAt:              now,
					UpdatedAt:              now,
				},
				{
					ID:                     "task-b",
					Name:                   "Task B",
					PrerequisiteIDs:        []string{},
					DownstreamRequiredIDs:  []string{},
					DownstreamSuggestedIDs: []string{"task-c"},
					CreatedAt:              now,
					UpdatedAt:              now,
				},
				{
					ID:                     "task-c",
					Name:                   "Task C",
					PrerequisiteIDs:        []string{},
					DownstreamRequiredIDs:  []string{},
					DownstreamSuggestedIDs: []string{"task-a"},
					CreatedAt:              now,
					UpdatedAt:              now,
				},
			},
			wantError: true,
			errorMsg:  "Downstream Suggested DAG:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()

			// Add all tasks to the manager directly (bypassing AddTask validation
			// since we're testing DetectCycles itself and need to create invalid graphs)
			for _, task := range tt.tasks {
				manager.tasks[task.ID] = task
			}

			// Check for cycles
			err := manager.DetectCycles()

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error containing %q, but got no error", tt.errorMsg)
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
