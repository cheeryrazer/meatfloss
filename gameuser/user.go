package gameuser

import (
	"sync"
)

// User ...
type User struct {
	Lock     sync.RWMutex
	UserID   int
	Profile  *Profile
	Bag      *Bag
	TaskBox  *TaskBox
	NewsBox  *NewsBox
	EventBox *EventBox
}

// NewUser ...
func NewUser(userID int) (user *User) {
	user = &User{}
	user.UserID = userID
	user.Profile = NewProfile(userID)
	user.Bag = NewBag(userID)
	user.TaskBox = NewTaskBox(userID)
	user.NewsBox = NewNewsBox(userID)
	user.EventBox = NewEventBox(userID)
	return
}
