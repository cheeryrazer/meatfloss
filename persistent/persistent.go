package persistent

import (
	"meatfloss/gameredis"
	"meatfloss/gameuser"
	"sync"
	"time"
)

var (
	lock         sync.RWMutex
	changedUsers map[int]*gameuser.User
)

func init() {
	changedUsers = make(map[int]*gameuser.User)
}

// Start ...
func Start() {
	go Worker()
}

func persistUsers() {
	lock.Lock()
	users := changedUsers
	changedUsers = make(map[int]*gameuser.User)
	_ = users
	lock.Unlock()
	for userID, user := range users {
		for {
			err := gameredis.PersistUser(userID, user)
			if err == nil {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// Worker ...
func Worker() {
	for {
		time.Sleep(1 * time.Second)
		persistUsers()
	}
}

// AddUser ...;
func AddUser(userID int, user *gameuser.User) {
	lock.Lock()
	defer lock.Unlock()
	oldUser, ok := changedUsers[userID]
	if ok {
		return
	}
	_ = oldUser
	changedUsers[userID] = user
}

// LoadUser ...
func LoadUser(userID int) (user *gameuser.User) {
	return gameredis.LoadUser(userID)
}
