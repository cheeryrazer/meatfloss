package gameuser

// EventBox ...
type EventBox struct {
	UserID int
}

// NewEventBox ...
func NewEventBox(userID int) *EventBox {
	box := &EventBox{}
	box.UserID = userID
	return box
}
