package workflow

import (
	"time"

	"github.com/google/uuid"
)

type Node struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Title       string     `db:"title" json:"title"`
	Type        string     `db:"type" json:"type"`
	Description string     `db:"description,omitempty" json:"description,omitempty"`
	CreatedAt   *time.Time `db:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt   *time.Time `db:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type NodeClosure struct {
	Ancestor   uuid.UUID `db:"ancestor" json:"ancestor"`
	Descendant uuid.UUID `db:"descendant" json:"descendant"`
	Depth      int       `db:"depth" json:"depth"`
}

type Task struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Type      string     `db:"type" json:"type"`
	Action    string     `db:"action" json:"action"`
	Params    string     `db:"params" json:"params"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type NodeTask struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	NodeID    uuid.UUID  `db:"node_id" json:"node_id"`
	TaskID    uuid.UUID  `db:"task_id" json:"task_id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type Workflow struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description string     `db:"description" json:"description"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

type WorkflowNode struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	WorkflowID     uuid.UUID  `db:"workflow_id" json:"workflow_id"`
	StartingNodeID uuid.UUID  `db:"node_id" json:"strating_node_id"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at"`
}
