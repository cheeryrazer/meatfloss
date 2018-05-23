package gameuser

import (
	"meatfloss/common"
	"meatfloss/message"
	"sync"
)

// User ...
type User struct {
	Lock            sync.RWMutex
	UserID          int
	Profile         *Profile
	Bag             *common.Bag
	TaskBox         *TaskBox
	NewsBox         *NewsBox
	EventBox        *EventBox
	Layout          *message.ClientLayout
	LoginTime       *LoginTime
	GuajiOutputBox  *GuajiOutputBox
	ClickOutputBox  *ClickOutputBox
	GuajiSettlement *GuajiSettlement //挂机结算暂存数据
	GuajiProfile    *GuajiProfile    //挂机需要的配置的信息
}

// NewUser ...
func NewUser(userID int) (user *User) {
	user = &User{}
	user.UserID = userID
	user.Profile = NewProfile(userID)
	user.Bag = common.NewBagWithInitialData(userID)
	user.TaskBox = NewTaskBox(userID)
	user.NewsBox = NewNewsBox(userID)
	user.EventBox = NewEventBox(userID)
	user.Layout = message.NewClientLayout()
	user.LoginTime = NewLoginTime(userID)
	user.GuajiOutputBox = NewGuajiOutputBox(userID)
	user.ClickOutputBox = NewClickOutputBox(userID)
	user.GuajiSettlement = NewGuajiSettlement(userID)
	user.GuajiProfile = NewGuajiProfile(userID)
	return
}
