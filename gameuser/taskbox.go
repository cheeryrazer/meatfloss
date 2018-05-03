package gameuser

// TaskBox ...
type TaskBox struct {
	UserID int
}

// NewTaskBox ...
func NewTaskBox(userID int) *TaskBox {
	box := &TaskBox{}
	box.UserID = userID
	return box
}
