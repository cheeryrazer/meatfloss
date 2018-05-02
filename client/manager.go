package client

import (
	"sync"
)

var (
	// Mgr is the unique instance of Manager.
	Mgr *Manager
)

// Manager ...
type Manager struct {
	rwMutex        sync.RWMutex
	allOnlineUsers map[int]*GameClient
}

func init() {
	Mgr = newManager()
}

// NewManager xxx
func newManager() *Manager {
	mgr := &Manager{}
	mgr.allOnlineUsers = make(map[int]*GameClient)
	return mgr
}

func (m *Manager) onNewLogin(cli *GameClient) {
	m.rwMutex.Lock()
	oldCli, ok := m.allOnlineUsers[cli.UserID]
	if !ok {
		// 没找到， 说明我就是最新用户
		m.allOnlineUsers[cli.UserID] = cli
		m.rwMutex.Unlock()
		return
	}
	m.allOnlineUsers[cli.UserID] = cli
	m.rwMutex.Unlock()
	oldCli.kickOff()
	kickOff := <-oldCli.KickOffChan
	_ = kickOff
}

func (m *Manager) onLogout(cli *GameClient) {
	m.rwMutex.Lock()
	oldCli, ok := m.allOnlineUsers[cli.UserID]
	if !ok || cli.UniqueID != oldCli.UniqueID {
		m.rwMutex.Unlock()
		return
	}
	delete(m.allOnlineUsers, cli.UserID)
	m.rwMutex.Unlock()
}

// Broadcast ...
func (m *Manager) Broadcast(msg interface{}) {
	m.rwMutex.RLock()
	users := make(map[int]*GameClient)
	for k, v := range m.allOnlineUsers {
		users[k] = v
	}
	m.rwMutex.RUnlock()

	for userID, client := range users {
		_ = userID
		_ = client
		client.TrySendMsg(msg)
	}
}

// SendToClient ...
func (m *Manager) SendToClient(userID int, msg interface{}) {
	m.rwMutex.Lock()
	client, ok := m.allOnlineUsers[userID]
	if !ok {
		m.rwMutex.Unlock()
		return
	}
	m.rwMutex.Unlock()
	client.TrySendMsg(msg)
}
