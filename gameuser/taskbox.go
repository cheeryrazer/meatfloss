package gameuser

// TaskBox ...
type TaskBox struct {
	UserID int
	Tasks  []*TaskInfo
}

// NewTaskBox ...
func NewTaskBox(userID int) *TaskBox {
	box := &TaskBox{}
	box.Tasks = make([]*TaskInfo, 0)
	box.UserID = userID
	return box
}

// TaskInfo ...
type TaskInfo struct {
	TaskID    string
	Timestamp int
	PreTime   int
	UserID    int
	ID        int64
	NPCID     string
}
