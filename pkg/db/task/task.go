package task

// tableName is the name of the environment variable to retrieve the task
// table's name from.
const tableName = "DYNAMODB_TABLE_TASK"

// Task defines a task in the task table.
type Task struct {
	ID           string //guid
	Title        string
	Description  string
	Order        int
	Subtasks     []Subtask
	BoardID      string //guid
	ColumnNumber int
}

// Subtask defines a subtask in a task's subtasks.
type Subtask struct {
	Title  string
	IsDone bool
}
