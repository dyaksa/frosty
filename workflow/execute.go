package workflow

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func ExecuteWorkflow(db *sql.DB, workflowID uuid.UUID) error {
	// Fetch starting node
	startNode, err := GetStartingNode(db, workflowID)
	if err != nil {
		return fmt.Errorf("failed to fetch starting node: %v", err)
	}

	// Update workflow status
	err = UpdateWorkflowStatus(db, workflowID, "executing")
	if err != nil {
		return fmt.Errorf("failed to update workflow status: %v", err)
	}

	// Log workflow execution
	err = LogWorkflowExecution(db, workflowID, startNode.ID, nil, "executing", "Starting workflow execution", nil, nil, err)
	if err != nil {
		return fmt.Errorf("failed to log workflow execution: %v", err)
	}

	// Initialize a queue for BFS traversal
	nodeQueue := []uuid.UUID{startNode.ID}
	visited := make(map[uuid.UUID]bool)

	// Execute nodes using BFS
	for len(nodeQueue) > 0 {
		currentNodeID := nodeQueue[0]
		nodeQueue = nodeQueue[1:]

		if visited[currentNodeID] {
			continue
		}

		visited[currentNodeID] = true

		err = ExecuteNode(db, currentNodeID, workflowID, &nodeQueue, visited)
		if err != nil {
			UpdateWorkflowStatus(db, workflowID, "error")
			return fmt.Errorf("workflow execution failed: %v", err)
		}
	}

	// Update workflow status
	err = UpdateWorkflowStatus(db, workflowID, "completed")
	if err != nil {
		return fmt.Errorf("failed to update workflow status: %v", err)
	}

	// Log workflow execution
	err = LogWorkflowExecution(db, workflowID, startNode.ID, nil, "completed", "Starting workflow execution", nil, nil, err)
	if err != nil {
		return fmt.Errorf("failed to log workflow execution: %v", err)
	}

	return nil
}

func ExecuteNode(db *sql.DB, nodeID uuid.UUID, workflowID uuid.UUID, queue *[]uuid.UUID, visited map[uuid.UUID]bool) error {
	fmt.Printf("Executing node %s\n", nodeID)

	// Retrieve nodeTask associated with the node
	nodeTasks, err := GetNodeTasks(db, nodeID)

	if err != nil {
		return fmt.Errorf("error retrieving tasks for node %s: %v", nodeID, err)
	}

	fmt.Printf("Node %s has %d tasks\n", nodeID, len(nodeTasks))
	fmt.Printf("Executing tasks in node %s\n", nodeID)
	fmt.Printf("Tasks: %v\n", nodeTasks)

	// Execute each task in sequence
	for _, nodeTask := range nodeTasks {
		err := ExecuteTask(db, workflowID, nodeID, nodeTask.Task, nodeTask.RetryCount)
		if err != nil {
			return fmt.Errorf("task %s execution failed in node %s: %v", nodeTask.ID, nodeID, err)
		}

		// Log workflow execution
		err = LogWorkflowExecution(db, workflowID, nodeID, &nodeTask.TaskID, "completed", fmt.Sprintf("Task %s execution completed", nodeTask.TaskID), nil, nil, err)
		if err != nil {
			return fmt.Errorf("failed to log task execution: %v", err)
		}
	}

	fmt.Printf("All tasks in node %s completed successfully\n", nodeID)

	// Log workflow node execution
	err = LogWorkflowExecution(db, workflowID, nodeID, nil, "completed", "Node execution completed", nil, nil, err)
	if err != nil {
		return fmt.Errorf("failed to log workflow node execution: %v", err)
	}

	// Get the next node(s) to execute
	descendants, err := GetDescendants(db, nodeID)
	if err != nil {
		return fmt.Errorf("error retrieving next nodes for node %s: %v", nodeID, err)
	}

	// Evaluate condition for next node execution
	nextNodeID, err := evaluateCondition(db, nodeID, descendants)
	if err != nil {
		return fmt.Errorf("error evaluating condition for next node execution: %v", err)
	}

	// Add the next node to the queue if not visited
	if nextNodeID != uuid.Nil && !visited[nextNodeID] {
		*queue = append(*queue, nextNodeID)
	}

	return nil
}

func ExecuteTask(db *sql.DB, workflowID uuid.UUID, nodeID uuid.UUID, task Task, retryCount int) error {
	fmt.Printf("Executing task %s\n", task.Title)

	// Retry logic
	retryLimit := task.MaxRetries
	for retry := 0; retry <= retryLimit; retry++ {
		// Simulate task execution
		response, httpCode, err := performTask(task)

		// Handle success
		if err == nil {
			fmt.Printf("Task %s executed successfully on attempt %d\n", task.ID, retry+1)

			// Update task status and response
			err := UpdateTaskStatus(db, task.ID, "completed", retryCount)
			if err != nil {
				return fmt.Errorf("failed to update task %s status: %v", task.ID, err)
			}

			err = LogWorkflowExecution(db, workflowID, nodeID, &task.ID, "completed", fmt.Sprintf("Task \"%s\" performed successfully: [%d] %s", task.Title, httpCode, response), &httpCode, &response, err)
			if err != nil {
				return fmt.Errorf("failed to log workflow node execution: %v", err)
			}

			return nil
		}

		// Handle retry
		fmt.Printf("Task %s failed on attempt %d: %v\n", task.ID, retry+1, err)

		// Update task retry count and error
		err = UpdateTaskStatus(db, task.ID, "failed", retryCount)
		if err != nil {
			return fmt.Errorf("failed to update task %s status on failure: %v", task.ID, err)
		}

		err = LogWorkflowExecution(db, workflowID, nodeID, &task.ID, "failed", fmt.Sprintf("Task \"%s\" (retry no: %d) performed failed: [%d] %s", task.Title, retry, httpCode, response), &httpCode, &response, err)
		if err != nil {
			return fmt.Errorf("failed to log workflow node execution: %v", err)
		}

		if retry < retryLimit {
			fmt.Printf("Retrying task %s (%d/%d)\n", task.ID, retry+2, retryLimit+1)
		} else {
			break
		}
	}

	// Task failed after retries
	fmt.Printf("Task %s failed after %d attempts\n", task.ID, retryLimit+1)
	return fmt.Errorf("task %s execution failed after maximum retries", task.ID)
}

func performTask(task Task) (string, int, error) {
	// Create the HTTP request
	req, err := http.NewRequest(task.HttpMethod, task.Action, bytes.NewBuffer([]byte(task.Params)))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	// Set headers if needed
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	// Read the response body
	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	// Handle the response
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return responseBody.String(), resp.StatusCode, nil
	}

	// Handle errors
	return "", resp.StatusCode, fmt.Errorf("Task %s failed with status code: %d", task.Title, resp.StatusCode)
}

func evaluateCondition(db *sql.DB, currentNodeID uuid.UUID, descendants []Node) (uuid.UUID, error) {
	// Check the status of the current node tasks
	nodeTasks, err := GetNodeTasks(db, currentNodeID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error retrieving tasks for node %s: %v", currentNodeID, err)
	}

	// If all tasks are completed successfully, execute the next node
	allTasksCompleted := true
	for _, nodeTask := range nodeTasks {
		if nodeTask.Status != "completed" {
			allTasksCompleted = false
			break
		}
	}

	// Return the next node if all tasks are completed for the current node
	if allTasksCompleted && len(descendants) > 0 {
		return descendants[0].ID, nil
	}

	// If no condition matches, return nil
	return uuid.Nil, nil
}

func ExecuteWorkflowByExecutionID(db *sql.DB, executionID uuid.UUID) error {
	// Fetch the workflow execution details
	execution, err := GetWorkflowExecutionByID(db, executionID)
	if err != nil {
		return fmt.Errorf("failed to fetch workflow execution: %v", err)
	}

	// Update workflow execution status
	err = UpdateWorkflowExecutionStatus(db, executionID, "executing")
	if err != nil {
		return fmt.Errorf("failed to update workflow execution status: %v", err)
	}

	// Fetch starting node
	startNode, err := GetStartingNode(db, execution.WorkflowID)
	if err != nil {
		return fmt.Errorf("failed to fetch starting node: %v", err)
	}

	// Log workflow execution
	err = LogWorkflowExecution(db, execution.WorkflowID, startNode.ID, nil, "executing", "Starting workflow execution", nil, nil, err)
	if err != nil {
		return fmt.Errorf("failed to log workflow execution: %v", err)
	}

	// Initialize a queue for BFS traversal
	nodeQueue := []uuid.UUID{startNode.ID}
	visited := make(map[uuid.UUID]bool)

	// Execute nodes using BFS
	for len(nodeQueue) > 0 {
		currentNodeID := nodeQueue[0]
		nodeQueue = nodeQueue[1:]

		if visited[currentNodeID] {
			continue
		}

		visited[currentNodeID] = true

		err = ExecuteNode(db, currentNodeID, execution.WorkflowID, &nodeQueue, visited)
		if err != nil {
			UpdateWorkflowExecutionStatus(db, executionID, "error")
			return fmt.Errorf("workflow execution failed: %v", err)
		}
	}

	// Update workflow execution status
	err = UpdateWorkflowExecutionStatus(db, executionID, "completed")
	if err != nil {
		return fmt.Errorf("failed to update workflow execution status: %v", err)
	}

	// Log workflow execution
	err = LogWorkflowExecution(db, execution.WorkflowID, startNode.ID, nil, "completed", "Starting workflow execution", nil, nil, err)
	if err != nil {
		return fmt.Errorf("failed to log workflow execution: %v", err)
	}

	return nil
}
