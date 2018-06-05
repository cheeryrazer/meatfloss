package gameuser

import (
	"meatfloss/common"
)

// CollectionBox ...
type CollectionBox struct {
	UserID      int
	Collections []*common.Collections
}

// NewCollectionBox ...
func NewCollectionBox(userID int) (collection *CollectionBox) {
	collection = &CollectionBox{}
	collection.UserID = userID
	collection.Collections = make([]*common.Collections, 0)
	return
}
