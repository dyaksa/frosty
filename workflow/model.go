package workflow

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Node struct {
	ID            uuid.UUID  `db:"id" json:"id"`                                             // Unique identifier for the node
	Title         string     `db:"title" json:"title"`                                       // Title of the node
	Type          string     `db:"type" json:"type"`                                         // Type of the node
	Description   string     `db:"description,omitempty" json:"description,omitempty"`       // Description of the node
	RollbackScope string     `db:"rollback_scope,omitempty" json:"rollback_scope,omitempty"` // Rollback scope for the node: start, finish, or immediate_ancestor
	CreatedAt     *time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`         // Creation timestamp
	UpdatedAt     *time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`         // Update timestamp
	DeletedAt     *time.Time `db:"deleted_at,omitempty" json:"deleted_at,omitempty"`         // Deletion timestamp
}

type NodeClosure struct {
	Ancestor   uuid.UUID `db:"ancestor" json:"ancestor"`     // Unique identifier for the ancestor node
	Descendant uuid.UUID `db:"descendant" json:"descendant"` // Unique identifier for the descendant node
	Depth      int       `db:"depth" json:"depth"`           // Depth of the descendant node
}

type NodeTask struct {
	ID         uuid.UUID  `db:"id" json:"id"`                   // Unique identifier for the node-task relationship
	NodeID     uuid.UUID  `db:"node_id" json:"node_id"`         // ID of the node
	TaskID     uuid.UUID  `db:"task_id" json:"task_id"`         // ID of the task
	TaskOrder  int        `db:"order" json:"task_order"`        // Order of task execution within the node
	Status     string     `db:"status" json:"status"`           // Current status of the task (e.g., pending, completed)
	RetryCount int        `db:"retry_count" json:"retry_count"` // Number of retries for this task
	Data       string     `db:"data" json:"data"`               // Data associated with the task
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`   // Creation timestamp
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`   // Update timestamp
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`   // Deletion timestamp
	Task       Task       `json:"task"`                         // Embedded Task details
}

type Task struct {
	ID         uuid.UUID  `db:"id" json:"id"`       // Unique identifier for the task
	Title      string     `db:"title" json:"title"` // Title of the task
	Type       string     `db:"type" json:"type"`
	HttpMethod string     `db:"http_method" json:"http_method"` // Type of the task
	Action     string     `db:"action" json:"action"`           // Action to be performed by the task
	Params     string     `db:"params" json:"params"`           // Parameters for the action
	MaxRetries int        `db:"max_retries" json:"max_retries"` // Maximum number of retries for this task
	CreatedAt  *time.Time `db:"created_at" json:"created_at"`   // Creation timestamp
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at"`   // Update timestamp
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at"`   // Deletion timestamp
}

type Workflow struct {
	ID             uuid.UUID  `db:"id" json:"id"`                             // Unique identifier for the workflow
	Name           string     `db:"name" json:"name"`                         // Name of the workflow
	StartingNodeID uuid.UUID  `db:"starting_node_id" json:"starting_node_id"` // ID of the related node
	Description    string     `db:"description" json:"description"`           // Description of the workflow
	CreatedAt      *time.Time `db:"created_at" json:"created_at"`             // Creation timestamp
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at"`             // Update timestamp
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at"`             // Deletion timestamp
}
type WorkflowLog struct {
	ID           uuid.UUID      `db:"id" json:"id"`                       // Unique identifier for the log entry, referenced by the caller (can be order, provisioning etc)
	WorkflowID   uuid.UUID      `db:"workflow_id" json:"workflow_id"`     // ID of the workflow this log belongs to
	NodeID       *uuid.UUID     `db:"node_id" json:"node_id"`             // ID of the node being logged
	TaskID       *uuid.UUID     `db:"task_id" json:"task_id"`             // ID of the task being logged
	Status       string         `db:"status" json:"status"`               // Status of the node execution (e.g., "success", "failed", "rollback")
	Message      string         `db:"message" json:"message"`             // Message of the node execution
	ExecutedAt   time.Time      `db:"executed_at" json:"executed_at"`     // Timestamp of when the node was executed
	CompletedAt  sql.NullTime   `db:"completed_at" json:"completed_at"`   // Timestamp of when the node execution was completed
	ErrorMessage sql.NullString `db:"error_message" json:"error_message"` // Details of any error during execution
	ActionType   string         `db:"action_type" json:"action_type"`     // Type of action performed (e.g., "execution", "rollback")
	Metadata     sql.NullString `db:"metadata" json:"metadata"`           // Additional metadata (e.g., execution context)
	CreatedAt    time.Time      `db:"created_at" json:"created_at"`       // Log creation timestamp
	UpdatedAt    time.Time      `db:"updated_at" json:"updated_at"`       // Log update timestamp
	HttpCode     *int           `db:"http_code" json:"http_code"`         // HTTP status code returned by the task
	Response     *string        `db:"response" json:"response"`           // Response from the task
	Error        *string        `db:"error" json:"error,omitempty"`       // Error message, if any
}

type WorkflowExecution struct {
	ID                  uuid.UUID    `db:"id" json:"id"`                                         // Unique identifier for the workflow execution
	WorkflowID          uuid.UUID    `db:"workflow_id" json:"workflow_id"`                       // ID of the workflow this log belongs to
	LastExecutedNodeID  *uuid.UUID   `db:"last_executed_node_id" json:"last_executed_node_id"`   // ID of the node being logged
	LastExecutedTaskID  *uuid.UUID   `db:"last_executed_task_id" json:"last_executed_task_id"`   // ID of the task being logged
	ReferenceNumber     string       `db:"reference_number" json:"reference_number"`             // Reference number of the node execution, could be order id, provisioning id, etc
	Status              string       `db:"status" json:"status"`                                 // Status of the node execution (e.g., "success", "failed", "rollback")
	Message             string       `db:"message" json:"message"`                               // Message of the node execution
	LastNodeExecutedAt  time.Time    `db:"last_node_executed_at" json:"last_node_executed_at"`   // Timestamp of when the node was executed
	LastNodeCompletedAt sql.NullTime `db:"last_node_completed_at" json:"last_node_completed_at"` // Timestamp of when the node execution was completed
	LastTaskExecutedAt  time.Time    `db:"last_task_executed_at" json:"last_task_executed_at"`   // Timestamp of when the task was executed
	LastTaskCompletedAt sql.NullTime `db:"last_task_completed_at" json:"last_task_completed_at"` // Timestamp of when the task execution was completed
	CreatedAt           time.Time    `db:"created_at" json:"created_at"`                         // Log creation timestamp
	UpdatedAt           time.Time    `db:"updated_at" json:"updated_at"`                         // Log update timestamp
}
