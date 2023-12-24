// Package tasktbl contains code to interact with the task table in DynamoDB.
package tasktbl

// tableName is the name of the environment variable to retrieve the task
// table's name from.
const tableName = "TASK_TABLE_NAME"

// Task defines the task entity - the primary entity of task domain.
type Task struct {
	TeamID       string //guid
	BoardID      string //guid
	ColumnNumber int
	ID           string //guid
	Title        string
	Description  string
	Order        int
	Subtasks     []Subtask
}

// NewTask creates and returns a new Task.
func NewTask(
	teamID string,
	boardID string,
	colNo int,
	id string,
	title string,
	descr string,
	order int,
	subtasks []Subtask,
) Task {
	return Task{
		TeamID:       teamID,
		BoardID:      boardID,
		ColumnNumber: colNo,
		ID:           id,
		Title:        title,
		Description:  descr,
		Order:        order,
		Subtasks:     subtasks,
	}
}

// Subtask defines the subtask entity which a task may own one/many of.
type Subtask struct {
	Title  string
	IsDone bool
}

// NewSubtask creates and returns a new Subtask.
func NewSubtask(title string, isDone bool) Subtask {
	return Subtask{
		Title:  title,
		IsDone: isDone,
	}
}
