package gameuser

import (
	"meatfloss/common"
	"meatfloss/message"
	"sync"
)

// User ...
type User struct {
	Lock     sync.RWMutex
	UserID   int
	Profile  *Profile
	Bag      *common.Bag
	TaskBox  *TaskBox
	NewsBox  *NewsBox
	EventBox *EventBox
	Layout   *message.ClientLayout
}

// NewUser ...
func NewUser(userID int) (user *User) {
	user = &User{}
	user.UserID = userID
	user.Profile = NewProfile(userID)
	user.Bag = common.NewBagWithInitialData(userID)
	user.TaskBox = NewTaskBox(userID)
	user.NewsBox = NewNewsBox(userID)
	user.EventBox = NewEventBox(userID)
	user.Layout = message.NewClientLayout()
	return
}
