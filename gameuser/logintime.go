package gameuser

import "time"

// LoginTime ...
type LoginTime struct {
	UserID int // user id
	Time   string
}

// NewLoginTime ...
func NewLoginTime(userID int) (login *LoginTime) {
	login = &LoginTime{}
	login.UserID = userID
	login.Time = time.Now().Format("2006-01-02 15:04:05")
	return
}
