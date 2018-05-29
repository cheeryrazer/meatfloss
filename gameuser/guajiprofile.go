package gameuser

// GuajiProfile ...
type GuajiProfile struct {
	UserID             int          //用户的ID
	Coin               int          //金币
	MachineLevel       int          // 机器等级
	UpgradeTime        int          //升级剩余事件
	Upgrade            int          //1升级中2不是升级中3表示未开启升级
	EmployeeBox        *EmployeeBox // 雇员
	CurrentTemperature float64      //当前温度
	TemperaturePercent float64      //温度百分比
	CDTemperature      int          //剩余冷却时间
	CDPick             int          //剩余捡起时间
	ClickTime          int64        //温度产生时间
	MessageSequenceID  int64          //唯一id
}

// NewGuajiProfile ...
func NewGuajiProfile(userID int) (guajiprofile *GuajiProfile) {
	guajiprofile = &GuajiProfile{}
	guajiprofile.UserID = userID
	guajiprofile.EmployeeBox = NewEmployeeBox()
	return
}
