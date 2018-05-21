package gameuser

import (
	"meatfloss/common"
)

// GuajiOutputBox ...
type GuajiOutputBox struct {
	UserID       int // 用户的id
	GuajiOutputs []*common.GuajiOutputInfo
}

// NewGuajiOutputBox ...
func NewGuajiOutputBox(userID int) (guajioutputbox *GuajiOutputBox) {
	guajioutputbox = &GuajiOutputBox{}
	guajioutputbox.UserID = userID
	guajioutputbox.GuajiOutputs = make([]*common.GuajiOutputInfo, 0)
	return
}
