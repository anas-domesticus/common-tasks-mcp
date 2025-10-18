package graph_manager

import (
	"fmt"
	"os"
	"path/filepath"

	"common-tasks-mcp/pkg/graph_manager/types"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// Manager handles task graph operations
type Manager struct {
	nodes    map[string]*types.Node
	tagCache map[string][]*types.Node
	logger   *zap.Logger
}

// NewManager creates a new task manager instance
func NewManager(logger *zap.Logger) *Manager {
	logger.Debug("Creating new task manager")
	return &Manager{
		nodes:    make(map[string]*types.Node),
		tagCache: make(map[string][]*types.Node),
		logger:   logger,
	}
}

// AddNode adds a task to the manager.
// It uses a clone-validate-commit pattern to ensure the addition doesn't introduce cycles.
func (m *Manager) AddNode(task *types.Node) error {
	m.logger.Debug("Adding task")

	if task == nil {
		m.logger.Error("Attempted to add nil task")
		return fmt.Errorf("task cannot be nil")
	}
	if task.ID == "" {
		m.logger.Error("Attempted to add task with empty ID")
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.nodes[task.ID]; exists {
		m.logger.Warn("Node already exists", zap.String("task_id", task.ID))
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	m.logger.Debug("Validating task addition for cycles", zap.String("task_id", task.ID))

	// Clone the manager to test the addition
	testManager := m.Clone()

	// Perform the addition in the test manager
	testManager.nodes[task.ID] = task

	// Check for cycles in the test manager
	if err := testManager.DetectCycles(); err != nil {
		m.logger.Error("Node addition would introduce cycle",
			zap.String("task_id", task.ID),
			zap.Error(err),
		)
		return fmt.Errorf("addition would introduce cycle: %w", err)
	}

	// If no cycles detected, commit the addition to the original manager
	m.nodes[task.ID] = task
	m.logger.Debug("Node added to internal storage", zap.String("task_id", task.ID))

	// Update tag cache with the new task
	m.PopulateTagCache()
	m.logger.Info("Node added successfully",
		zap.String("task_id", task.ID),
		zap.String("task_name", task.Name),
		zap.Int("total_nodes", len(m.nodes)),
	)

	return nil
}

// UpdateNode updates an existing task in the manager.
// It uses a clone-validate-commit pattern to ensure the update doesn't introduce cycles,
// and automatically refreshes all task pointers to prevent stale references.
func (m *Manager) UpdateNode(task *types.Node) error {
	m.logger.Debug("Updating task")

	if task == nil {
		m.logger.Error("Attempted to update with nil task")
		return fmt.Errorf("task cannot be nil")
	}
	if task.ID == "" {
		m.logger.Error("Attempted to update task with empty ID")
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.nodes[task.ID]; !exists {
		m.logger.Warn("Node not found for update", zap.String("task_id", task.ID))
		return fmt.Errorf("task with ID %s not found", task.ID)
	}

	m.logger.Debug("Validating task update for cycles", zap.String("task_id", task.ID))

	// Clone the manager to test the update
	testManager := m.Clone()

	// Perform the update in the test manager
	testManager.nodes[task.ID] = task

	// Check for cycles in the test manager
	if err := testManager.DetectCycles(); err != nil {
		m.logger.Error("Node update would introduce cycle",
			zap.String("task_id", task.ID),
			zap.Error(err),
		)
		return fmt.Errorf("update would introduce cycle: %w", err)
	}

	// If no cycles detected, commit the update to the original manager
	m.nodes[task.ID] = task
	m.logger.Debug("Node updated in internal storage", zap.String("task_id", task.ID))

	// Resolve all task pointers to fix stale references
	// This ensures that any nodes pointing to the updated task get fresh pointers
	m.logger.Debug("Resolving task pointers after update")
	if err := m.ResolveNodePointers(); err != nil {
		m.logger.Error("Failed to resolve task pointers", zap.Error(err))
		return err
	}

	// Update tag cache since tags may have changed
	m.PopulateTagCache()
	m.logger.Info("Node updated successfully",
		zap.String("task_id", task.ID),
		zap.String("task_name", task.Name),
	)

	return nil
}

// DeleteNode removes a task from the manager and cleans up all references to it
// from other nodes' edge lists.
func (m *Manager) DeleteNode(id string) error {
	m.logger.Debug("Deleting task", zap.String("task_id", id))

	if id == "" {
		m.logger.Error("Attempted to delete task with empty ID")
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := m.nodes[id]; !exists {
		m.logger.Warn("Node not found for deletion", zap.String("task_id", id))
		return fmt.Errorf("task with ID %s not found", id)
	}

	// Purge the node from the graph (removes all edges and the node itself)
	m.purgeNode(id)

	// Update tag cache since a task was removed
	m.PopulateTagCache()
	m.logger.Info("Node deleted successfully",
		zap.String("task_id", id),
		zap.Int("remaining_nodes", len(m.nodes)),
	)

	return nil
}

// purgeNode removes a node from the graph and cleans up all edges pointing to it.
// This is an internal method used by DeleteNode and other operations.
// It does NOT validate that the node exists - caller must check.
func (m *Manager) purgeNode(id string) {
	m.logger.Debug("Purging node from graph", zap.String("node_id", id))

	// Track how many edges were cleaned up
	edgesRemoved := 0

	// Remove all edges pointing to this node from other nodes
	for _, node := range m.nodes {
		if node.ID == id {
			continue // Skip the node being deleted
		}

		// Iterate through all edge types in this node
		if node.EdgeIDs != nil {
			for relationshipName, targetIDs := range node.EdgeIDs {
				// Remove the deleted node's ID from this edge list
				cleaned := removeStringFromSlice(targetIDs, id)

				// Update if we removed anything
				if len(cleaned) != len(targetIDs) {
					node.EdgeIDs[relationshipName] = cleaned
					edgesRemoved++

					// Also clean up the resolved Edges map if it exists
					if node.Edges != nil {
						node.Edges[relationshipName] = removeEdgeByNodeID(node.Edges[relationshipName], id)
					}
				}
			}
		}
	}

	m.logger.Debug("Removed edges pointing to node",
		zap.String("node_id", id),
		zap.Int("edges_removed", edgesRemoved),
	)

	// Delete the node itself from the graph
	delete(m.nodes, id)

	m.logger.Debug("Node purged from graph", zap.String("node_id", id))
}

// removeStringFromSlice removes all occurrences of a string from a slice
func removeStringFromSlice(slice []string, value string) []string {
	if slice == nil {
		return nil
	}

	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}

	// Return nil if the result is empty to maintain nil vs empty slice distinction
	if len(result) == 0 && slice != nil {
		return []string{}
	}
	return result
}

// removeEdgeByNodeID removes all edges whose To node has the given ID
func removeEdgeByNodeID(edges []types.Edge, nodeID string) []types.Edge {
	if edges == nil {
		return nil
	}

	result := make([]types.Edge, 0, len(edges))
	for _, edge := range edges {
		if edge.To != nil && edge.To.ID != nodeID {
			result = append(result, edge)
		}
	}

	// Return nil if the result is empty to maintain nil vs empty slice distinction
	if len(result) == 0 && edges != nil {
		return []types.Edge{}
	}
	return result
}

// ListAllNodes returns all nodes in the manager
func (m *Manager) ListAllNodes() []*types.Node {
	nodes := make([]*types.Node, 0, len(m.nodes))
	for _, task := range m.nodes {
		nodes = append(nodes, task)
	}
	return nodes
}

// GetNode retrieves a task by ID
func (m *Manager) GetNode(id string) (*types.Node, error) {
	if id == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}

	task, exists := m.nodes[id]
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", id)
	}

	return task, nil
}

// getNodes retrieves multiple nodes by their IDs
func (m *Manager) getNodes(ids []string) ([]*types.Node, error) {
	if len(ids) == 0 {
		return []*types.Node{}, nil
	}

	nodes := make([]*types.Node, 0, len(ids))
	var notFound []string

	for _, id := range ids {
		if id == "" {
			return nil, fmt.Errorf("task ID cannot be empty")
		}

		task, exists := m.nodes[id]
		if !exists {
			notFound = append(notFound, id)
			continue
		}

		nodes = append(nodes, task)
	}

	if len(notFound) > 0 {
		return nodes, fmt.Errorf("nodes not found: %v", notFound)
	}

	return nodes, nil
}

// ResolveNodePointers populates the Edges map for all nodes by looking up the
// corresponding IDs in EdgeIDs and creating Edge objects with resolved pointers.
// Should be called after loading nodes from disk to restore the pointer relationships.
// Returns an error if any referenced node IDs cannot be found.
func (m *Manager) ResolveNodePointers() error {
	m.logger.Debug("Resolving node pointers for all nodes", zap.Int("node_count", len(m.nodes)))

	for _, node := range m.nodes {
		if node.EdgeIDs == nil {
			continue
		}

		// Initialize Edges map if needed
		if node.Edges == nil {
			node.Edges = make(map[string][]types.Edge)
		}

		// Iterate through all relationship types in this node
		for relationshipName, targetIDs := range node.EdgeIDs {
			if len(targetIDs) == 0 {
				continue
			}

			// Look up the target nodes
			targetNodes, err := m.getNodes(targetIDs)
			if err != nil {
				return fmt.Errorf("failed to resolve %s for node %s: %w", relationshipName, node.ID, err)
			}

			// Create Edge objects for each target
			edges := make([]types.Edge, len(targetNodes))
			for i, targetNode := range targetNodes {
				edges[i] = types.Edge{
					To: targetNode,
					// TODO: Populate Type field with the Relationship object once we implement
					// relationship configuration loading. For now, edges have resolved node
					// pointers which is sufficient for graph traversal.
					Type: nil,
				}
			}

			// Store the resolved edges
			node.Edges[relationshipName] = edges
		}
	}

	m.logger.Debug("Node pointers resolved successfully")
	return nil
}

// Clone creates a deep copy of the manager and all its nodes.
// The cloned manager has independent nodes with resolved pointers.
// This is useful for making transactional changes that can be validated before committing.
func (m *Manager) Clone() *Manager {
	if m == nil {
		return nil
	}

	m.logger.Debug("Cloning manager", zap.Int("task_count", len(m.nodes)))

	// Create new manager with same logger
	clone := &Manager{
		nodes:    make(map[string]*types.Node),
		tagCache: make(map[string][]*types.Node),
		logger:   m.logger,
	}

	// Clone all nodes
	for id, task := range m.nodes {
		clonedNode := task.Clone()
		clone.nodes[id] = clonedNode
	}

	// Resolve task pointers in the cloned manager
	// Note: We ignore errors here because if the original manager was valid,
	// the clone should also be valid. If there are resolution errors, they
	// would have existed in the original manager too.
	_ = clone.ResolveNodePointers()

	// Clone tag cache (we'll just rebuild it)
	clone.PopulateTagCache()

	m.logger.Debug("Manager cloned successfully")

	return clone
}

// DetectCycles checks all relationship types for cycles. Since all relationships form DAGs,
// any cycle in any relationship type is invalid. Returns an error if any cycles are detected,
// with detailed information about all cycles found.
func (m *Manager) DetectCycles() error {
	m.logger.Debug("Detecting cycles in graph")

	var allCycles []string

	// Collect all unique relationship types across all nodes
	relationshipTypes := make(map[string]bool)
	for _, node := range m.nodes {
		if node.EdgeIDs != nil {
			for relationshipName := range node.EdgeIDs {
				relationshipTypes[relationshipName] = true
			}
		}
	}

	m.logger.Debug("Checking relationship types for cycles",
		zap.Int("relationship_count", len(relationshipTypes)))

	// Check each relationship type for cycles
	for relationshipName := range relationshipTypes {
		relName := relationshipName // Capture for closure
		cycles := m.detectCyclesInDAG(relName, func(node *types.Node) []string {
			if node.EdgeIDs == nil {
				return nil
			}
			return node.EdgeIDs[relName]
		})

		if len(cycles) > 0 {
			for _, cycle := range cycles {
				allCycles = append(allCycles, fmt.Sprintf("%s: %s", relName, cycle))
			}
		}
	}

	if len(allCycles) > 0 {
		msg := fmt.Sprintf("detected %d cycle(s):\n", len(allCycles))
		for i, cycle := range allCycles {
			msg += fmt.Sprintf("  %d. %s\n", i+1, cycle)
		}
		m.logger.Error("Cycles detected in graph", zap.Int("cycle_count", len(allCycles)))
		return fmt.Errorf("%s", msg)
	}

	return nil
}

// detectCyclesInDAG performs cycle detection on a specific DAG using DFS
// Returns a slice of cycle descriptions (e.g., "task-a -> task-b -> task-c -> task-a")
func (m *Manager) detectCyclesInDAG(dagName string, getEdges func(*types.Node) []string) []string {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var cycles []string
	path := []string{}

	// Check each task as a potential starting point
	for taskID := range m.nodes {
		if !visited[taskID] {
			m.findCyclesDFS(taskID, visited, recStack, &path, &cycles, getEdges)
		}
	}

	return cycles
}

// findCyclesDFS performs depth-first search to find all cycles
func (m *Manager) findCyclesDFS(taskID string, visited, recStack map[string]bool, path *[]string, cycles *[]string, getEdges func(*types.Node) []string) {
	// Mark current node as visited and add to recursion stack
	visited[taskID] = true
	recStack[taskID] = true
	*path = append(*path, taskID)

	// Get the task
	task, exists := m.nodes[taskID]
	if exists {
		// Get edges for this task based on the DAG we're checking
		edges := getEdges(task)

		// Recursively check all adjacent nodes
		for _, adjacentID := range edges {
			// If adjacent node is not visited, recurse on it
			if !visited[adjacentID] {
				m.findCyclesDFS(adjacentID, visited, recStack, path, cycles, getEdges)
			} else if recStack[adjacentID] {
				// If adjacent node is in recursion stack, we found a cycle
				// Find where the cycle starts in the path
				cycleStart := -1
				for i, id := range *path {
					if id == adjacentID {
						cycleStart = i
						break
					}
				}

				// Build the cycle description
				if cycleStart >= 0 {
					cyclePath := append((*path)[cycleStart:], adjacentID)
					cycleDesc := ""
					for i, id := range cyclePath {
						if i > 0 {
							cycleDesc += " -> "
						}
						cycleDesc += id
					}
					*cycles = append(*cycles, cycleDesc)
				}
			}
		}
	}

	// Remove from recursion stack and path before returning
	recStack[taskID] = false
	*path = (*path)[:len(*path)-1]
}

// LoadFromDir reads all YAML files from the specified directory and loads nodes
func (m *Manager) LoadFromDir(dirPath string) error {
	m.logger.Info("Loading nodes from directory", zap.String("path", dirPath))

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		m.logger.Error("Failed to create directory", zap.String("path", dirPath), zap.Error(err))
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Read all .yaml files in the directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		m.logger.Error("Failed to read directory", zap.String("path", dirPath), zap.Error(err))
		return fmt.Errorf("failed to read directory: %w", err)
	}

	m.logger.Debug("Found directory entries", zap.Int("count", len(entries)))

	nodesLoaded := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" {
			m.logger.Debug("Skipping non-YAML file", zap.String("filename", entry.Name()))
			continue
		}

		filename := filepath.Join(dirPath, entry.Name())
		m.logger.Debug("Reading task file", zap.String("filename", filename))

		data, err := os.ReadFile(filename)
		if err != nil {
			m.logger.Error("Failed to read file", zap.String("filename", entry.Name()), zap.Error(err))
			return fmt.Errorf("failed to read file %s: %w", entry.Name(), err)
		}

		var task types.Node
		if err := yaml.Unmarshal(data, &task); err != nil {
			m.logger.Error("Failed to unmarshal task",
				zap.String("filename", entry.Name()),
				zap.Error(err),
			)
			return fmt.Errorf("failed to unmarshal task from %s: %w", entry.Name(), err)
		}

		m.nodes[task.ID] = &task
		nodesLoaded++
		m.logger.Debug("Loaded task from file",
			zap.String("task_id", task.ID),
			zap.String("task_name", task.Name),
			zap.String("filename", entry.Name()),
		)
	}

	m.logger.Info("Finished loading task files", zap.Int("nodes_loaded", nodesLoaded))

	// Detect cycles before resolving pointers
	m.logger.Debug("Detecting cycles in task graph")
	if err := m.DetectCycles(); err != nil {
		m.logger.Error("Cycle detected in task graph", zap.Error(err))
		return fmt.Errorf("cycle detected in task graph: %w", err)
	}
	m.logger.Debug("No cycles detected")

	// Resolve task pointers after loading all nodes and validating no cycles
	m.logger.Debug("Resolving task pointers")
	if err := m.ResolveNodePointers(); err != nil {
		m.logger.Error("Failed to resolve task pointers", zap.Error(err))
		return err
	}
	m.logger.Debug("Node pointers resolved")

	// Populate tag cache for efficient tag-based lookups
	m.logger.Debug("Populating tag cache")
	m.PopulateTagCache()

	m.logger.Info("Successfully loaded nodes from directory",
		zap.String("path", dirPath),
		zap.Int("total_nodes", len(m.nodes)),
	)

	return nil
}

// PersistToDir writes all nodes to the specified directory as YAML files
func (m *Manager) PersistToDir(dirPath string) error {
	m.logger.Info("Persisting nodes to directory",
		zap.String("path", dirPath),
		zap.Int("task_count", len(m.nodes)),
	)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		m.logger.Error("Failed to create directory", zap.String("path", dirPath), zap.Error(err))
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write each task as a separate YAML file
	nodesPersisted := 0
	for id, task := range m.nodes {
		filename := filepath.Join(dirPath, fmt.Sprintf("%s.yaml", id))

		m.logger.Debug("Marshaling task", zap.String("task_id", id))
		data, err := yaml.Marshal(task)
		if err != nil {
			m.logger.Error("Failed to marshal task", zap.String("task_id", id), zap.Error(err))
			return fmt.Errorf("failed to marshal task %s: %w", id, err)
		}

		m.logger.Debug("Writing task file", zap.String("filename", filename))
		if err := os.WriteFile(filename, data, 0644); err != nil {
			m.logger.Error("Failed to write task file",
				zap.String("task_id", id),
				zap.String("filename", filename),
				zap.Error(err),
			)
			return fmt.Errorf("failed to write task %s: %w", id, err)
		}

		nodesPersisted++
	}

	m.logger.Info("Successfully persisted nodes to directory",
		zap.String("path", dirPath),
		zap.Int("nodes_persisted", nodesPersisted),
	)

	return nil
}
