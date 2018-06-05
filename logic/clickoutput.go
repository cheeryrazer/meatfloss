package logic

import (
	"fmt"
	"math/rand"
	"meatfloss/common"
	"meatfloss/gameconf"
	"meatfloss/gameuser"
	"strings"
	"time"

	"github.com/golang/glog"
)

// RandOutputInfo ...
func RandOutputInfo(c *gameuser.User) (err error) {
	c.GuajiProfile.ClickTime = time.Now().Unix()
	fmt.Println(c.GuajiProfile.CurrentTemperature)
	//取出当前的等级
	level := c.GuajiProfile.MachineLevel
	//根据等级取出当前的机器的结算数据
	machine := gameconf.AllGuajis
	// 默认等级0+1
	machineInfo := machine[level]
	//取出机器温度
	// currentTemperature := c.GuajiProfile.CurrentTemperature
	// 判断机器是否在cd

	if c.GuajiProfile.CDTemperature > 0 {
		return
	}
	//捡起时间+1s
	//增温
	CurrentTemperature := c.GuajiProfile.CurrentTemperature
	CurrentTemperature += float64(machineInfo.TemperaturePerClick)
	c.GuajiProfile.CurrentTemperature = CurrentTemperature
	c.GuajiProfile.TemperaturePercent = (CurrentTemperature / float64(machineInfo.MaxTemperature)) * float64(100)
	if c.GuajiProfile.TemperaturePercent > float64(100) {
		c.GuajiProfile.TemperaturePercent = float64(100)
	}
	fmt.Println(c.GuajiProfile.CurrentTemperature)
	if c.GuajiProfile.CurrentTemperature >= float64(machineInfo.MaxTemperature) {
		// c.GuajiProfile.CDTemperature =machineInfo.CD
		// 11111
		c.GuajiProfile.CDTemperature = 10
		c.ClickOutputBox.ClickOutput = &common.ClickOutputInfo{}
		c.GuajiProfile.CurrentTemperature = float64(machineInfo.MaxTemperature)
		return
	}

	// 暴击概率产出
	n := rand.Intn(100)
	Clickoutput := &common.ClickOutputInfo{}
	if n < machineInfo.CritProbability {
		goods := strings.Split(machineInfo.CritOutput, "|")
		ln := len(goods)
		// 随机取出一件物品
		goodIndex := rand.Intn(ln)

		good := goods[goodIndex]
		goodDetail := strings.Split(good, ";")

		Clickoutput.GoodID = goodDetail[0]
		Clickoutput.GoodNum = goodDetail[1]
		Clickoutput.Time = int(time.Now().Unix())
		Clickoutput.Type = 0
		Clickoutput.UserID = c.UserID
		Clickoutput.MessageSequenceID = c.GuajiProfile.MessageSequenceID
		c.ClickOutputBox.ClickOutput = Clickoutput
		c.ClickOutputBox.ClickOutputs = append(c.ClickOutputBox.ClickOutputs, Clickoutput)
		fmt.Println("暴击")
		fmt.Println(goodIndex)
	} else {
		goods := strings.Split(machineInfo.ClickOutput, "|")
		ln := len(goods)
		//随机取出一件物品
		goodIndex := rand.Intn(ln)
		good := goods[goodIndex]
		goodDetail := strings.Split(good, ";")
		Clickoutput.GoodID = goodDetail[0]
		Clickoutput.GoodNum = goodDetail[1]
		Clickoutput.Time = int(time.Now().Unix())
		Clickoutput.Type = 0
		Clickoutput.UserID = c.UserID
		Clickoutput.MessageSequenceID = c.GuajiProfile.MessageSequenceID
		c.ClickOutputBox.ClickOutput = Clickoutput
		c.ClickOutputBox.ClickOutputs = append(c.ClickOutputBox.ClickOutputs, Clickoutput)
		fmt.Println("非暴击")
		fmt.Println(goodIndex)
		fmt.Println(machineInfo.CritProbability)

	}
	fmt.Println("-----------点击")
	glog.Info(c.ClickOutputBox.ClickOutputs)

	return

}
