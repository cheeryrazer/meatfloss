package persistent

import (
	"encoding/json"
	"fmt"
	"meatfloss/gameuser"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	lock         sync.RWMutex
	changedUsers map[int]*gameuser.User
	levelDB      *leveldb.DB
)

func init() {
	changedUsers = make(map[int]*gameuser.User)
}

// Start ...
func Start() {
	leveldbDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if !strings.HasSuffix(leveldbDir, string(os.PathSeparator)) {
		leveldbDir += string(os.PathSeparator)
	}
	leveldbDir += "leveldb"
	db, err := leveldb.OpenFile(leveldbDir, nil)
	if err != nil {
		glog.Error("leveldb.OpenFile failed: ", err)
	}
	levelDB = db
	go Worker()
}

func persistUsers() {
	lock.Lock()
	users := changedUsers
	changedUsers = make(map[int]*gameuser.User)
	_ = users
	lock.Unlock()
	batch := new(leveldb.Batch)
	for userID, usr := range changedUsers {
		_ = userID
		_ = usr
		usr.Lock.RLock()
		data, err := json.Marshal(usr)
		usr.Lock.RUnlock()
		if err != nil {
			glog.Info("json.Marshal failed in persistUsers")
			continue
		}
		key := fmt.Sprintf("user:%d", userID)
		batch.Put([]byte(key), data)
	}
	if batch.Len() > 0 {
		err := levelDB.Write(batch, nil)
		if err != nil {
			glog.Info("levelDB.Write failed.")
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
	changedUsers[userID] = user
}

// LoadUser ...
func LoadUser(userID int) (user *gameuser.User) {
	key := fmt.Sprintf("user:%d", userID)
	value, err := levelDB.Get([]byte(key), nil)
	if err != nil {
		glog.Info("load user from levelDB failed")
		return nil
	}
	user = &gameuser.User{}
	err = json.Unmarshal(value, user)
	if err != nil {
		glog.Info("json.Unmarshal failed in LoadUser")
		return nil
	}
	return
}
