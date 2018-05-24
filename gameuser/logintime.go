package gameuser

// LoginTime ...
type LoginTime struct {
	UserID int // user id
	Time   string
}

// NewLoginTime ...
func NewLoginTime(userID int) (login *LoginTime) {
	login = &LoginTime{}
	login.UserID = userID

	return
}
