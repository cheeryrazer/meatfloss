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
	"meatfloss/logic"
	"meatfloss/message"
	"meatfloss/persistent"
	"meatfloss/usermgr"
	"meatfloss/utils"
	"net"
	"runtime"
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

	//	go c.HandleHelperGuajiTemperature()
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
		case <-time.After(time.Second * 1):
			c.onPeriod()

		}
		if exit {
			break
		}
	}
	glog.Info("HandleHelper exit")
	c.waitGroup.Done()
}

// // HandleHelperGuajiTemperature ...
// func (c *GameClient) HandleHelperGuajiTemperature() {
// 	for {
// 		exit := false
// 		select {
// 		case <-c.helperChan:
// 			exit = true
// 			break
// 		case <-time.After(time.Second * 1):
// 			c.coolTemperature()
// 		}
// 		if exit {
// 			break
// 		}
// 	}
// 	glog.Info("HandleHelperGuajiTemperature exit")
// 	c.waitGroup.Done()
// }
func (c *GameClient) coolTemperature() {
	// fmt.Println(c.user.GuajiProfile.CurrentTemperature)
	//取出当前的等级
	userProfile := c.user.Profile
	//根据等级取出当前的机器的结算数据
	machine := gameconf.AllGuajis
	// 默认等级0+1
	machineInfo := machine[userProfile.Level+1]
	if c.user.GuajiProfile.CDPick > 1 {
		c.user.GuajiProfile.CDPick--
	}
	if c.user.GuajiProfile.CDPick == 1 {
		c.user.GuajiProfile.CDPick--
		c.SendMsg(&message.ClickStatusReq{Status: 1})
	}
	// 取出cd
	cd := c.user.GuajiProfile.CDTemperature
	if cd == 1 {
		c.user.GuajiProfile.CurrentTemperature = float64(machineInfo.InitialTemperature)

	}
	currentTemperature := c.user.GuajiProfile.CurrentTemperature

	if cd >= 1 {
		c.user.GuajiProfile.CDTemperature = cd - 1
		return
	}

	cdPerDegree := machineInfo.CDPerDegree
	c.user.GuajiProfile.CurrentTemperature = currentTemperature - float64(float64(1)/float64(cdPerDegree))
	if c.user.GuajiProfile.CurrentTemperature < float64(machineInfo.InitialTemperature) {
		c.user.GuajiProfile.CurrentTemperature = float64(machineInfo.InitialTemperature)
	}
	// fmt.Println(currentTemperature)
}
func (c *GameClient) onPeriod() {
	c.lock.Lock()
	if c.UserID != 0 {
		c.periodCheck()
	}
	c.lock.Unlock()
}

func (c *GameClient) periodCheck() {
	c.checkTasks()
	c.checkGuajiOutput()
	c.coolTemperature()
}

func (c *GameClient) checkGuajiOutput() {

	//	取出最后的一条产出的记录，根据时间判断是否进行产出的结算
	outPut := c.user.GuajiOutputBox.GuajiOutputs
	if len(outPut) != 0 {
		//当前的时间戳
		timestamp := time.Now().Unix()
		fmt.Println(timestamp)
		//上次结算的时间戳
		toBeCharge := outPut[len(outPut)-1].Time
		timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
		loc, _ := time.LoadLocation("Local")                            //重要：获取时区
		theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
		sr := theTime.Unix()                                            //转化为时间戳 类型是int64
		//如果两个的时间差小于10秒就不执行下面的代码
		if (timestamp - sr) < 10 {
			return
		}
	}

	//通过匹配等级和雇员的token值判定是否需要更新计算的产出的参数
	guajisettlement := c.user.GuajiSettlement
	//guajisettlement.MinLevel
	EmployeeInfo := gameconf.AllEmployees
	_ = EmployeeInfo
	c.PushSettlement()
	//根据guajisettlement来计算产出
	var size int = len(c.user.GuajiOutputBox.GuajiOutputs)
	if size == 100 {
		for i := 0; i < size-1; i++ {
			c.user.GuajiOutputBox.GuajiOutputs[i] = c.user.GuajiOutputBox.GuajiOutputs[i+1]
		}
		c.user.GuajiOutputBox.GuajiOutputs = append(c.user.GuajiOutputBox.GuajiOutputs[:99], c.user.GuajiOutputBox.GuajiOutputs[100:]...)
	}
	// a.	每10点运气值提高1%正向事件触发概率，降低1%负向事件触发概率（数值都可调整）
	// b.	最高可提高40%的正向事件触发概率，降低40%负向事件触发概率，无法再堆叠（数值都可调整）
	// c.	运气值提升的概率只作用于意外事件本来配置的%概率
	// 正向事件原有概率 + 运气值增加的概率 = 正向事件的总触发概率
	// 负向事件原有概率 - 运气值增加的概率 = 负向事件的总触发概率

	// 1）	[印刷质量 * 印刷速度 + 意外产出（负像为减少，正向为相加）]  = 挂机产出（金币）
	// 2）	每10分钟，计算一次挂机产出，并会自动收取
	// 3）	总产出最高可收入24小时的产出，超过则无法再挂机获得

	//计算本次结算的运气值，正向和负向的触发概率
	//长正向的概率
	var Probability1 int = guajisettlement.Probability1
	Probability1 += (guajisettlement.Luck / 10)
	//负向的概率
	var Probability2 int = guajisettlement.Probability2
	Probability2 -= (guajisettlement.Luck / 10)
	fmt.Println(Probability1)
	fmt.Println(Probability2)

	// //根据概率确定
	// var gailv [5]byte

	var n [3]int /* n 是一个长度为 3 的数组 */
	var j int
	n[0] = Probability1
	n[1] = Probability2
	n[2] = (100 - Probability1 - Probability2)
	rand.Seed(time.Now().Unix())
	var result int = 0
	_ = result
	var sum_all int = 100
	_ = sum_all
	/* 输出每个数组元素的值 */
	for j = 0; j < 3; j++ {
		rnd := rand.Intn(sum_all)

		fmt.Println(rnd)

		if rnd <= n[j] {
			result = j
			break
		} else {
			sum_all -= n[j]
		}
	}
	var coinNum int
	_ = coinNum

	fmt.Println("_______________")
	fmt.Println(guajisettlement.Quality)
	fmt.Println(guajisettlement.Speed)
	fmt.Println("_______________")

	coinNum = guajisettlement.Quality * guajisettlement.Speed
	oneEvent := &common.GuajiOutputInfo{}
	var gailv int

	if result == 1 {
		gailv = 100 - n[1]
		oneEvent.Type = "f"
	}
	if result == 0 {
		gailv = 100 + n[0]
		oneEvent.Type = "z"
	}

	if result == 2 {
		gailv = 100
		oneEvent.Type = "n"
	}
	coinNum = coinNum * gailv / 100
	oneEvent.UserID = c.user.UserID
	oneEvent.Name = c.user.Profile.Name
	coinNums := strconv.Itoa(coinNum)
	_ = coinNums
	oneEvent.Items = "产出" + coinNums + "金币"
	oneEvent.Time = time.Now().Format("2006-01-02 15:04:05")
	c.user.GuajiOutputBox.GuajiOutputs = append(c.user.GuajiOutputBox.GuajiOutputs, oneEvent)
	//用户金币数的增加
	c.user.GuajiProfile.Coin += coinNum
	c.persistGuajiProfile()
	fmt.Println(len(c.user.GuajiOutputBox.GuajiOutputs))
	c.persistOutput()
}

func (c *GameClient) persistGuajiProfile() {
	newGuajiProfile := &gameuser.User{}
	newGuajiProfile.UserID = c.UserID
	cpy := deepcopy.Copy(c.user.GuajiProfile)
	output, _ := cpy.(*gameuser.GuajiProfile)
	newGuajiProfile.GuajiProfile = output
	persistent.AddUser(c.UserID, newGuajiProfile)
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
	case message.MsgTypeClickOutputReq:
		return c.HandleClickOutputReq(metaData, rawMsg)
	case message.MsgTypeEmployeeListReq:
		return c.HandleEmployeeListReq(metaData, rawMsg)
	case message.MsgTypeEmployeeAdjustReq:
		return c.HandleEmployeeAdjustReq(metaData, rawMsg)
	case message.MsgTypeMyEmployeeReq:
		return c.HandleMyEmployeeReq(metaData, rawMsg)
	case message.MsgTypePickReq:
		return c.HandlePickReq(metaData, rawMsg)
	}

	return
}

// HandleEmployeeAdjustReq ...
func (c *GameClient) HandleMyEmployeeReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.MyEmployeeNotify{}
	reply.Meta.MessageType = "MyEmployeeNotify"
	reply.Meta.MessageTypeID = message.MsgMyEmployeeNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	var num int = len(c.user.GuajiProfile.EmployeeBox.EmployeesInfo)
	var numB int = len(c.user.Bag.BagEmployee)
	if num == 0 && numB == 0 {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "invalid request"
		c.SendMsg(reply)
		return
	}
	if num != 0 {
		//工作的
		for a := 0; a < numB; a++ {
			//	go func(who int) {
			myEmployee := &message.Employeeinfo{}
			numid := c.user.Bag.BagEmployee[a].EmployeesID
			myEmployee.Speed = gameconf.AllEmployees[numid].Speed
			myEmployee.Quality = gameconf.AllEmployees[numid].Quality
			myEmployee.Number = gameconf.AllEmployees[numid].Number
			myEmployee.Luck = gameconf.AllEmployees[numid].Luck
			myEmployee.Introdution = gameconf.AllEmployees[numid].Introdution
			myEmployee.EmployeeName = gameconf.AllEmployees[numid].EmployeeName
			myEmployee.AvatarImage = gameconf.AllEmployees[numid].AvatarImage
			reply.Data.EmployeeBack = append(reply.Data.EmployeeBack, myEmployee)
			//}(a)
		}
	}
	if numB != 0 {
		//背包
		for a := 0; a < num; a++ {
			//	go func(who int) {
			myEmployee := &message.Employeeinfo{}
			numid := c.user.GuajiProfile.EmployeeBox.EmployeesInfo[a].EmployeesID
			myEmployee.Speed = gameconf.AllEmployees[numid].Speed
			myEmployee.Quality = gameconf.AllEmployees[numid].Quality
			myEmployee.Number = gameconf.AllEmployees[numid].Number
			myEmployee.Luck = gameconf.AllEmployees[numid].Luck
			myEmployee.Introdution = gameconf.AllEmployees[numid].Introdution
			myEmployee.EmployeeName = gameconf.AllEmployees[numid].EmployeeName
			myEmployee.AvatarImage = gameconf.AllEmployees[numid].AvatarImage
			reply.Data.EmployeeWork = append(reply.Data.EmployeeWork, myEmployee)
			//}(a)
		}
	}
	c.SendMsg(reply)
	return
}

// HandleEmployeeAdjustReq ...
func (c *GameClient) HandleEmployeeAdjustReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.EmployeeAdjustNotify{}
	reply.Meta.MessageType = "EmployeeAdjustNotify"
	reply.Meta.MessageTypeID = message.MsgEmployeeAdjustNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	req := &message.SaveEmployeeAdjustReq{}

	fmt.Println(req.Data.EmployeeAdjust)

	err = json.Unmarshal(rawMsg, req)
	if err != nil {
		reply.Meta.Error = true
		reply.Meta.ErrorMessage = "invalid request"
		c.SendMsg(reply)
		return
	}
	//产生随机的标示值
	r := rand.New(rand.NewSource(time.Now().Unix()))
	fmt.Println(r.Intn(10000)) // [0,100)的随机值，返回值为int
	string := strconv.Itoa(r.Intn(10000))
	str := "token" + string

	cpy := deepcopy.Copy(req.Data.EmployeeAdjust)
	layout, _ := cpy.(*message.EmployeeAdjust)
	//加入工作中
	if len(layout.Employee) > 0 {
		c.user.GuajiProfile.EmployeeBox.EmployeesInfo = make([]*common.EmployeesInfo, 0)
		c.user.GuajiProfile.EmployeeBox.EmployeesToken = str
		for a := len(layout.Employee); a >= 1; a-- {
			go func(who int) {
				onet := &common.EmployeesInfo{}
				onet.EmployeesID = layout.Employee[who]
				fmt.Println("+++++++++++++++++++")
				fmt.Println(layout.Employee[who])
				c.user.GuajiProfile.EmployeeBox.EmployeesInfo = append(c.user.GuajiProfile.EmployeeBox.EmployeesInfo, onet)

				time.Sleep(10 * time.Nanosecond)
			}(a)
		}
		runtime.Gosched()
	}

	//加入背包
	if len(layout.Back) > 0 {
		c.user.Bag.BagEmployee = make([]*common.EmployeesInfo, 0)

		for a := len(layout.Back); a >= 1; a-- {
			go func(who int) {
				onet := &common.EmployeesInfo{}
				onet.EmployeesID = layout.Employee[who]
				fmt.Println("+++++++++++++++++++")
				fmt.Println(layout.Employee[who])
				c.user.Bag.BagEmployee = append(c.user.Bag.BagEmployee, onet)

				time.Sleep(10 * time.Nanosecond)
			}(a)
		}
		runtime.Gosched()
	}
	fmt.Println(len(layout.Employee))
	c.SendMsg(reply)
	c.persistEmployee()
	c.persistBagBox()
	return
}

func (c *GameClient) persistEmployee() {

	Employee := &gameuser.User{}

	cpy := deepcopy.Copy(c.user.GuajiProfile)
	adjust, _ := cpy.(*gameuser.GuajiProfile)
	Employee.GuajiProfile = adjust
	persistent.AddUser(c.UserID, Employee)

}

// HandleEmployeeListReq ...
func (c *GameClient) HandleEmployeeListReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	reply := &message.EmployeeListNotify{}
	reply.Meta.MessageType = "EmployeeListNotify"
	reply.Meta.MessageTypeID = message.MsgEmployeeListNotify
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	reply.Data.Employee = make([]*message.Employeeinfo, 0)
	fmt.Println(len(c.user.GuajiOutputBox.GuajiOutputs))

	for a := 1; a <= 10; a++ {
		//	go func(who int) {
		myEmployee := &message.Employeeinfo{}
		var str = "gy00"
		if a <= 9 {
			str = "gy00"
		} else {
			str = "gy0"
		}
		d := strconv.Itoa(a)
		str += d
		myEmployee.Speed = gameconf.AllEmployees[str].Speed
		myEmployee.Quality = gameconf.AllEmployees[str].Quality
		myEmployee.Number = gameconf.AllEmployees[str].Number
		myEmployee.Luck = gameconf.AllEmployees[str].Luck
		myEmployee.Introdution = gameconf.AllEmployees[str].Introdution
		myEmployee.EmployeeName = gameconf.AllEmployees[str].EmployeeName
		myEmployee.AvatarImage = gameconf.AllEmployees[str].AvatarImage
		reply.Data.Employee = append(reply.Data.Employee, myEmployee)
		//}(a)
	}
	c.SendMsg(reply)
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

	fmt.Println(c.user.GuajiOutputBox.GuajiOutputs)
	reply.Data.GuajiOutputs = make([]common.GuajiOutputInfo, 0)
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

// HandleClickOutputReq ...
func (c *GameClient) HandleClickOutputReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {

	ClickOutput := c.user.ClickOutputBox.ClickOutput
	ClickOutput.GoodID = "wp0001"
	ClickOutput.Time = int(time.Now().Unix())
	ClickOutput.Type = 0
	ClickOutput.UserID = c.user.UserID
	reply := &message.ClickOutputReq{}
	reply.Meta.MessageType = "ClickOutputReq"
	reply.Meta.MessageTypeID = message.MsgTypeClickOutputReq
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID

	// if len(c.user.GuajiOutputBox.GuajiOutputs) == 0 {
	// 	reply.Meta.Error = true
	// 	reply.Meta.ErrorMessage = "GuajiOutputs don't exits"
	// 	c.SendMsg(reply)
	// 	return
	// }
	//消息的推送
	//Events: = make([]common.EventInfo, 0)
	logic.RandOutputInfo(c.user)
	reply.Data.GoodID = c.user.ClickOutputBox.ClickOutput.GoodID
	reply.Data.Temperature = c.user.GuajiProfile.CurrentTemperature
	reply.Data.Num = c.user.ClickOutputBox.ClickOutput.GoodNum
	reply.Data.CD = c.user.GuajiProfile.CDTemperature
	reply.Data.Percent = c.user.GuajiProfile.TemperaturePercent
	// fmt.Println(len(c.user.GuajiOutputBox.GuajiOutputs))
	// for _, myOutputs := range c.user.GuajiOutputBox.GuajiOutputs {
	// 	reply.Data.GuajiOutputs = append(reply.Data.GuajiOutputs, *myOutputs)
	// }

	c.persistClikOutput()
	gameredis.PersistUser(c.user.UserID, c.user)
	fmt.Println((c.user.ClickOutputBox))
	c.SendMsg(reply)

	return
}
func (c *GameClient) persistClikOutput() {
	newUser := &gameuser.User{}
	newUser.UserID = c.UserID

	cpy := deepcopy.Copy(c.user.ClickOutputBox)
	output, _ := cpy.(*gameuser.ClickOutputBox)
	newUser.ClickOutputBox = output
	persistent.AddUser(c.UserID, newUser)
}

// HandlePickReq ...
func (c *GameClient) HandlePickReq(metaData message.ReqMetaData, rawMsg []byte) (err error) {
	fmt.Println("捡起")
	reply := &message.PickReq{}
	reply.Meta.MessageType = "PickReq"
	reply.Meta.MessageTypeID = message.MsgTypePickReq
	reply.Meta.MessageSequenceID = metaData.MessageSequenceID
	reply.Data.GoodID = c.user.ClickOutputBox.ClickOutput.GoodID
	reply.Data.Num = c.user.ClickOutputBox.ClickOutput.GoodNum
	reply.Data.Status = 1
	//5秒内捡起
	if c.user.GuajiProfile.CDPick > 0 {
		reply.Data.Status = 2
		c.persistPick()
		c.user.GuajiProfile.CDPick = 0
	}

	c.SendMsg(reply)
	return
}
func (c *GameClient) persistPick() {
	newUser := &gameuser.User{}
	newUser.UserID = c.UserID
	if c.user.Bag != nil {
		for _, v := range c.user.Bag.Cells {
			if v.GoodsID == c.user.ClickOutputBox.ClickOutput.GoodID {
				fmt.Println("你好啊啊啊啊啊")
				num, err := strconv.Atoi(c.user.ClickOutputBox.ClickOutput.GoodNum)
				if err == nil {
					v.Count += num
				}

			}
		}
	}
	cpy := deepcopy.Copy(c.user.Bag)
	newUserBag, _ := cpy.(*common.Bag)
	newUser.Bag = newUserBag
	persistent.AddUser(c.UserID, newUser)
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
	fmt.Println(c.user.GuajiProfile)
	fmt.Println(c.user.LoginTime.Time)
	fmt.Println("我是哈哈啊啊啊啊啊啊啊-----")
	//send message
	err = c.PushRoleInfo()
	if err != nil {
		return
	}
	//等级信息等的初始化
	err = c.InitializationInfo()
	if err != nil {
		return
	}
	err = c.LoadGuajiProfile()
	if err != nil {
		return
	}
	return
}

// InitializationInfo ...
func (c *GameClient) InitializationInfo() (err error) {

	//第一次就初始化等级为1
	if c.user.LoginTime.Time == "" {
		c.user.Profile.Level = 1
	} else {
		//不是第一次登陆，查看上次的登陆时间，如果差值大于一天，取上限24小时，否则，取上次的登陆时间进行运算
		//当前的时间戳
		timestampnow := time.Now().Unix()
		//上次登陆的时间戳
		toBeCharge := c.user.LoginTime.Time
		fmt.Println(toBeCharge)
		timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
		loc, _ := time.LoadLocation("Local")                            //重要：获取时区
		theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
		//var logintime int64 = 0
		logintime := theTime.Unix() //转化为时间戳 类型是int64
		//判断另个时间的差值
		var xunhuanTime int64 = 0
		_ = xunhuanTime
		if (timestampnow - logintime) >= 86400 {
			xunhuanTime = 86400 / 10
		} else {
			xunhuanTime = (timestampnow - logintime) / 10
		}
		var a int64 = 0
		//更新计算的暂存值
		c.PushSettlement()
		guajisettlement := c.user.GuajiSettlement
		for a = 0; a < xunhuanTime; a++ {
			timestamp := logintime + a*10
			//格式化为字符串,tm为Time类型
			tm := time.Unix(timestamp, 0)
			//根据guajisettlement来计算产出
			var size int = len(c.user.GuajiOutputBox.GuajiOutputs)
			if size == 100 {
				for i := 0; i < size-1; i++ {
					c.user.GuajiOutputBox.GuajiOutputs[i] = c.user.GuajiOutputBox.GuajiOutputs[i+1]
				}
				c.user.GuajiOutputBox.GuajiOutputs = append(c.user.GuajiOutputBox.GuajiOutputs[:99], c.user.GuajiOutputBox.GuajiOutputs[100:]...)
			}
			//计算本次结算的运气值，正向和负向的触发概率
			//长正向的概率
			var Probability1 int = guajisettlement.Probability1
			Probability1 += (guajisettlement.Luck / 10)
			//负向的概率
			var Probability2 int = guajisettlement.Probability2
			Probability2 -= (guajisettlement.Luck / 10)
			// //根据概率确定
			// var gailv [5]byte
			var n [3]int /* n 是一个长度为 3 的数组 */
			var j int
			n[0] = Probability1
			n[1] = Probability2
			n[2] = (100 - Probability1 - Probability2)
			rand.Seed(time.Now().Unix())
			var result int = 0
			_ = result
			var sum_all int = 100
			_ = sum_all
			/* 输出每个数组元素的值 */
			for j = 0; j < 3; j++ {
				rnd := rand.Intn(sum_all)
				if rnd <= n[j] {
					result = j
					break
				} else {
					sum_all -= n[j]
				}
			}
			var coinNum int
			_ = coinNum

			fmt.Println(guajisettlement.Quality)
			fmt.Println(guajisettlement.Speed)

			coinNum = guajisettlement.Quality * guajisettlement.Speed
			oneEvent := &common.GuajiOutputInfo{}
			var gailv int

			if result == 1 {
				gailv = 100 - n[1]
				oneEvent.Type = "f"
			}
			if result == 0 {
				gailv = 100 + n[0]
				oneEvent.Type = "z"
			}

			if result == 2 {
				gailv = 100
				oneEvent.Type = "n"
			}
			coinNum = coinNum * gailv / 100
			oneEvent.UserID = c.user.UserID
			oneEvent.Name = c.user.Profile.Name
			coinNums := strconv.Itoa(coinNum)
			_ = coinNums
			oneEvent.Items = "产出" + coinNums + "金币"
			oneEvent.Time = tm.Format("2006-01-02 15:04:05")
			c.user.GuajiOutputBox.GuajiOutputs = append(c.user.GuajiOutputBox.GuajiOutputs, oneEvent)
			fmt.Println(len(c.user.GuajiOutputBox.GuajiOutputs))
			//用户金币数的增加
			c.user.GuajiProfile.Coin += coinNum
		}
		c.persistGuajiProfile()
		c.persistOutput()
	}
	return
}

//  LoadGuajiProfile ...
func (c *GameClient) LoadGuajiProfile() (err error) {
	// redis取出温度计算温度
	timeNow := time.Now().Unix()
	if c.user.GuajiProfile.CDTemperature > 0 {
		//取出当前的等级
		userProfile := c.user.Profile
		//根据等级取出当前的机器的结算数据
		machine := gameconf.AllGuajis
		// 默认等级0+1
		machineInfo := machine[userProfile.Level+1]
		//取出机器温度
		c.user.GuajiProfile.CurrentTemperature = float64(machineInfo.InitialTemperature)
		if int(timeNow-c.user.GuajiProfile.ClickTime) >= c.user.GuajiProfile.CDTemperature {
			c.user.GuajiProfile.CurrentTemperature = float64(machineInfo.InitialTemperature)
		}
	} else {
		TimeDecr := timeNow - c.user.GuajiProfile.ClickTime
		//取出当前的等级
		userProfile := c.user.Profile
		//根据等级取出当前的机器的结算数据
		machine := gameconf.AllGuajis
		// 默认等级0+1
		machineInfo := machine[userProfile.Level+1]
		c.user.GuajiProfile.CurrentTemperature -= float64(TimeDecr / int64(machineInfo.CDPerDegree))
		if c.user.GuajiProfile.CurrentTemperature < float64(machineInfo.InitialTemperature) {
			c.user.GuajiProfile.CurrentTemperature = float64(machineInfo.InitialTemperature)
		}
		fmt.Println(c.user.GuajiProfile.CurrentTemperature)
		fmt.Println("")
	}
	return
}

// PushSettlement ...
func (c *GameClient) PushSettlement() (err error) {

	//取出当前的等级
	userProfile := c.user.Profile
	_ = userProfile
	//根据等级取出当前的机器的结算数据
	machine := gameconf.AllGuajis
	machineInfo := machine[userProfile.Level]
	//判断雇员的索引是为空，不为空的话，查出雇员的信息
	employer := c.user.GuajiProfile.EmployeeBox
	//通过匹配等级和雇员的token值判定是否需要更新计算的产出的参数
	guajisettlement := c.user.GuajiSettlement
	//guajisettlement.MinLevel
	EmployeeInfo := gameconf.AllEmployees
	if employer.EmployeesToken != guajisettlement.SettlementToken || guajisettlement.MachineLevel != userProfile.Level {
		//更新需要计算的数据
		//如果雇员的数量大于0，循环计算出雇员的计算值
		guajisettlement.Luck = 0

		guajisettlement.Quality = 0

		guajisettlement.Speed = 0
		if len(c.user.GuajiProfile.EmployeeBox.EmployeesInfo) > 0 {

			for _, myEmployer := range c.user.GuajiProfile.EmployeeBox.EmployeesInfo {

				fmt.Println("_______________")
				fmt.Println(EmployeeInfo[myEmployer.EmployeesID].Quality)
				fmt.Println(EmployeeInfo[myEmployer.EmployeesID].Speed)
				fmt.Println(EmployeeInfo[myEmployer.EmployeesID].Luck)

				fmt.Println("_______________")

				guajisettlement.Luck += EmployeeInfo[myEmployer.EmployeesID].Luck
				guajisettlement.Quality += EmployeeInfo[myEmployer.EmployeesID].Quality
				guajisettlement.Speed += EmployeeInfo[myEmployer.EmployeesID].Speed
			}

		}

		guajisettlement.Luck += machineInfo.Luck
		guajisettlement.MachineLevel = userProfile.Level
		guajisettlement.Quality += machineInfo.Quality
		guajisettlement.SettlementToken = employer.EmployeesToken
		guajisettlement.Speed += machineInfo.Speed
		guajisettlement.OppositeOutput = machineInfo.OppositeOutput
		guajisettlement.PositiveOutput = machineInfo.PositiveOutput
		guajisettlement.Probability1 = machineInfo.Probability1
		guajisettlement.Probability2 = machineInfo.Probability2
	}

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
