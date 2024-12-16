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
			if !allParentsCompleted(db, currentID) { // Implement this function
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

		default:
			return fmt.Errorf("unsupported node type: %s", node.Type)
		}
	}
	return nil
}

func ValidateWorkflow(db *sql.DB, startNode uuid.UUID) error {
	rows, err := db.Query(`
		SELECT COUNT(*)
		FROM node_closure
		WHERE ancestor = descendant AND ancestor = ?
	`, startNode)
	if err != nil {
		return err
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return err
		}
		if count > 1 {
			return fmt.Errorf("cyclic dependency detected")
		}
	}
	return nil
}

func evaluateCondition(node Node, child Node) bool {
	// Example: Evaluate based on some attributes or external data
	return true // Replace with actual condition logic
}

func allParentsCompleted(db *sql.DB, nodeID uuid.UUID) bool {
	var count int
	err := db.QueryRow(`
        SELECT COUNT(*)
        FROM node_closure nc
        JOIN nodes n ON nc.ancestor = n.id
        WHERE nc.descendant = ? AND n.type != 'End'
    `, nodeID).Scan(&count)

	if err != nil {
		return false
	}
	return count == 0
}
