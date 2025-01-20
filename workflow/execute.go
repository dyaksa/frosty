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
	err = LogWorkflowNodeExecution(db, workflowID, startNode.ID, "executing", "Starting workflow execution")
	if err != nil {
		return fmt.Errorf("failed to log workflow execution: %v", err)
	}

	// Initialize a queue for BFS traversal
	queue := []uuid.UUID{startNode.ID}
	visited := make(map[uuid.UUID]bool)

	// Execute nodes using BFS
	for len(queue) > 0 {
		currentNodeID := queue[0]
		queue = queue[1:]

		if visited[currentNodeID] {
			continue
		}

		visited[currentNodeID] = true

		err = ExecuteNode(db, currentNodeID, workflowID, &queue, visited)
		if err != nil {
			return fmt.Errorf("workflow execution failed: %v", err)
		}
	}

	return nil
}

func ExecuteNode(db *sql.DB, nodeID uuid.UUID, workflowID uuid.UUID, queue *[]uuid.UUID, visited map[uuid.UUID]bool) error {
	fmt.Printf("Executing node %s\n", nodeID)

	// Retrieve nodeTask associated with the node
	nodeTask, err := GetNodeTasks(db, nodeID)
	if err != nil {
		return fmt.Errorf("error retrieving tasks for node %s: %v", nodeID, err)
	}

	// Execute each task in sequence
	for _, task := range nodeTask {
		err := ExecuteTask(db, task.Task, task.RetryCount)
		if err != nil {
			return fmt.Errorf("task %s execution failed in node %s: %v", task.ID, nodeID, err)
		}
	}

	fmt.Printf("All tasks in node %s completed successfully\n", nodeID)

	// Log workflow node execution
	err = LogWorkflowNodeExecution(db, workflowID, nodeID, "completed", "Node execution completed")
	if err != nil {
		return fmt.Errorf("failed to log workflow node execution: %v", err)
	}

	// Get the next node(s) to execute
	descendants, err := GetDescendants(db, nodeID)
	if err != nil {
		return fmt.Errorf("error retrieving next nodes for node %s: %v", nodeID, err)
	}

	// Add descendants to the queue if not visited
	for _, descendant := range descendants {
		if !visited[descendant.ID] {
			*queue = append(*queue, descendant.ID)
		}
	}

	return nil
}

func ExecuteTask(db *sql.DB, task Task, retryCount int) error {
	fmt.Printf("Executing task %s\n", task.Title)

	// Retry logic
	retryLimit := 3
	for retry := 0; retry <= retryLimit; retry++ {
		// Simulate task execution
		err, response, httpCode := performTask(task)

		// Handle success
		if err == nil {
			fmt.Printf("Task %s executed successfully on attempt %d\n", task.ID, retry+1)

			// Update task status and response
			err := UpdateTaskStatusAndResponse(db, task.ID, "completed", response, "", httpCode)
			if err != nil {
				return fmt.Errorf("failed to update task %s status: %v", task.ID, err)
			}

			return nil
		}

		// Handle retry
		fmt.Printf("Task %s failed on attempt %d: %v\n", task.ID, retry+1, err)

		// Update task retry count and error
		err = UpdateTaskStatusAndResponse(db, task.ID, "failed", "", err.Error(), httpCode)
		if err != nil {
			return fmt.Errorf("failed to update task %s status on failure: %v", task.ID, err)
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

func performTask(task Task) (error, string, int) {
	// Create the HTTP request
	req, err := http.NewRequest(task.HttpMethod, task.Action, bytes.NewBuffer([]byte(task.Params)))
	if err != nil {
		return err, "", http.StatusInternalServerError
	}

	// Set headers if needed
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err, "", http.StatusInternalServerError
	}
	defer resp.Body.Close()

	// Read the response body
	var responseBody bytes.Buffer
	_, err = responseBody.ReadFrom(resp.Body)
	if err != nil {
		return err, "", http.StatusInternalServerError
	}

	// Handle the response
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices {
		return nil, responseBody.String(), resp.StatusCode
	}

	// Handle errors
	return fmt.Errorf("Task %s failed with status code: %d", task.Title, resp.StatusCode), "", resp.StatusCode
}

func evaluateCondition(node Node, child Node) bool {
	// Example: Evaluate based on some attributes or external data
	return true // Replace with actual condition logic
}
