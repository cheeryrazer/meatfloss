package message

import (
	"meatfloss/common"
)

// ReqMetaData ...
type ReqMetaData struct {
	MessageType       string `json:"messageType"`       // 消息的名字， 填字符串
	MessageTypeID     int32  `json:"messageTypeId"`     // 消息的ID
	MessageSequenceID int32  `json:"messageSequenceId"` //客户端在一个会话里面保持自增即可
}

// ReqMeta ...
type ReqMeta struct {
	Meta ReqMetaData `json:"Meta"`
}

// MsgTypeLoginReq ...
// 登录请求
const MsgTypeLoginReq int32 = 1000

// LoginReq ...
// 登录请求内部的具体细节目前暂未定义
type LoginReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		Source  string `json:"source"`
		Account string `json:"account"`
		Token   string `json:"token"`
	} `json:"data"`
}

// MsgTypeWgLoginReq ...
// 外挂的登录请求
const MsgTypeWgLoginReq int32 = 1010

// WgLoginReq ...
// 登录请求内部的具体细节目前暂未定义
type WgLoginReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		Source  string `json:"source"`
		Account string `json:"account"`
		Token   string `json:"token"`
		Type    string `json:"type"`   //方便加数据测试加的   1就是正常的登陆   2就是外挂加数据的
		UserID  string `json:"userID"` //方便加数据测试加的  正常的登陆用户id为空，
	} `json:"data"`
}

// MsgTypeOutputReq ...
// 产出信息的请求
const MsgTypeOutputReq int32 = 2100

// OutputReq ...
// 前段用户的获取产出的信息
type OutputReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		OUTPUT string `json:"outPut"` // npcId 目前只有一个1
	} `json:"data"`
}

// MsgTypeEmployeeListReq ...
// 查看配置雇员的请求
const MsgTypeEmployeeListReq int32 = 1481

// EmployeeListReq ...
// 前端用户的获取雇员的信息系
type EmployeeListReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		EmployeeList string `json:"employeelist"` // npcId 目前只有一个1
	} `json:"data"`
}

// MsgTypeEmployeeAdjustReq ...
// 个人雇员的调整
const MsgTypeEmployeeAdjustReq int32 = 1482

// SaveEmployeeAdjustReq ...
// 个人雇员的调整的请求
type SaveEmployeeAdjustReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		EmployeeAdjust *EmployeeAdjust `json:"employeeadjust"` // npcId 目前只有一个1
	} `json:"data"`
}

// MsgTypeMyEmployeeReq ...
// 查看自己的雇员
const MsgTypeMyEmployeeReq int32 = 1483

// ShowMyEmployeeListReq ...
// 查看自己的雇员
type ShowMyEmployeeListReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		EmployeeList string `json:"employeelist"` // npcId 目前只有一个1
	} `json:"data"`
}

// MsgTypeMachineUpgradeReq ...
// 机器升级
const MsgTypeMachineUpgradeReq int32 = 1484

// ShowMachineUpgradeReq ...
// 机器升级的请求
type ShowMachineUpgradeReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		MachineUpgradeApply   string `json:"machineupgradeapply"`   //1升级的请求2不处理
		MachineUpgradeConfirm string `json:"machineupgradeconfirm"` //1升级的确认2不处理
	} `json:"data"`
}

// MsgTypeWgReq ...
// 外挂的请求
const MsgTypeWgReq int32 = 1888

// ShowWgReq ...
// 外挂的请求
type ShowWgReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		UserID  string `json:"UserID"`  //用户的ID
		GoodsID string `json:"goodsID"` //物品的ID
		Num     string `json:"num"`     //数量
	} `json:"data"`
}

// MsgTypeMakeLatticeReq ...
// 制作的请求
const MsgTypeMakeLatticeReq int32 = 2810

// MakeLatticeReq ...
// 制作页面的请求
type MakeLatticeReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		Res string `json:"res"` //请求
	} `json:"data"`
}

// MsgTypeMakingReq ...
// 制作的请求
const MsgTypeMakingReq int32 = 2812

// MakingReq ...
// 制作的请求
type MakingReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		Type    string `json:"type"`    //1添加2完成3解锁4加速
		GoodsID string `json:"goodsID"` //请求
		Lattice int    `json:"lattice"` //请求
		CoinNum int    `json:"coinNum"` //请求
	} `json:"data"`
}

// MsgTypeCollectionReq ...
// 收藏的请求
const MsgTypeCollectionReq int32 = 2811

// CollectionReq ...
// 收藏的请求
type CollectionReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		Type    string `json:"type"`    //1添加2取消
		GoodsID string `json:"goodsID"` //请求
		Img     string `json:"img"`     //请求
	} `json:"data"`
}

// +++++++++++++

// MsgTypeMarkNewsAsReadReq ...
// 将一个新闻标记为已读
const MsgTypeMarkNewsAsReadReq int32 = 1100

// MarkNewsAsReadReq ...
type MarkNewsAsReadReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		PushID    string `json:"pushId"`
		ArticleID string `json:"articleId"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeFinishEventReq ...
// 完成一个事件
const MsgTypeFinishEventReq int32 = 1200

// FinishEventReq ...
type FinishEventReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		EventGenID string `json:"eventGenId"`
		Choice     int    `json:"choice"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeCreateTaskReq ...
// 创建任务
const MsgTypeCreateTaskReq int32 = 1300

// CreateTaskReq ...
type CreateTaskReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		NPCID string `json:"npcId"` // npcId 目前只有一个1
	} `json:"data"`
}

// +++++++++++++

// MsgTypeSaveClientLayoutReq ...
// 创建任务
const MsgTypeSaveClientLayoutReq int32 = 1400

// SaveClientLayoutReq ...
type SaveClientLayoutReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		Layout *ClientLayout `json:"layout"` // npcId 目前只有一个1
	} `json:"data"`
}

// response
//--------------------------------------------------------------

// ReplyMetaData ...
type ReplyMetaData struct {
	MessageType       string `json:"messageType"`
	MessageTypeID     int32  `json:"messageTypeId"`     // 消息类型
	MessageSequenceID int32  `json:"messageSequenceId"` // 客户端传什么， 服务端回什么
	Error             bool   `json:"error"`             // true表示发生了错误
	ErrorMessage      string `json:"errorMessage"`      // 错误信息
}

// +++++++++++++

// MsgTypeLoginReply ...
// 登录请求
const MsgTypeLoginReply int32 = 2000

// LoginReply ...
//  登录的响应
type LoginReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

const MsgLoginInitReply int32 = 2500

type LoginInitReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Level              int     `json:"level"` // 等级
		Exp                int     `json:"exp"`   // 经验
		NextExp            int     `json:"nextexp"`
		Coin               int     `json:"coin"`    // 金币
		Diamond            int     `json:"diamond"` //钻石
		MachineLevel       int     `json:"machineLevel"`
		CDTemperature      int     `json:"cdTemperature"`
		Temperature        float64 `json:"temperature"`
		TemperaturePercent float64 `json:"temperaturePercent"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeWgLoginReply ...
// 外挂的登录请求
const MsgTypeWgLoginReply int32 = 2010

// WgLoginReply ...
//  外挂的登录的响应
type WgLoginReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeMakeLatticeReply ...
// 制作页面的请求
const MsgTypeMakeLatticeReply int32 = 2011

// MakeLatticeReply ...
//  制作页面的响应
type MakeLatticeReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeMakingReply ...
// 制作的请求
const MsgTypeMakingReply int32 = 2012

// MakingReply ...
//  制作的响应
type MakingReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeCollectionReply ...
// 收藏的请求
const MsgTypeCollectionReply int32 = 2012

// CollectionReply ...
//  制作的响应
type CollectionReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeOutputReply ...
// 产出报告请求
const MsgTypeOutputReply int32 = 5000

// OutPutReply ...
//  产出报告的响应
type OutPutReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeEmployeeListReply ...
// 雇员列表的请求
const MsgTypeEmployeeListReply int32 = 5600

// EmployeeListReply ...
//  查看所有可配置的雇员的列表
type EmployeeListReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeEmployeeAdjustReply ...
// 雇员列表的请求修改
const MsgTypeEmployeeAdjustReply int32 = 5700

// EmployeeAdjustReply ...
//  更改雇员配置
type EmployeeAdjustReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeMyEmployeeReply ...
// 我的雇员请请求查看
const MsgTypeMyEmployeeReply int32 = 5800

// MyEmployeeReply ...
//  查看自己的雇员
type MyEmployeeReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeMachineUpgradeReply ...
// 机器升级的请求
const MsgTypeMachineUpgradeReply int32 = 5900

//MachineUpgradeReply ...
//  机器的升级
type MachineUpgradeReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeWgReply ...
// 外挂的请求
const MsgTypeWgReply int32 = 5888

//WgReply ...
//  外挂的回复
type WgReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeMarkNewsAsReadReply ...
// 将一个新闻标记为已读的响应
const MsgTypeMarkNewsAsReadReply int32 = 2100

// MarkNewsAsReadReply ...
type MarkNewsAsReadReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		PushID    string `json:"pushId"`
		ArticleID string `json:"articleId"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeFinishEventReply ...
// 完成一个事件的响应
const MsgTypeFinishEventReply int32 = 2200

// FinishEventReply ...
type FinishEventReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		EventGenID string `json:"eventGenId"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeCreateTaskReply ...
const MsgTypeCreateTaskReply int32 = 2300

// CreateTaskReply ...
// 创建一个任务的响应
type CreateTaskReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		TaskID string `json:"taskId"`
	} `json:"data"`
}

// MsgTypeSaveClientLayoutReply ...
const MsgTypeSaveClientLayoutReply int32 = 2400

// SaveClientLayoutReply ...
// 创建一个任务的响应
type SaveClientLayoutReply struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// notifications
//--------------------------------------------------------------

// MsgTypeKickOffNotify ...
const MsgTypeKickOffNotify int32 = 3100

// KickOffNotify ...
// 通知这个用户被别人踢掉了
type KickOffNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Dummy string `json:"dummy"`
	} `json:"data"`
}

// MsgTypePushNewsNotify ...
const MsgTypePushNewsNotify int32 = 3100

// PushNewsNotify ...
//  服务端推送新闻
type PushNewsNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Articles []ArticleInfo `json:"articles"`
	} `json:"data"`
}

// MsgTypeEventNotify ...
const MsgTypeEventNotify int32 = 3200

// EventNotify ...
//  服务端推送时间
type EventNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Events []*EventInfo `json:"events"`
		UserID int          `json:"userId"`
	} `json:"data"`
}

// MsgTypeUpdateGoodsNotify ...
const MsgTypeUpdateGoodsNotify int32 = 3201

// UpdateGoodsNotify ...
// 推送物品信息
type UpdateGoodsNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		List []GoodsUpdateInfo `json:"list"`
		Type string            `json:"type"`
	} `json:"data"`
}

// MsgTypeGameBaseInfoNotify ...
const MsgTypeGameBaseInfoNotify int32 = 3300

// GameBaseInfoNotify ...
// 服务端推送角色信息
type GameBaseInfoNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Profile *RoleProfile      `json:"profile"`
		Bag     RoleBag           `json:"bag"`
		Tasks   []common.TaskInfo `json:"tasks"`
		Events  []EventInfo       `json:"events"`
		Layout  *ClientLayout     `json:"layout"`
	} `json:"data"`
}

// MsgTypeNPCGuestNotify ...
const MsgTypeNPCGuestNotify int32 = 3301

// NPCGuestNotify ...
type NPCGuestNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		NPCList []string `json:"npcList"`
	} `json:"data"`
}

// MsgTypeTaskNotify ...
const MsgTypeTaskNotify int32 = 3400

// TaskNotify ...
//  服务端推送时间
type TaskNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Tasks  []*common.TaskInfo `json:"tasks"`
		UserID int                `json:"userId"`
	} `json:"data"`
}

// MsgTypeOutputNotify ...
const MsgTypeOutputNotify int32 = 3500

// OutputNotify ...
//  结算信息的推送
type OutputNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		GuajiOutputs []*common.GuajiOutputInfo `json:"outputs"`
	} `json:"data"`
}

// MsgEmployeeListNotify ...
const MsgEmployeeListNotify int32 = 3560

// EmployeeListNotify ...
//  可配置雇员的推送
type EmployeeListNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Employee []*Employeeinfo `json:"employeesinfo"`
	} `json:"data"`
}

// MsgEmployeeAdjustNotify ...
const MsgEmployeeAdjustNotify int32 = 3570

// EmployeeAdjustNotify ...
//  雇员调整的推送
type EmployeeAdjustNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		EmployeeWork []*Employeeinfo        `json:"employeeWork"`
		EmployeeBack []*Employeeinfo        `json:"employeeBack"`
		Machine      []*RoleGuajiSettlement `json:"machine"`
	} `json:"data"`
}

// MsgWgtNotify ...
const MsgWgtNotify int32 = 3888

// WgNotify ...
//  外挂的推送
type WgNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Res string `json:"res"`
	} `json:"data"`
}

// MsgMakeLatticeNotify ...
const MsgMakeLatticeNotify int32 = 3910

// MakeLatticNotify ...
//  制作的推送
type MakeLatticNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		MakeLatticeBack *MakeLatticeBack      `json:"makeLatticeBack"`
		Collection      []*common.Collections `json:"collection"`
		Lattice         []*common.Lattice     `json:"lattice"`
	} `json:"data"`
}

// MsgMakingNotify ...
const MsgMakingNotify int32 = 3911

// MakingNotify ...
//  制作的推送
type MakingNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Lattice         []*common.Lattice     `json:"lattice"`
		MakeLatticeBack *MakeLatticeBack      `json:"makeLatticeBack"`
		Collection      []*common.Collections `json:"collection"`
	} `json:"data"`
}

// MsgCollectionNotify ...
const MsgCollectionNotify int32 = 3581

// CollectionNotify ...
//  收藏的推送
type CollectionNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Collection []*common.Collections `json:"collection"`
	} `json:"data"`
}

// MsgMyEmployeeNotify ...
const MsgMyEmployeeNotify int32 = 3580

// MyEmployeeNotify ...
//  自己雇员的推送
type MyEmployeeNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		EmployeeWork []*Employeeinfo        `json:"employeeWork"`
		EmployeeBack []*Employeeinfo        `json:"employeeBack"`
		Machine      []*RoleGuajiSettlement `json:"machine"`
	} `json:"data"`
}

// MsgMyUpgradeNotify ...
const MsgMyUpgradeNotify int32 = 3810

// UpgradeNotify ...
//  机器升级的主动推送
type UpgradeNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		// EmployeeWork []*Employeeinfo        `json:"employeeWork"`
		// EmployeeBack []*Employeeinfo        `json:"employeeBack"`
		// Machine      []*RoleGuajiSettlement `json:"machine"`
		Upgrade     string `json:"upgrade"`
		UpgradeTime int    `json:"upgradeTime"`
	} `json:"data"`
}

// MsgMyCoolingNotify ...
const MsgMyCoolingNotify int32 = 3811

// CooliNotify ...
//  机器降温的主动推送的
type CooliNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Upgrade         string `json:"upgrade"`
		UpgradeTime     int    `json:"upgradeTime"`
		InitTemperature int    `json:"inittemperature"`
	} `json:"data"`
}

// MsgMachineUpgradeNotify ...
const MsgMachineUpgradeNotify int32 = 3581

// MachineUpgradeNotify ...
//  机器确认升级的推送
type MachineUpgradeNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		MachineNow            []*RoleGuajiSettlement `json:"machinenow"`
		MachineUpgrade        []*RoleGuajiSettlement `json:"machineupgrade"`
		MachineUpgradeType    string                 `json:"machineupgradetype"` //yes可升级no不可升级
		MachineUpgradeMessage string                 `json:"machineupgrademessage"`
	} `json:"data"`
}

// +++++++++++++

// MsgTypeClickOutputReq ...
// 点击产出

// ReplyClickOutputReq ...
type ReplyClickOutputReq struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Output struct {
			GoodID            string  `json:"goodId"`            // 物品Id
			Temperature       float64 `json:"temperature"`       // 温度
			Num               string  `json:"num"`               // 物品数
			CD                int     `json:"cd"`                // 机器过热cd
			Percent           float64 `json:"percent"`           // 温度百分比
			MessageSequenceID int64   `json:"messageSequenceId"` //客户端在一个会话里面保持自增即可
		} `json:"output"`
	} `json:"data"`
}

const MsgTypeClickOutputReq int32 = 6000

// ClickOutputReq ...
type ClickOutputReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		Output struct {
			GoodID      string  `json:"goodId"`      // 物品Id
			Temperature float64 `json:"temperature"` // 温度
			Num         string  `json:"num"`         // 物品数
			CD          int     `json:"cd"`          // 机器过热cd
			Percent     float64 `json:"percent"`     // 温度百分比
		} `json:"output"`
	} `json:"data"`
}

//  ReplyPickReq ...
type ReplyPickReq struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Status struct {
			GoodID            string `json:"goodId"`            // 物品Id
			Num               int    `json:"num"`               // 物品数
			Status            int    `json:"status"`            // 状态
			MessageSequenceID int64  `json:"messageSequenceId"` //客户端在一个会话里面保持自增即可
		} `json:"status"`
	} `json:"data"`
}

// MsgTypePickReq ...
//捡起
const MsgTypePickReq int32 = 7300

// PickReq ...
type PickReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		GoodID string `json:"goodId"` // npcId 目前只有一个1
		Num    int    `json:"num"`    // 物品数
	} `json:"data"`
}

// ClickStatusReq ...
type ClickStatusReq struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Status struct {
			Status            int   `json:"status"`            // 物品Id
			MessageSequenceID int64 `json:"messageSequenceId"` //客户端在一个会话里面保持自增即可
		} `json:"status"`
	} `json:"data"`
}

// ReplyUserNotify ...
type ReplyUserNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Level   int `json:"level"` // 等级
		Exp     int `json:"exp"`   // 经验
		NextExp int `json:"nextexp"`
		Coin    int `json:"coin"`    // 金币
		Diamond int `json:"diamond"` //钻石
	} `json:"data"`
}

// MsgTypePickCoinReq ...
const MsgTypePickCoinReq int32 = 7400

// PickCoinReq ...
type PickCoinReq struct {
	MetaData ReqMetaData `json:"meta"`
	Data     struct {
		Num  int `json:"num"`
		Type int `json:"type"` //1.金币2钻石
	} `json:"data"`
}

// ReplyPickCoinReq ...
type ReplyPickCoinReq struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Coin int `json:"coin"`
	} `json:"data"`
}

// +++++++++++++
// common structs
// 下面的都是子结构体
//--------------------------------------------------------------

// RoleBag ...
type RoleBag struct {
	Cells []CellInfo `json:"cells"`
}

// CellInfo ...
type CellInfo struct {
	GoodsID  string `json:"goodsId"`
	Count    int    `json:"count"`
	UniqueID int64  `json:"uniqueId"`
}

// RoleGoods ...
type RoleGoods struct {
	Index int    `json:"index"`
	ID    string `json:"id"`
	Count int    `json:"count"`
}

// RoleProfile ...
type RoleProfile struct {
	UserID       int    `json:"userId"`
	Name         string `json:"name"`
	Gender       int    `json:"gender"`
	Level        int    `json:"level"`
	Spine        string `json:"spine"`
	Intelligence int    `json:"intelligence"`
	Intimacy     int    `json:"intimacy"`
	Stamina      int    `json:"stamina"`
	Experience   int    `json:"experience"`
}

// RoleGuajiSettlement ...
type RoleGuajiSettlement struct {
	UserID              int    //用户的ID
	Number              string // 编号
	MachineLevel        int    // 机器等级
	MinLevel            int    // 需要等级
	Speed               int    // 速度
	Quality             int    // 质量
	Luck                int    // 运气
	InitialTemperature  int    // 初始温度
	MaxTemperature      int    // 最高温度
	CDPerDegree         int    // 每度冷却时间（s)
	CD                  int    // 冷却时间
	TemperaturePerClick int    // 每次点击温度
	MachineImage        string // 机器图片
	NumEmployees        int    // 可雇佣数
	PositiveOutput      string // 正向事件产出
	Probability1        int    //触发概率1
	OppositeOutput      string // 负向事件产出
	Probability2        int    // 触发概率2
	ClickOutput         string // 每次点击产出
	CritProbability     int    // 暴击概率
	CritOutput          string // 暴击产出
	Upmaterial          string // 升级材料
	Uptime              int    // 升级时间
}

// RoleGuajiProfile ...
type RoleGuajiProfile struct {
	UserID    int    //用户的ID
	Employees string // 雇员
}

// ArticleInfo ...
type ArticleInfo struct {
	PushID    string `json:"pushId"`
	ArticleID string `json:"articleId"`
	PicURL    string `json:"picurl"`
	Tags      string `json:"tags"`
	Title     string `json:"title"`
}

// EventInfo ...
type EventInfo struct {
	Type    string `json:"type"` // normal普通事件， select 选择事件
	EventID string `json:"eventId"`
	GenID   string `json:"genId"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    string `json:"time"`
	UserID  int    `json:"userId"`
}

// SingleReward ...
type SingleReward struct {
	GoodsID  string `json:"goodsId"`
	GoodsNum int    `json:"goodsNum"`
}

// Reward ...
type Reward struct {
	List []SingleReward `json:"list"`
	Exp  int            `json:"exp"`
}

// GoodsUpdateInfo ...
type GoodsUpdateInfo struct {
	GoodsID       string `json:"goodsId"`
	GoodsNumDelta int    `json:"goodsNumDelta"`
	GoodsNum      int    `json:"goodsNum"`
	UniqueID      int64  `json:"uniqueId"`
}

// RunningTask ...
type RunningTask struct {
	TaskID   string //
	CreateAt string
	PreTime  int
}

// ClientLayout ...
type ClientLayout struct {
	Floor1 map[string]string `json:"floor1"`
	Floor2 map[string]string `json:"floor2"`
	Floor3 map[string]string `json:"floor3"`
	Dress  map[string]string `json:"dress"`
}

// EmployeeAdjust ...
type EmployeeAdjust struct {
	Employee map[int]string `json:"employee"`
	Back     map[int]string `json:"back"`
}

// NewClientLayout ...
func NewClientLayout() *ClientLayout {
	layout := &ClientLayout{}
	layout.Floor1 = make(map[string]string)
	layout.Floor2 = make(map[string]string)
	layout.Floor3 = make(map[string]string)
	layout.Dress = make(map[string]string)
	return layout
}

// SingleGuaji ...
type SingleGuaji struct {
	GoodsID  string `json:"goodsId"`
	GoodsNum int    `json:"goodsNum"`
}

// Guaji ...
type Guaji struct {
	List []SingleGuaji `json:"list"`
}

// Employeeinfo ...
type Employeeinfo struct {
	Number       string // 编号
	AvatarImage  string // 头像图片
	EmployeeName string // 雇员名字
	Speed        int    // 速度
	Quality      int    // 质量
	Luck         int    // 运气
	Introdution  string // 介绍
}

// EmployeeinfoId ...
type EmployeeinfoId struct {
	Number string // 编号
}

// EmployeeinfoId1 ...
type EmployeeinfoId1 struct {
	Number string // 编号
}

// // ClickOutputs ...
// type ClickOutputs struct {

// }

// MakeLatticeBack ...
type MakeLatticeBack struct {
	Furniture map[string]string `json:"furniture"`
	Clothes   map[string]string `json:"clothes"`
}
