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
	"meatfloss/gameuser"
	"meatfloss/message"
	"meatfloss/persistent"
	"meatfloss/usermgr"
	"meatfloss/utils"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/mohae/deepcopy"
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
	go c.HandleHelperGuaji()
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

// HandleHelperGuaji ...
func (c *GameClient) HandleHelperGuaji() {
	for {
		exit := false
		select {
		case <-c.helperChan:
			exit = true
			break
		case <-time.After(time.Second * 10):
			c.onPeriodGuaji()
		}
		if exit {
			break
		}
	}
	glog.Info("HandleHelperGuaji exit")
	c.waitGroup.Done()
}

func (c *GameClient) onPeriod() {
	c.lock.Lock()
	if c.UserID != 0 {
		c.periodCheck()
	}
	c.lock.Unlock()
}

func (c *GameClient) onPeriodGuaji() {
	c.lock.Lock()
	if c.UserID != 0 {
		c.periodCheckGuaji()
	}
	c.lock.Unlock()
}
func (c *GameClient) periodCheckGuaji() {
	c.checkGuajiOutput()
}

func (c *GameClient) periodCheck() {
	c.checkTasks()

}

func (c *GameClient) checkGuajiOutput() {
	fmt.Println("你好吗")

	var size int = len(c.user.GuajiOutputBox.GuajiOutputs)

	if size == 5 {
		for i := 0; i < size-1; i++ {
			//arr[base] = c.user.GuajiOutputBox.GuajiOutputs[i].Items
			c.user.GuajiOutputBox.GuajiOutputs[i] = c.user.GuajiOutputBox.GuajiOutputs[i+1]
		}
		c.user.GuajiOutputBox.GuajiOutputs = append(c.user.GuajiOutputBox.GuajiOutputs[:4], c.user.GuajiOutputBox.GuajiOutputs[5:]...)
	}
	oneEvent := &common.GuajiOutputInfo{}
	oneEvent.UserID = c.user.UserID
	oneEvent.Type = "z"
	oneEvent.Name = "小明"
	oneEvent.Items = "你好"
	oneEvent.Time = time.Now().Format("2006-01-02 15:04:05")
	c.user.GuajiOutputBox.GuajiOutputs = append(c.user.GuajiOutputBox.GuajiOutputs, oneEvent)
	fmt.Println(len(c.user.GuajiOutputBox.GuajiOutputs))
	c.persistOutput()
}
func (c *GameClient) checkTasks() {
	if len(c.user.TaskBox.Tasks) == 0 {
		return
	}

	glog.Info("found task finished ...................")

	now := int(time.Now().Unix())
	for userID, taskInfo := range c.user.TaskBox.Tasks {
		if now <= taskInfo.Timestamp+taskInfo.PreTime {
			continue
		}
		taskInfo.UserID = userID
		var msg = message.EventNotify{}
		oneEvent := &message.EventInfo{}
		msg.Meta.MessageType = "EventNotify"
		msg.Meta.MessageTypeID = message.MsgTypeEventNotify
		msg.Data.UserID = userID
		oneEvent.Type = "normal"
		oneEvent.Title = "任务完成"
		oneEvent.Content = "任务已经完成了"
		oneEvent.Time = time.Now().Format("2006-01-02 15:04:05")
		oneEvent.UserID = userID
		oneEvent.EventID = taskInfo.TaskID
		taskID := gameredis.GetUniqueID()
		oneEvent.GenID = strconv.FormatInt(taskID, 10)
		msg.Data.Events = append(msg.Data.Events, oneEvent)
		{
			cpy := deepcopy.Copy(oneEvent)
			event, _ := cpy.(*message.EventInfo)
			c.user.EventBox.Events[oneEvent.GenID] = event
		}
		c.persistEventBox()
		c.SendMsg(msg)
	}

	c.user.TaskBox.Tasks = make([]*common.TaskInfo, 0)
	c.persistTaskBox()
}

func (c *GameClient) persistTaskBox() {
	newUser := &gameuser.User{}
	newUser.UserID = c.UserID

	cpy := deepcopy.Copy(c.user.TaskBox)
	taskBox, _ := cpy.(*gameuser.TaskBox)
	newUser.TaskBox = taskBox
	persistent.AddUser(c.UserID, newUser)
}

func (c *GameClient) persistEventBox() {
	newUser := &gameuser.User{}
	newUser.UserID = c.UserID
	//newUser.TaskBox = c.user.TaskBox

	cpy := deepcopy.Copy(c.user.EventBox)
	eventBox, _ := cpy.(*gameuser.EventBox)
	newUser.EventBox = eventBox
	persistent.AddUser(c.UserID, newUser)
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
		return c.HandleCreateTaskReq(metaData, rawMsg)
	case message.MsgTypeFinishEventReq:
		return c.HandleFinishEventReq(metaData, rawMsg)
	case message.MsgTypeSaveClientLayoutReq:
		return c.HandleSaveClientLayoutReq(metaData, rawMsg)
	case message.MsgTypeOutputReq:
		return c.HandleOutputReq(metaData, rawMsg)
	}

	return
}

//回复前端的信息
//HandleOutputReq  ...
func (c *GameClient) HandleOutputReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.OutputNotify{}
	reply.Meta.MessageType = "OutputNotify"
	reply.Meta.MessageTypeID = message.MsgTypeOutputNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	if len(c.user.GuajiOutputBox.GuajiOutputs) == 0 {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "GuajiOutputs don't exits"
		c.SendMsg(reply)
		return
	}
	//消息的推送
	//Events: = make([]common.EventInfo, 0)
	reply.Data.GuajiOutputs = make([]common.GuajiOutputInfo, len(c.user.GuajiOutputBox.GuajiOutputs))
	fmt.Println(len(c.user.GuajiOutputBox.GuajiOutputs))
	for _, myOutputs := range c.user.GuajiOutputBox.GuajiOutputs {
		reply.Data.GuajiOutputs = append(reply.Data.GuajiOutputs, *myOutputs)
	}
	fmt.Println(len(c.user.GuajiOutputBox.GuajiOutputs))
	c.SendMsg(reply)

	return
}

func (c *GameClient) persistOutput() {
	newOutput := &gameuser.User{}
	newOutput.UserID = c.UserID

	cpy := deepcopy.Copy(c.user.GuajiOutputBox)
	output, _ := cpy.(*gameuser.GuajiOutputBox)
	newOutput.GuajiOutputBox = output
	persistent.AddUser(c.UserID, newOutput)
}

//HandleSaveClientLayoutReq  ...
func (c *GameClient) HandleSaveClientLayoutReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {
	reply := &message.SaveClientLayoutReply{}
	reply.Meta.MessageType = "SaveClientLayoutReply"
	reply.Meta.MessageTypeID = message.MsgTypeSaveClientLayoutReply
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	fmt.Println("in SaveClientLayoutReply")

	req := &message.SaveClientLayoutReq{}
	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "invalid request"
		c.SendMsg(reply)
		return
	}
	cpy := deepcopy.Copy(req.Data.Layout)
	layout, _ := cpy.(*message.ClientLayout)
	c.user.Layout = layout
	c.persistLayout()
	c.SendMsg(reply)
	return
}

func (c *GameClient) persistLayout() {
	newUser := &gameuser.User{}
	newUser.UserID = c.UserID

	cpy := deepcopy.Copy(c.user.Layout)
	Layout, _ := cpy.(*message.ClientLayout)
	newUser.Layout = Layout
	persistent.AddUser(c.UserID, newUser)
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
	eventInfo, ok := c.user.EventBox.Events[req.Data.EventGenID]
	if !ok {
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
		delete(c.user.EventBox.Events, req.Data.EventGenID)
		c.persistEventBox()
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
	c.persistBagBox()
	// 删除事件
	c.persistEventBox()
	c.SendMsg(notify)
}

func (c *GameClient) persistBagBox() {
	newUser := &gameuser.User{}
	newUser.UserID = c.UserID

	cpy := deepcopy.Copy(c.user.Bag)
	bag, _ := cpy.(*common.Bag)
	newUser.Bag = bag
	persistent.AddUser(c.UserID, newUser)
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
			_, ok := c.user.Bag.Cells[goods.UniqueID]
			if !ok {
				cellNumAtLeast++
			}
		}
	}

	if cellNumAtLeast+len(c.user.Bag.Cells) > 81 {
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
			c.user.Bag.Cells[uniqueID] = cell

			gui := &message.GoodsUpdateInfo{}
			gui.GoodsID = cell.GoodsID
			gui.GoodsNum = goodsCounts[i] // actually, it is 1.
			gui.GoodsNumDelta = goodsCounts[i]
			gui.UniqueID = uniqueID
			deltaMap[uniqueID] = gui
		} else {
			// 如果允许重叠
			cell, ok := c.user.Bag.Cells[goods.UniqueID]
			if !ok {
				cell = &common.BagCell{}
				cell.Count = goodsCounts[i]
				cell.GoodsID = goodsIDs[i]
				cell.UniqueID = goods.UniqueID
				c.user.Bag.Cells[goods.UniqueID] = cell
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

	if len(c.user.TaskBox.Tasks) != 0 {
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

	newTaskInfo := &common.TaskInfo{}
	newTaskInfo.TaskID = task.ID
	newTaskInfo.ID = gameredis.GetUniqueID()
	newTaskInfo.NPCID = npc.ID
	newTaskInfo.Timestamp = int(time.Now().Unix())
	newTaskInfo.PreTime = task.PreTime
	c.user.TaskBox.Tasks = append(c.user.TaskBox.Tasks, newTaskInfo)
	c.persistTaskBox()
	reply.Data.TaskID = strconv.FormatInt(newTaskInfo.ID, 10)
	c.SendMsg(reply)
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
	c.persistLoginTime()
	fmt.Println(c.UserID)
	fmt.Println("小花花花花花花花花")
	return
}

func (c *GameClient) persistLoginTime() {
	newUser := &gameuser.User{}
	newUser.UserID = c.UserID
	cpy := deepcopy.Copy(c.user.LoginTime)
	logintime, _ := cpy.(*gameuser.LoginTime)
	newUser.LoginTime = logintime
	persistent.AddUser(c.UserID, newUser)
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
		c.user = user
	}
	//send message
	err = c.PushRoleInfo()
	if err != nil {
		return
	}
	//数据本地化
	err = c.DownLocalizationInfo()
	if err != nil {
		return
	}
	return
}

// DownLocalizationInfo ...
func (c *GameClient) DownLocalizationInfo() (err error) {

	return
}

// PushRoleInfo ...
func (c *GameClient) PushRoleInfo() (err error) {
	// send base info.

	msg := &message.GameBaseInfoNotify{}
	msg.Meta.MessageType = "GameBaseInfoNotify"
	msg.Meta.MessageTypeID = message.MsgTypeGameBaseInfoNotify

	// copy profile.
	msg.Data.Profile = &message.RoleProfile{}
	msg.Data.Profile.UserID = c.UserID
	msg.Data.Profile.Experience = c.user.Profile.Experience
	msg.Data.Profile.Gender = c.user.Profile.Gender
	msg.Data.Profile.Intelligence = c.user.Profile.Intelligence
	msg.Data.Profile.Intimacy = c.user.Profile.Intimacy
	msg.Data.Profile.Level = c.user.Profile.Level
	msg.Data.Profile.Name = c.user.Profile.Name
	msg.Data.Profile.Spine = c.user.Profile.Spine
	msg.Data.Profile.Stamina = c.user.Profile.Stamina

	// bag.
	msg.Data.Bag.Cells = make([]message.CellInfo, 0)
	for _, mycell := range c.user.Bag.Cells {
		cell := message.CellInfo{}
		cell.Count = mycell.Count
		cell.GoodsID = mycell.GoodsID
		cell.UniqueID = mycell.UniqueID
		msg.Data.Bag.Cells = append(msg.Data.Bag.Cells, cell)
	}

	{
		cpy := deepcopy.Copy(c.user.Layout)
		layout, _ := cpy.(*message.ClientLayout)
		msg.Data.Layout = layout
	}

	// tasks
	msg.Data.Tasks = make([]common.TaskInfo, 0)
	for _, myTask := range c.user.TaskBox.Tasks {
		msg.Data.Tasks = append(msg.Data.Tasks, *myTask)
	}
	// events
	msg.Data.Events = make([]message.EventInfo, 0)
	for _, myEvent := range c.user.EventBox.Events {
		msg.Data.Events = append(msg.Data.Events, *myEvent)
	}

	//layout

	c.SendMsg(msg)
	// var msgEvent = message.EventNotify{}
	// msgEvent.Meta.MessageType = "EventNotify"
	// msgEvent.Meta.MessageTypeID = message.MsgTypeEventNotify
	// msgEvent.Data.Events = make([]*message.EventInfo, 0)
	// for _, myEvent := range c.user.EventBox.Events {
	// 	msgEvent.Data.Events = append(msgEvent.Data.Events, myEvent)

	// }
	// msgEvent.Data.UserID = c.user.UserID
	// c.SendMsg(msgEvent)
	// //task
	// var msgTask = message.TaskNotify{}
	// msgTask.Meta.MessageType = "TaskNotify"
	// msgTask.Meta.MessageTypeID = message.MsgTypeTaskNotify
	// msgTask.Data.Tasks = make([]*common.TaskInfo, 0)

	// for _, myTask := range c.user.TaskBox.Tasks {
	// 	// 必要的时候, deepcopy一份
	// 	msgTask.Data.Tasks = append(msgTask.Data.Tasks, myTask)
	// }
	// msgTask.Data.UserID = c.user.UserID
	// c.SendMsg(msgTask)

	// // send events.
	// {

	// 	var msg = message.EventNotify{}
	// 	msg.Meta.MessageType = "EventNotify"
	// 	msg.Meta.MessageTypeID = message.MsgTypeEventNotify
	// 	msg.Data.Events = make([]*message.EventInfo, 0)
	// 	for _, myEvent := range c.user.EventBox.Events {
	// 		msg.Data.Events = append(msg.Data.Events, myEvent)

	// 	}
	// 	msg.Data.UserID = c.user.UserID
	// 	c.SendMsg(msg)
	// }

	// {

	// 	var msg = message.TaskNotify{}
	// 	msg.Meta.MessageType = "TaskNotify"
	// 	msg.Meta.MessageTypeID = message.MsgTypeTaskNotify
	// 	msg.Data.Tasks = make([]*common.TaskInfo, 0)

	// 	for _, myTask := range c.user.TaskBox.Tasks {
	// 		// 必要的时候, deepcopy一份
	// 		msg.Data.Tasks = append(msg.Data.Tasks, myTask)
	// 	}
	// 	msg.Data.UserID = c.user.UserID
	// 	c.SendMsg(msg)
	// }

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
