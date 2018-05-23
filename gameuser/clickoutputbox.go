package gameuser

import (
	"meatfloss/common"
)

// ClickOutputBox ...
type ClickOutputBox struct {
	UserID       int // 用户的id
	ClickOutput *common.ClickOutputInfo // 点击产出
}

// NewClickOutputBox ...
func NewClickOutputBox (userID int) (c *ClickOutputBox) {
	c = &ClickOutputBox{}
	c.UserID = userID
	c.ClickOutput = &common.ClickOutputInfo {}
	return
}
