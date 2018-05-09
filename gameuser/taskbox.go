package gameuser

import "meatfloss/common"

// TaskBox ...
type TaskBox struct {
	UserID int
	Tasks  []*common.TaskInfo
}

// NewTaskBox ...
func NewTaskBox(userID int) *TaskBox {
	box := &TaskBox{}
	box.Tasks = make([]*common.TaskInfo, 0)
	box.UserID = userID
	return box
}
