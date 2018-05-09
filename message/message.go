package message

import "meatfloss/common"

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
	} `json:"data"`
}

// MsgTypeGameBaseInfoNotify ...
const MsgTypeGameBaseInfoNotify int32 = 3300

// GameBaseInfoNotify ...
// 服务端推送角色信息
type GameBaseInfoNotify struct {
	Meta ReplyMetaData `json:"meta"`
	Data struct {
		Profile *RoleProfile  `json:"profile"`
		Bag     RoleBag       `json:"bag"`
		Tasks   []RunningTask `json:"tasks"`
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
