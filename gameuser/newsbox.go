package gameuser

// NewsBox ...
type NewsBox struct {
	UserID int
}

// NewNewsBox ...
func NewNewsBox(userID int) *NewsBox {
	box := &NewsBox{}
	box.UserID = userID
	return box
}
