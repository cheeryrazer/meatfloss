package gameuser

import "meatfloss/message"

// EventBox ...
type EventBox struct {
	UserID int
	Events map[string]*message.EventInfo
}

// NewEventBox ...
func NewEventBox(userID int) *EventBox {
	box := &EventBox{}
	box.UserID = userID
	box.Events = make(map[string]*message.EventInfo)
	return box
}
