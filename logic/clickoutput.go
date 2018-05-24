// package logic

// import (
// 	"fmt"
// 	"math/rand"
// 	"meatfloss/common"
// 	"meatfloss/gameconf"
// 	"meatfloss/gameuser"
// 	"strings"
// 	"time"
// )

// // RandOutputInfo ...
// func RandOutputInfo(c *gameuser.User) (err error) {
// 	fmt.Println(c.GuajiProfile.CurrentTemperature)
// 	//取出当前的等级
// 	userProfile := c.Profile
// 	//根据等级取出当前的机器的结算数据
// 	machine := gameconf.AllGuajis
// 	// 默认等级0+1
// 	machineInfo := machine[userProfile.Level+1]
// 	//取出机器温度
// 	// currentTemperature := c.GuajiProfile.CurrentTemperature
// 	// 判断机器是否在cd

// 	if c.GuajiProfile.CDTemperature > 0 {
// 		c.ClickOutputBox.ClickOutput = &common.ClickOutputInfo{}
// 		return
// 	}
// 	//捡起时间+1s
// 	c.GuajiProfile.CDPick =5+1
// 	//增温
// 	CurrentTemperature:=c.GuajiProfile.CurrentTemperature
// 	CurrentTemperature+=float64(machineInfo.TemperaturePerClick)
// 	c.GuajiProfile.CurrentTemperature = CurrentTemperature
// 	c.GuajiProfile.TemperaturePercent = (CurrentTemperature/float64(machineInfo.MaxTemperature))*float64(100)
// 	if c.GuajiProfile.TemperaturePercent>float64(100){
// 		c.GuajiProfile.TemperaturePercent = float64(100)
// 	}
// 	fmt.Println(c.GuajiProfile.CurrentTemperature)
// 	if c.GuajiProfile.CurrentTemperature >= float64(machineInfo.MaxTemperature) {
// 		// c.GuajiProfile.CDTemperature =machineInfo.CD
// 		// 11111
// 		fmt.Println("suc")
// 		c.GuajiProfile.CDTemperature = 10
// 		c.ClickOutputBox.ClickOutput = &common.ClickOutputInfo{}
// 		c.GuajiProfile.CurrentTemperature = float64(machineInfo.MaxTemperature)
// 		return
// 	}

// 	// 暴击概率产出
// 	n := rand.Intn(100)
// 	if n < machineInfo.CritProbability {
// 		goods := strings.Split(machineInfo.CritOutput, "|")
// 		ln := len(goods)
// 		// 随机取出一件物品
// 		goodIndex := rand.Intn(ln)

// 		good := goods[goodIndex]
// 		goodDetail := strings.Split(good, ";")

// 		ClickOutputInfo := c.ClickOutputBox.ClickOutput
// 		ClickOutputInfo.GoodID = goodDetail[0]
// 		ClickOutputInfo.GoodNum = goodDetail[1]
// 		ClickOutputInfo.Time = int(time.Now().Unix())
// 		ClickOutputInfo.Type = 0
// 		ClickOutputInfo.UserID = c.UserID
// 		fmt.Println("暴击")
// 		fmt.Println(goodIndex)
// 	} else {
// 		goods := strings.Split(machineInfo.ClickOutput, "|")
// 		ln := len(goods)
// 		//随机取出一件物品
// 		goodIndex := rand.Intn(ln)
// 		good := goods[goodIndex]
// 		goodDetail := strings.Split(good, ";")

// 		ClickOutputInfo := c.ClickOutputBox.ClickOutput
// 		ClickOutputInfo.GoodID = goodDetail[0]
// 		ClickOutputInfo.GoodNum = goodDetail[1]
// 		ClickOutputInfo.Time = int(time.Now().Unix())
// 		ClickOutputInfo.Type = 0
// 		ClickOutputInfo.UserID = c.UserID
// 		fmt.Println("非暴击")
// 		fmt.Println(goodIndex)
// 		fmt.Println(machineInfo.CritProbability)

// 	}
// 	return

// }
