package gameuser

import (
	"meatfloss/common"
)

// ClickOutputBox ...
type ClickOutputBox struct {
	UserID       int                       // 用户的id
	ClickOutputs []*common.ClickOutputInfo // 点击产出
	ClickOutput  *common.ClickOutputInfo
}

// NewClickOutputBox ...
func NewClickOutputBox(userID int) (c *ClickOutputBox) {
	c = &ClickOutputBox{}
	c.UserID = userID
	c.ClickOutputs = make([]*common.ClickOutputInfo, 0)
	c.ClickOutput = &common.ClickOutputInfo{}
	return
}
