package task

// tableName is the name of the environment variable to retrieve the task
// table's name from.
const tableName = "DYNAMODB_TABLE_TASK"

// Task defines the task entity - the primary entity of task domain.
type Task struct {
	ID           string //guid
	Title        string
	Description  string
	Order        int
	Subtasks     []Subtask
	BoardID      string //guid
	ColumnNumber int
}

// Subtask defines the subtask entity which a task may own one/many of.
type Subtask struct {
	Title  string
	IsDone bool
}
