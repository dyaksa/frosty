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
