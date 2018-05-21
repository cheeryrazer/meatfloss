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
	if !ok {
		changedUsers[userID] = user
		return
	}

	if user.Profile != nil {
		oldUser.Profile = user.Profile
	}

	if user.Bag != nil {
		oldUser.Bag = user.Bag
	}

	if user.TaskBox != nil {
		oldUser.TaskBox = user.TaskBox
	}

	if user.NewsBox != nil {
		oldUser.NewsBox = user.NewsBox
	}

	if user.EventBox != nil {
		oldUser.EventBox = user.EventBox
	}

	if user.Layout != nil {
		oldUser.Layout = user.Layout
	}

	if user.LoginTime != nil {
		oldUser.LoginTime = user.LoginTime
	}

	if user.GuajiOutputBox != nil {
		oldUser.GuajiOutputBox = user.GuajiOutputBox
	}

	if user.GuajiSettlement != nil {
		oldUser.GuajiSettlement = user.GuajiSettlement
	}

	if user.GuajiProfile != nil {
		oldUser.GuajiProfile = user.GuajiProfile
	}

}

// // GetUser ...
// func GetUser(userID int, name string) {
// 	lock.Lock()
// 	defer lock.Unlock()
// 	user := &gameuser.User{}
// 	oldUser, ok := changedUsers[userID]
// 	if !ok {
// 		changedUsers[userID] = user
// 		return
// 	}
// 	if name == "OutputBox" {
// 		return oldUser.OutputBox.Outputs
// 	}

// }

// LoadUser ...
func LoadUser(userID int) (user *gameuser.User) {
	return gameredis.LoadUser(userID)
}
