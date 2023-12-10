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

// NewTask creates and returns a new Task.
func NewTask(
	id string,
	title string,
	descr string,
	order int,
	subtasks []Subtask,
	boardID string,
	colNo int,
) Task {
	return Task{
		ID:           id,
		Title:        title,
		Description:  descr,
		Order:        order,
		Subtasks:     subtasks,
		BoardID:      boardID,
		ColumnNumber: colNo,
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
