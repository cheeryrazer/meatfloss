package gameuser

// GuajiProfile ...
type GuajiProfile struct {
	UserID      int          //用户的ID
	EmployeeBox *EmployeeBox // 雇员
}

}
 // ceshi ...
// NewGuajiProfile ...
func NewGuajiProfile(userID int) (guajiprofile *GuajiProfile) {
	guajiprofile = &GuajiProfile{}
	guajiprofile.UserID = userID
	guajiprofile.EmployeeBox = NewEmployeeBox()
	return
}
