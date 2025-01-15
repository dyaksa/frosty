package workflow

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func RollbackWorkflow(db *sql.DB, workflowID uuid.UUID, currentNodeID uuid.UUID, rollbackScope RollbackScope) error {
	fmt.Printf("Starting rollback for workflow %s from node %s with scope %s\n", workflowID, currentNodeID, rollbackScope)

	switch rollbackScope {
	case RollbackToStart:
		return rollbackToStart(db, workflowID, currentNodeID)
	case RollbackOne:
		return rollbackOneAncestor(db, workflowID, currentNodeID)
	case RollbackFinish:
		return rollbackFinish(db, workflowID, currentNodeID)
	default:
		return fmt.Errorf("unknown rollback scope: %s", rollbackScope)
	}
}

// Rollback to the very start of the workflow
func rollbackToStart(db *sql.DB, workflowID uuid.UUID, currentNodeID uuid.UUID) error {
	fmt.Printf("Rolling back to the start of the workflow %s\n", workflowID)

	// Traverse back to the starting node
	for {
		node, err := GetNode(db, currentNodeID)
		if err != nil {
			return fmt.Errorf("error retrieving node %s: %v", currentNodeID, err)
		}

		// Perform rollback logic for the current node
		err = RollbackNode(db, node.ID)
		if err != nil {
			return fmt.Errorf("rollback failed for node %s: %v", node.ID, err)
		}

		if node.Type == "Start" {
			// Reached the start node
			fmt.Printf("Reached the starting node %s\n", node.ID)
			break
		}

		// Get the ancestor node
		ancestor, err := GetImmediateAncestor(db, node.ID)
		if err != nil {
			return fmt.Errorf("error retrieving ancestor for node %s: %v", node.ID, err)
		}
		currentNodeID = ancestor.ID
	}

	fmt.Printf("Rollback to start completed for workflow %s\n", workflowID)
	return nil
}

// Rollback one ancestor node
func rollbackOneAncestor(db *sql.DB, workflowID uuid.UUID, currentNodeID uuid.UUID) error {
	fmt.Printf("Rolling back one ancestor for workflow %s from node %s\n", workflowID, currentNodeID)

	// Get the current node
	node, err := GetNode(db, currentNodeID)
	if err != nil {
		return fmt.Errorf("error retrieving node %s: %v", currentNodeID, err)
	}

	// Perform rollback logic for the current node
	err = RollbackNode(db, node.ID)
	if err != nil {
		return fmt.Errorf("rollback failed for node %s: %v", node.ID, err)
	}

	// Get the ancestor node
	ancestor, err := GetImmediateAncestor(db, node.ID)
	if err != nil {
		return fmt.Errorf("error retrieving ancestor for node %s: %v", node.ID, err)
	}

	// Rollback the ancestor
	err = RollbackNode(db, ancestor.ID)
	if err != nil {
		return fmt.Errorf("rollback failed for ancestor node %s: %v", ancestor.ID, err)
	}

	fmt.Printf("Rollback one ancestor completed for workflow %s\n", workflowID)
	return nil
}

// Stop/finish the workflow upon rollback
func rollbackFinish(db *sql.DB, workflowID uuid.UUID, currentNodeID uuid.UUID) error {
	fmt.Printf("Finishing workflow %s from node %s\n", workflowID, currentNodeID)

	// Perform rollback logic for the current node
	node, err := GetNode(db, currentNodeID)
	if err != nil {
		return fmt.Errorf("error retrieving node %s: %v", currentNodeID, err)
	}

	err = RollbackNode(db, node.ID)
	if err != nil {
		return fmt.Errorf("rollback failed for node %s: %v", currentNodeID, err)
	}

	// Mark the workflow as "finished"
	err = UpdateWorkflowStatus(db, workflowID, "finished")
	if err != nil {
		return fmt.Errorf("failed to mark workflow %s as finished: %v", workflowID, err)
	}

	fmt.Printf("Workflow %s has been finished\n", workflowID)
	return nil
}

// Helper function to perform rollback on a single node
func RollbackNode(db *sql.DB, nodeID uuid.UUID) error {
	fmt.Printf("Rolling back node %s\n", nodeID)

	// Perform rollback tasks for the node
	tasks, err := GetNodeTasks(db, nodeID)
	if err != nil {
		return fmt.Errorf("error retrieving tasks for node %s: %v", nodeID, err)
	}

	for _, task := range tasks {
		err := rollbackTask(db, task.ID)
		if err != nil {
			return fmt.Errorf("rollback failed for task %s: %v", task.ID, err)
		}
	}

	fmt.Printf("Rollback completed for node %s\n", nodeID)
	return nil
}

// Helper function to rollback a single task
func rollbackTask(db *sql.DB, taskID uuid.UUID) error {
	fmt.Printf("Rolling back task %s\n", taskID)

	// Logic to revert task (could involve deleting records, undoing changes, etc.)
	// Example: Mark task as reverted or execute rollback actions
	err := UpdateTaskStatus(db, taskID, "reverted")
	if err != nil {
		return fmt.Errorf("failed to mark task %s as reverted: %v", taskID, err)
	}

	fmt.Printf("Task %s rolled back successfully\n", taskID)
	return nil
}
