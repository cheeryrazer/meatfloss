package usermgr

import (
	"meatfloss/gameuser"
	"meatfloss/persistent"
	"sync"
	"fmt"
)

var (
	allUsersLock sync.RWMutex
	// 所有已经加载到内存中的用户
	allUsers map[int]*gameuser.User
)

func init() {
	allUsers = make(map[int]*gameuser.User)
}

// GetUser ...
func GetUser(userID int) *gameuser.User {
	allUsersLock.RLock()
	user, ok := allUsers[userID]
	if ok {
		allUsersLock.RUnlock()
		return user
	}
	allUsersLock.RUnlock()

	user = persistent.LoadUser(userID)
	fmt.Println(user.GuajiSettlement)
	return user
}

// NewUser ...
func NewUser(userID int) *gameuser.User {
	user := gameuser.NewUser(userID)
	allUsersLock.Lock()
	allUsers[userID] = user
	allUsersLock.Unlock()
	return user
}
