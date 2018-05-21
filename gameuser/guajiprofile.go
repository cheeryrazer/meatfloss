package gameuser

// GuajiProfile ...
type GuajiProfile struct {
	UserID    int    //用户的ID
	Employees string // 雇员
}

// NewGuajiProfile ...
func NewGuajiProfile(userID int) (guajiprofile *GuajiProfile) {
	guajiprofile = &GuajiProfile{}
	guajiprofile.UserID = userID
	return
}
