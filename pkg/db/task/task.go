package task

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
