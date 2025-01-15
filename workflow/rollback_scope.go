package workflow

type RollbackScope string

const (
	RollbackToStart RollbackScope = "start"  // Rollback to the start of the workflow
	RollbackOne     RollbackScope = "one"    // Rollback one ancestor
	RollbackFinish  RollbackScope = "finish" // Stop/finish the workflow on rollback
)
