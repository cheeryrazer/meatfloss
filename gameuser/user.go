package gameuser

import (
	"sync"
)

// User ...
type User struct {
	Lock    sync.RWMutex
	Profile *Profile
	Bag     *Bag
}

// NewUser ...
func NewUser(userID int) (user *User) {
	user = &User{}
	user.Profile = NewProfile(userID)
	user.Bag = NewBag()
	return
}
