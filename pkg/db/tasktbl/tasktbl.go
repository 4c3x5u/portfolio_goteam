// Package tasktbl contains code to interact with the task table in DynamoDB.
package tasktbl

// tableName is the name of the environment variable to retrieve the task
// table's name from.
const tableName = "TASK_TABLE_NAME"

// Task defines the task entity - the primary entity of task domain.
type Task struct {
	TeamID      string    `json:"teamID"`                          // guid
	BoardID     string    `json:"boardID" dynamodbav:",omitempty"` // guid
	ColNo       int       `json:"colNo" dynamodbav:",omitempty"`
	ID          string    `json:"id"` // guid
	Title       string    `json:"title"`
	Description string    `json:"description" dynamodbav:",omitempty"`
	Order       int       `json:"order"`
	Subtasks    []Subtask `json:"subtasks" dynamodbav:",omitempty"`
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
		TeamID:      teamID,
		BoardID:     boardID,
		ColNo:       colNo,
		ID:          id,
		Title:       title,
		Description: descr,
		Order:       order,
		Subtasks:    subtasks,
	}
}

// Subtask defines the subtask entity which a task may own one/many of.
type Subtask struct {
	Title  string `json:"title"`
	IsDone bool   `json:"done"`
}

// NewSubtask creates and returns a new Subtask.
func NewSubtask(title string, isDone bool) Subtask {
	return Subtask{
		Title:  title,
		IsDone: isDone,
	}
}
