package workflow

import (
	"fmt"

	"github.com/google/uuid"
)

func ExecuteWorkflow(workflowID uuid.UUID) error {
	// Fetch all nodes in the workflow
	nodes, err := GetWorkflowNodes(workflowID)
	if err != nil {
		return fmt.Errorf("failed to fetch workflow nodes: %v", err)
	}

	// Fetch starting node
	startNode, err := GetStartingNode(workflowID)
	if err != nil {
		return fmt.Errorf("failed to fetch starting node: %v", err)
	}

	// Execute nodes recursively
	err = ExecuteNode(startNode.ID)
	if err != nil {
		return fmt.Errorf("workflow execution failed: %v", err)
	}

	return nil
}

func ExecuteNode(nodeID uuid.UUID) error {
	// Fetch node details
	node, err := GetNode(nodeID)
	if err != nil {
		return fmt.Errorf("failed to fetch node %v: %v", nodeID, err)
	}

	// Execute tasks associated with the node
	tasks, err := GetNodeTasks(nodeID)
	if err != nil {
		return fmt.Errorf("failed to fetch tasks for node %v: %v", nodeID, err)
	}

	for _, task := range tasks {
		err = ExecuteTask(task.ID)
		if err != nil {
			// Handle task failure logic (e.g., retries or rollback)
			return fmt.Errorf("task %v failed: %v", task.ID, err)
		}
	}

	// Fetch descendants and execute them recursively
	descendants, err := GetDescendantNodes(nodeID)
	if err != nil {
		return fmt.Errorf("failed to fetch descendants for node %v: %v", nodeID, err)
	}

	for _, descendant := range descendants {
		err = ExecuteNode(descendant.ID)
		if err != nil {
			return fmt.Errorf("execution failed for descendant %v: %v", descendant.ID, err)
		}
	}

	return nil
}

func ExecuteNode(nodeID uuid.UUID) error {
	fmt.Printf("Executing node %s\n", nodeID)

	// Retrieve tasks associated with the node
	tasks, err := GetNodeTasks(nodeID)
	if err != nil {
		return fmt.Errorf("error retrieving tasks for node %s: %v", nodeID, err)
	}

	// Execute each task in sequence
	for _, task := range tasks {
		err := ExecuteTask(task.ID)
		if err != nil {
			return fmt.Errorf("task %s execution failed in node %s: %v", task.ID, nodeID, err)
		}
	}

	fmt.Printf("All tasks in node %s completed successfully\n", nodeID)

	// Get the next node(s) to execute
	nextNodes, err := GetNextNodes(nodeID)
	if err != nil {
		return fmt.Errorf("error retrieving next nodes for node %s: %v", nodeID, err)
	}

	// Execute the next node(s)
	for _, nextNode := range nextNodes {
		err := ExecuteNode(nextNode.ID)
		if err != nil {
			return fmt.Errorf("execution failed for next node %s: %v", nextNode.ID, err)
		}
	}

	return nil
}

func ExecuteTask(taskID uuid.UUID) error {
	fmt.Printf("Executing task %s\n", taskID)

	// Retrieve task details
	task, err := GetTask(taskID)
	if err != nil {
		return fmt.Errorf("error retrieving task %s: %v", taskID, err)
	}

	// Retry logic
	retryLimit := 3
	for retry := 0; retry <= retryLimit; retry++ {
		// Simulate task execution
		err, response, httpCode := performTask(task)

		// Handle success
		if err == nil {
			fmt.Printf("Task %s executed successfully on attempt %d\n", taskID, retry+1)

			// Update task status and response
			err := UpdateTaskStatusAndResponse(taskID, "completed", response, httpCode, "")
			if err != nil {
				return fmt.Errorf("failed to update task %s status: %v", taskID, err)
			}

			return nil
		}

		// Handle retry
		fmt.Printf("Task %s failed on attempt %d: %v\n", taskID, retry+1, err)

		// Update task retry count and error
		err = UpdateTaskStatusAndResponse(taskID, "failed", "", httpCode, err.Error())
		if err != nil {
			return fmt.Errorf("failed to update task %s status on failure: %v", taskID, err)
		}

		if retry < retryLimit {
			fmt.Printf("Retrying task %s (%d/%d)\n", taskID, retry+2, retryLimit+1)
		} else {
			break
		}
	}

	// Task failed after retries
	fmt.Printf("Task %s failed after %d attempts\n", taskID, retryLimit+1)
	return fmt.Errorf("task %s execution failed after maximum retries", taskID)
}

func evaluateCondition(node Node, child Node) bool {
	// Example: Evaluate based on some attributes or external data
	return true // Replace with actual condition logic
}
