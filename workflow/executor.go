package workflow

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

func ExecuteWorkflow(db *sql.DB, startNode uuid.UUID, action func(Node) error) error {
	queue := []uuid.UUID{startNode}
	visited := make(map[uuid.UUID]bool)

	for len(queue) > 0 {
		currentID := queue[0]
		queue = queue[1:]

		if visited[currentID] {
			continue
		}
		visited[currentID] = true

		node, err := GetNode(db, currentID)
		if err != nil {
			return err
		}

		// Perform the action at this node
		err = action(node)
		if err != nil {
			return err
		}

		switch node.Type {
		case NodeTypeTask:
			descendants, err := GetDescendants(db, currentID)
			if err != nil {
				return err
			}
			for _, child := range descendants {
				queue = append(queue, child.ID)
			}

		case NodeTypeDecision:
			// Custom logic for decision nodes
			descendants, err := GetDescendants(db, currentID)
			if err != nil {
				return err
			}
			for _, child := range descendants {
				conditionMet := evaluateCondition(node, child) // Implement this function
				if conditionMet {
					queue = append(queue, child.ID)
					break
				}
			}

		case NodeTypeFork:
			descendants, err := GetDescendants(db, currentID)
			if err != nil {
				return err
			}
			var wg sync.WaitGroup
			for _, child := range descendants {
				wg.Add(1)
				go func(childID uuid.UUID) {
					defer wg.Done()
					ExecuteWorkflow(db, childID, action) // Recursive execution
				}(child.ID)
			}
			wg.Wait()

		case NodeTypeJoin:
			// Wait for all parent nodes to complete
			if !AllParentsCompleted(db, currentID) { // Implement this function
				queue = append(queue, currentID)
				continue
			}
			descendants, err := GetDescendants(db, currentID)
			if err != nil {
				return err
			}
			for _, child := range descendants {
				queue = append(queue, child.ID)
			}

		case NodeTypeEnd:
			// End nodes don't have children
			continue

		case NodeTypeStart:
			// Start node, do nothing and continue
			continue

		default:
			return fmt.Errorf("unsupported node type: %s", node.Type)
		}

		err = LogNodeExecution(db, node.ID, "success", "Node executed successfully")
		if err != nil {
			return err
		}
	}
	return nil
}

func RollbackWorkflow(db *sql.DB, currentNode uuid.UUID, rollbackScope string, rollbackAction func(Node) error) error {
	// Get executed nodes based on rollback scope
	var nodesToRollback []Node
	var err error

	switch rollbackScope {
	case "ancestor":
		// Get only the immediate ancestor of the current node
		nodesToRollback, err = GetImmediateAncestors(db, currentNode)
		if err != nil {
			return fmt.Errorf("failed to retrieve immediate ancestor: %w", err)
		}

	case "start":
		// Get all executed nodes up to the starting node
		nodesToRollback, err = GetExecutedNodes(db, currentNode)
		if err != nil {
			return fmt.Errorf("failed to retrieve executed nodes: %w", err)
		}

	case "finish":
		// Get the current node
		existingCurrentNode, err := GetNode(db, currentNode)
		if err != nil {
			return fmt.Errorf("failed to retrieve node: %w", err)
		}
		// Log rollback success
		err = LogNodeExecution(db, currentNode, "rollback", "Node rolled back successfully")
		if err != nil {
			return fmt.Errorf("failed to log rollback for node %s: %w", existingCurrentNode.Title, err)
		}

		// Return early as no further rollback actions are necessary
		return nil

	default:
		return fmt.Errorf("invalid rollback scope: %s", rollbackScope)
	}

	// Traverse nodes in reverse order for rollback
	for i := len(nodesToRollback) - 1; i >= 0; i-- {
		node := nodesToRollback[i]
		fmt.Printf("Rolling back node: %s (%s)\n", node.Title, node.Type)

		err := rollbackAction(node)
		if err != nil {
			return fmt.Errorf("failed to rollback node %s: %w", node.Title, err)
		}

		// Log rollback success
		err = LogNodeExecution(db, node.ID, "rollback", "Node rolled back successfully")
		if err != nil {
			return fmt.Errorf("failed to log rollback for node %s: %w", node.Title, err)
		}
	}

	fmt.Println("Rollback completed successfully")
	return nil
}

func evaluateCondition(node Node, child Node) bool {
	// Example: Evaluate based on some attributes or external data
	return true // Replace with actual condition logic
}
