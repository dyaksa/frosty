package workflow

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func logWorkflowNode(db *sql.DB, workflowID uuid.UUID, nodeID uuid.UUID, status string, message string) error {
	// Log workflow execution
	err := LogWorkflowExecution(db, workflowID, nodeID, nil, status, message, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to log workflow execution: %v", err)
	}

	return nil
}

func logWorkflowNodeTask(db *sql.DB, workflowID uuid.UUID, nodeID uuid.UUID, taskID uuid.UUID, status string, message string) error {
	// Log workflow execution
	err := LogWorkflowExecution(db, workflowID, nodeID, &taskID, status, message, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to log workflow execution: %v", err)
	}

	return nil
}
