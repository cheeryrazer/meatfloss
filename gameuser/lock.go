package gameuser

import (
	"sync"
)

var (
	allUserLocks map[int]*sync.RWMutex
	mapLock      sync.Mutex
)

func init() {
	allUserLocks = make(map[int]*sync.RWMutex)
}

// GetLockByUserID ...
// 用户操作的锁， GetLockByUserID
func GetLockByUserID(userID int) *sync.RWMutex {
	mapLock.Lock()
	defer mapLock.Unlock()
	lck, ok := allUserLocks[userID]
	if ok {
		return lck
	}

	lck = &sync.RWMutex{}
	allUserLocks[userID] = lck
	return lck
}
