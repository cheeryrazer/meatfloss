package logic

import(
	"fmt"
	"meatfloss/common"
	"meatfloss/gameconf"
	"meatfloss/gameuser"
	"math/rand"
)
// RandOutputInfo ...
func RandOutputInfo(c *gameuser.User)  *common.ClickOutputInfo {

	gl:=gameconf.AllGuajis
	gl[c.GuajiSettlement.MachineLevel].CritProbability=15
	// 暴击概率产出
	n:=rand.Intn(100)
	if n<gl[c.GuajiSettlement.MachineLevel].CritProbability{

	} else {
		fmt.Println("非暴击")
	}
}
func (c *gameuser.User) cooling
