package gameuser

import (
	"sync"
)

// User ...
type User struct {
	Lock     sync.RWMutex
	Profile  *Profile
	Bag      *Bag
	TaskBox  *TaskBox
	NewsBox  *NewsBox
	EventBox *EventBox
}

// NewUser ...
func NewUser(userID int) (user *User) {
	user = &User{}
	user.Profile = NewProfile(userID)
	user.Bag = NewBag()
	user.TaskBox = NewTaskBox()
	user.NewsBox = NewNewsBox()
	user.EventBox = NewEventBox()
	return
}
