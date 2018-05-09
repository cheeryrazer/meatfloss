package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"meatfloss/db"
	"meatfloss/gameconf"
	"meatfloss/gameuser"
	"meatfloss/message"
	"meatfloss/persistent"
	"meatfloss/usermgr"
	"meatfloss/utils"
	"net"
	"sync"
	"time"

	"github.com/mohae/deepcopy"

	"github.com/golang/glog"

	"github.com/gorilla/websocket"
)

// GameClient  ...
type GameClient struct {
	lock        sync.RWMutex
	conn        *websocket.Conn
	UniqueID    uint64
	UserID      int
	kickOffFlag int32
	logined     bool
	replyChan   chan interface{}
	helperChan  chan bool
	waitGroup   sync.WaitGroup

	maxNewsPushID uint64
	KickOffChan   chan bool
	user          *gameuser.User
}

// HandleConnection ...
func (c *GameClient) HandleConnection(conn *websocket.Conn) {
	c.UniqueID = utils.GetUniqueID()
	c.conn = conn
	c.replyChan = make(chan interface{}, 128)
	c.KickOffChan = make(chan bool, 1)
	c.helperChan = make(chan bool, 1)

	//c.conn.SetReadDeadline((time.Now().Add(5 * time.Second)))
	c.waitGroup.Add(3)
	go c.HandleRead()
	go c.HandleWrite()
	go c.HandleHelper()

	if c.UserID != 0 {
		Mgr.onLogout(c)
	}
	c.waitGroup.Wait()
	c.conn.Close()

	c.KickOffChan <- true

	glog.Info("session is gonna close")
}

// HandleHelper ...
func (c *GameClient) HandleHelper() {
	for {
		exit := false
		select {
		case <-c.helperChan:
			exit = true
			break
		case <-time.After(time.Second * 5):
			c.onPeriod()
		}
		if exit {
			break
		}
	}
	glog.Info("HandleHelper exit")
	c.waitGroup.Done()
}

func (c *GameClient) onPeriod() {
	c.lock.Lock()
	if c.UserID != 0 {
		c.periodCheck()
	}
	c.lock.Unlock()
}

func (c *GameClient) periodCheck() {
	glog.Info("check..................................")
}

// HandleWrite ...
func (c *GameClient) HandleWrite() {
	for {
		msg := <-c.replyChan
		if msg == nil {
			break
		}

		result, err := json.Marshal(msg)
		_ = err
		fmt.Println(string(result))
		err = c.conn.WriteMessage(1, result)
		if err != nil {
			log.Println("write:", err)
		}
	}

	c.waitGroup.Done()
}

// HandleRead ...
func (c *GameClient) HandleRead() {
	conn := c.conn
	for {
		mt, message, err := conn.ReadMessage()
		if c.kickOffFlag == 1 {
			c.handleKickOff()
			break
		}

		if err != nil {
			log.Println("read:", err)
			break
		}
		_ = mt
		c.lock.Lock()
		err = c.HandleMessage(message)
		c.lock.Unlock()

		if err != nil {
			glog.Errorf("HandleMessage failed, error: %s", err)
			break
		}
	}
	c.replyChan <- nil
	c.helperChan <- true
	c.waitGroup.Done()
}

func (c *GameClient) handleKickOff() {
	msg := &message.KickOffNotify{}
	msg.Meta.MessageType = "KickOffNotify"
	c.SendMsg(msg)
}

// HandleMessage ...
func (c *GameClient) HandleMessage(rawMsg []byte) (err error) {
	glog.Info("HandleMessage called")
	meta := &message.ReqMeta{}
	err = json.Unmarshal(rawMsg, meta)
	if err != nil {
		return
	}

	metaData := meta.Meta
	if !c.logined && metaData.MessageType != "LoginReq" {
		glog.Info("non-login message received before login.")
		return
	}

	fmt.Printf("+%v", metaData)
	switch metaData.MessageTypeID {
	case message.MsgTypeLoginReq:
		return c.HandleLoginReq(metaData, rawMsg)
	case message.MsgTypeMarkNewsAsReadReq:
	//	return c.HandleMarkNewsAsReadReq(metaData, rawMsg)
	case message.MsgTypeCreateTaskReq:
	//	return c.HandleCreateTaskReq(metaData, rawMsg)
	case message.MsgTypeFinishEventReq:
		//	return c.HandleFinishEventReq(metaData, rawMsg)
	}

	return
}

// HandleLoginReq ...
func (c *GameClient) HandleLoginReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {
	if c.logined {
		// multiple login disallowed.
		return
	}
	req := &message.LoginReq{}
	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		return
	}

	{
		if addr, ok := c.conn.RemoteAddr().(*net.TCPAddr); ok {
			req.Data.Account = addr.IP.String()
		}
	}

	reply := &message.LoginReply{}
	reply.Meta.MessageType = "LoginReply"
	reply.Meta.MessageTypeID = message.MsgTypeLoginReply
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	fmt.Printf("%v\n", req.Data.Account)
	if len(req.Data.Account) < 2 || len(req.Data.Account) > 30 {
		reply.Meta.ErrorMessage = "Authorize failed."
		reply.Meta.Error = true
		c.SendMsg(reply)
		return errors.New("Authorize failed")
	}
	userID, err := db.GetUserID(req.Data.Account)
	if err != nil {
		userID, err = db.CreateAccount(req.Data.Account)
		if err != nil {
			reply.Meta.ErrorMessage = "Authorize failed."
			reply.Meta.Error = true
			c.SendMsg(reply)
			return errors.New("Authorize failed")
		}
		err = c.InitUser(userID)
		if err != nil {
			glog.Errorf("InitRole failed, userID: %d", userID)
			return errors.New("Authorize failed")
		}
	}
	c.logined = true
	c.UserID = userID
	_ = userID

	Mgr.onNewLogin(c)
	c.SendMsg(reply)
	err = c.AfterLogin()
	return
}

// HandleTimeout ... true ignore error, otherwise not.
func (c *GameClient) HandleTimeout() bool {
	//	c.conn.SetReadDeadline((time.Now().Add(5 * time.Second)))
	//	glog.Error("HandleTimeout")
	//return true
	return false
}

// AfterLogin ...
func (c *GameClient) AfterLogin() (err error) {
	if c.user == nil {
		user := usermgr.GetUser(c.UserID)
		if user == nil {
			glog.Errorf("usermgr.GetUser failed, userID: %d", c.UserID)
			return errors.New("Load user failed")
		}
	}
	return
}

// InitUser ...
func (c *GameClient) InitUser(userID int) (err error) {
	user := gameuser.NewUser(userID)
	// TODO, init user.
	c.user = user

	cpy := deepcopy.Copy(user)
	newUser, _ := cpy.(*gameuser.User)
	_ = newUser
	persistent.AddUser(userID, newUser)
	return
}

// SendMsg ...
func (c *GameClient) SendMsg(msg interface{}) {
	c.replyChan <- msg
}

// TrySendMsg ...
func (c *GameClient) TrySendMsg(msg interface{}) {
	select {
	case c.replyChan <- msg:
	default:
	}
}

func (c *GameClient) kickOff() {
	// TODO, add mutex.
	if c.kickOffFlag == 1 {
		return
	}
	c.kickOffFlag = 1
	c.conn.SetReadDeadline(time.Now().Add(0))
}

// GetTasksByNPC ...
func (c *GameClient) GetTasksByNPC(npc *gameconf.NPC) (event *gameconf.Task) {
	events := gameconf.GetTasksByNPC(npc.ID)
	if len(events) == 0 {
		return
	}

	// TODO， 根据概率， 等级等选择一个event.
	which := rand.Intn(len(events))
	event = events[which]
	return
}
