package gameuser

// GuajiProfile ...
type GuajiProfile struct {
	UserID      int          //用户的ID
	Coin        int          //金币
	EmployeeBox *EmployeeBox // 雇员
}

// NewGuajiProfile ...
func NewGuajiProfile(userID int) (guajiprofile *GuajiProfile) {
	guajiprofile = &GuajiProfile{}
	guajiprofile.UserID = userID
	guajiprofile.EmployeeBox = NewEmployeeBox()
	return
}
