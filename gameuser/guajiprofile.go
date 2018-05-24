package gameuser
// GuajiProfile ...
type GuajiProfile struct {
	UserID      int          //用户的ID
	Coin        int          //金币
	EmployeeBox *EmployeeBox // 雇员
	CurrentTemperature	float64	//当前温度
	TemperaturePercent	float64 //温度百分比
	CDTemperature		int //剩余冷却时间
	CDPick 				int //剩余捡起时间


}

// NewGuajiProfile ...
func NewGuajiProfile(userID int) (guajiprofile *GuajiProfile) {
	guajiprofile = &GuajiProfile{}
	guajiprofile.UserID = userID
	guajiprofile.EmployeeBox = NewEmployeeBox()
	return
}
