package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"meatfloss/common"
	"meatfloss/db"
	"meatfloss/gameconf"
	"meatfloss/gameredis"
	"meatfloss/message"
	"meatfloss/utils"
	"net"
	"os/user"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"

	"github.com/gorilla/websocket"
)

// GameClient  ...
type GameClient struct {
	rwmutex       sync.RWMutex
	conn          *websocket.Conn
	UniqueID      uint64
	UserID        int
	kickOffFlag   int32
	logined       bool
	replyChan     chan interface{}
	helperChan    chan interface{}
	waitGroup     sync.WaitGroup
	maxNewsPushID uint64
	KickOffChan   chan bool
	profile       message.RoleProfile
	bag           *common.Bag
	user          *user.User
}

// HandleConnection ...
func (c *GameClient) HandleConnection(conn *websocket.Conn) {
	c.UniqueID = utils.GetUniqueID()
	c.conn = conn
	c.replyChan = make(chan interface{}, 128)
	c.KickOffChan = make(chan bool, 1)
	go c.HandleRead()
	go c.HandleWrite()
	c.waitGroup.Add(2)

	if c.UserID != 0 {
		Mgr.onLogout(c)
	}
	c.waitGroup.Wait()
	c.conn.Close()

	c.KickOffChan <- true
	glog.Info("session is gonna close")
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
		err = c.HandleMessage(message)

		if err != nil {
			glog.Errorf("HandleMessage failed, error: %s", err)
			break
		}
	}
	c.replyChan <- nil
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
		return c.HandleMarkNewsAsReadReq(metaData, rawMsg)
	case message.MsgTypeCreateTaskReq:
		return c.HandleCreateTaskReq(metaData, rawMsg)
	case message.MsgTypeFinishEventReq:
		return c.HandleFinishEventReq(metaData, rawMsg)
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
		err = c.InitRole(userID)
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

// AfterLogin ...
func (c *GameClient) AfterLogin() (err error) {
	//guest.AddUserID(c.UserID)
	err = c.PushRoleInfo()
	if err != nil {
		return
	}
	c.PushNews()
	// TODO  close sesssion if error occurred
	return
}

// PushRoleInfo ...
func (c *GameClient) PushRoleInfo() (err error) {
	notify, events, bag, npcGuests, err := gameredis.GetRoleInfo(c.UserID)
	_ = npcGuests
	if err != nil {
		glog.Errorf("gameredis.GetRoleInfo failed, userID: %d", c.UserID)
		return errors.New("redis failure")
	}
	{
		if bag != nil {
			c.bag = bag
		}
	}

	{
		notify.Meta.MessageType = "GameBaseInfoNotify"
		notify.Meta.MessageTypeID = message.MsgTypeGameBaseInfoNotify
		c.SendMsg(notify)
	}

	{
		msg := message.EventNotify{}
		msg.Meta.MessageType = "EventNotify"
		msg.Meta.MessageTypeID = message.MsgTypeEventNotify
		msg.Data.Events = events
		c.SendMsg(msg)
	}

	if len(npcGuests) == 0 {
		newNPCGuests := make([]string, 0)
		newNPCGuests = append(newNPCGuests, "zk00009")
		newNPCGuests = append(newNPCGuests, "zk00010")
		newNPCGuests = append(newNPCGuests, "dummy")
		gameredis.SetNPCGuestList(c.UserID, newNPCGuests)
		npcGuests = newNPCGuests
	}

	{
		msg := message.NPCGuestNotify{}
		msg.Meta.MessageType = "NPCGuestNotify"
		msg.Meta.MessageTypeID = message.MsgTypeNPCGuestNotify
		msg.Data.NPCList = make([]string, 0)
		for _, guest := range npcGuests {
			if guest != "dummy" {
				msg.Data.NPCList = append(msg.Data.NPCList, guest)
			}
		}
		c.SendMsg(msg)
	}
	return
}

// InitRole ...
func (c *GameClient) InitRole(userID int) (err error) {
	return gameredis.InitRole(userID)
}

// PushNews ...
func (c *GameClient) PushNews() {
	articles := gameredis.GetAllNews(c.UserID)
	notify := &message.PushNewsNotify{}
	notify.Meta.MessageType = "PushNewsNotify"
	notify.Meta.MessageTypeID = message.MsgTypePushNewsNotify
	notify.Data.Articles = articles
	c.SendMsg(notify)
}

// HandleMarkNewsAsReadReq ...
func (c *GameClient) HandleMarkNewsAsReadReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {
	req := &message.MarkNewsAsReadReq{}
	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		return
	}

	gameredis.MarkNewsAdRead(c.UserID, req.Data.PushID, req.Data.ArticleID)

	reply := &message.MarkNewsAsReadReply{}
	reply.Meta.MessageType = "MarkNewsAsReadReply"
	reply.Meta.MessageTypeID = message.MsgTypeMarkNewsAsReadReply
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	reply.Data.ArticleID = req.Data.ArticleID
	reply.Data.PushID = req.Data.PushID
	c.SendMsg(reply)
	return
}

// HandleFinishEventReq ...
func (c *GameClient) HandleFinishEventReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {
	req := &message.FinishEventReq{}
	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		return
	}

	reply := &message.FinishEventReply{}
	reply.Meta.MessageType = "FinishEventReply"
	reply.Meta.MessageTypeID = message.MsgTypeFinishEventReply
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	if req.Data.Choice < 1 || req.Data.Choice > 3 {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "bad choice"
		c.SendMsg(reply)
		return errors.New("Bad choice")
	}

	eventInfo := gameredis.GetEvent(c.UserID, req.Data.EventGenID)
	if eventInfo == nil {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "event not found"
		c.SendMsg(reply)
		return errors.New("event not found")
	}
	reply.Data.EventGenID = req.Data.EventGenID
	c.SendMsg(reply)

	// 然后下就是推送奖励
	if eventInfo.Type == "normal" {
		c.OnFinishNormalEvent(eventInfo, req.Data.Choice)
	} else if eventInfo.Type == "select" {
		c.OnFinishRandomEvent(eventInfo)
	}
	return
}

// OnFinishNormalEvent ...
// TODO， 需要改成如果
func (c *GameClient) OnFinishNormalEvent(eventInfo *message.EventInfo, choice int) {
	taskEvent, ok := gameconf.AllTasks[eventInfo.EventID]
	if !ok {
		glog.Warning("OnFinishNormalEvent, taskEvent not found")
		return
	}
	_ = taskEvent
	// ok 得到一个普通事件

	reward := taskEvent.Rewards[choice-1]
	_ = reward

	// // 无任何奖励， 则直接发送
	if len(reward.List) == 0 {
		return
	}

	var goodsIDs []string
	var goodsCounts []int
	for _, sw := range reward.List {
		goodsIDs = append(goodsIDs, sw.GoodsID)
		goodsCounts = append(goodsCounts, sw.GoodsNum)

	}

	updateInfos, err := c.PutToBagBatch(goodsIDs, goodsCounts)
	if err != nil {
		glog.Warning("c.PutToBagBatch failed, error: %s", err)
	}
	_ = updateInfos
	notify := message.UpdateGoodsNotify{}
	notify.Meta.MessageType = "UpdateGoodsNotify"
	notify.Meta.MessageTypeID = message.MsgTypeUpdateGoodsNotify

	notify.Data.List = updateInfos
	err = gameredis.SaveBagInfo(c.UserID, c.bag)
	if err != nil {
		glog.Errorf("SaveBagInfo,  error : %s", err)
		return
	}
	// 调试， 暂时不删
	//	gameredis.DelEvent(c.UserID, eventInfo.GenID)
	c.SendMsg(notify)
}

// PutToBag ...
func (c *GameClient) PutToBag(goodsID string, goodsCount int) (err error) {
	return
}

// PutToBagBatch ...
func (c *GameClient) PutToBagBatch(goodsIDs []string, goodsCounts []int) (infos []message.GoodsUpdateInfo, err error) {
	var goodsList []*gameconf.SuperGoods
	var uniqueIDs []int64
	var cellNumAtLeast int
	var disallowPileupNum int
	for _, goodsID := range goodsIDs {
		goods, ok := gameconf.AllSuperGoods[goodsID]
		if !ok {
			glog.Warning("no such goods, goods id : %s", goods.ID)
			return nil, errors.New("no such goods")
		}
		goodsList = append(goodsList, goods)
		if goods.AllowPileup != 1 {
			cellNumAtLeast++
			disallowPileupNum++
		} else {
			// 如果允许重叠
			_, ok := c.bag.Cells[goods.UniqueID]
			if !ok {
				cellNumAtLeast++
			}
		}
	}

	if cellNumAtLeast+len(c.bag.Cells) > 81 {
		return nil, errors.New("insufficient cells")
	}
	_ = uniqueIDs

	deltaMap := make(map[int64]*message.GoodsUpdateInfo)

	for i, goods := range goodsList {
		if goods.AllowPileup != 1 {
			uniqueID := gameredis.GetGoodsUniqueID()
			cell := &common.BagCell{}
			cell.Count = goodsCounts[i]
			cell.GoodsID = goodsIDs[i]
			cell.UniqueID = goods.UniqueID
			c.bag.Cells[uniqueID] = cell

			gui := &message.GoodsUpdateInfo{}
			gui.GoodsID = cell.GoodsID
			gui.GoodsNum = goodsCounts[i] // actually, it is 1.
			gui.GoodsNumDelta = goodsCounts[i]
			gui.UniqueID = uniqueID
			deltaMap[uniqueID] = gui
		} else {
			// 如果允许重叠
			cell, ok := c.bag.Cells[goods.UniqueID]
			if !ok {
				cell = &common.BagCell{}
				cell.Count = goodsCounts[i]
				cell.GoodsID = goodsIDs[i]
				cell.UniqueID = goods.UniqueID
				c.bag.Cells[goods.UniqueID] = cell
			} else {
				cell.Count += goodsCounts[i]
			}

			gui, ok := deltaMap[goods.UniqueID]
			if !ok {
				// first time.
				gui = &message.GoodsUpdateInfo{}
				gui.GoodsID = cell.GoodsID
				gui.GoodsNum = cell.Count
				gui.GoodsNumDelta = goodsCounts[i]
				deltaMap[goods.UniqueID] = gui

			} else {
				gui.GoodsNum = cell.Count
				gui.GoodsNumDelta += goodsCounts[i]
			}
			gui.UniqueID = goods.UniqueID
			_ = gui
		}
	}

	for _, v := range deltaMap {
		infos = append(infos, *v)
	}
	return
}

// OnFinishRandomEvent ...
func (c *GameClient) OnFinishRandomEvent(eventInfo *message.EventInfo) {

}

// HandleCreateTaskReq ...
func (c *GameClient) HandleCreateTaskReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {
	req := &message.CreateTaskReq{}
	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		return
	}

	reply := &message.CreateTaskReply{}
	reply.Meta.MessageType = "CreateTaskReply"
	reply.Meta.MessageTypeID = message.MsgTypeCreateTaskReply
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	// 写死派遣任务
	// 写死一分钟

	taskInfo, e := gameredis.GetRunningTask(c.UserID)
	if e != nil {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "BackendFailure"
		c.SendMsg(reply)
		return
	}
	if taskInfo != "" {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "Task already exits"
		c.SendMsg(reply)
		return
	}

	npc := gameconf.GetNPC(req.Data.NPCID)
	if npc == nil {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "npc not found"
		c.SendMsg(reply)
		return
	}

	task := c.GetTasksByNPC(npc)
	if task == nil {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "could not create task"
		c.SendMsg(reply)
		return

	}
	_ = taskInfo

	newTaskInfo := &gameredis.RunningTaskInfo{}
	// TODO, check failure.
	newTaskInfo.TaskID = task.ID
	newTaskInfo.ID = gameredis.GetUniqueID()
	newTaskInfo.NPCID = npc.ID
	newTaskInfo.Timestamp = int(time.Now().Unix())
	newTaskInfo.PreTime = task.PreTime
	gameredis.SetRunningTask(c.UserID, newTaskInfo)
	reply.Data.TaskID = strconv.FormatInt(newTaskInfo.ID, 10)
	c.SendMsg(reply)
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
