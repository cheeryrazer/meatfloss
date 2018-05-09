package common

// BagCell ...
type BagCell struct {
	GoodsID  string `json:"goodsId"`
	Count    int    `json:"count"`
	UniqueID int64  `json:"uniqueId"` // 注, js只支持48位整形
}

// Bag ...
type Bag struct {
	UserID int
	Cells  map[int64]*BagCell
}

// NewBag ...
func NewBag(userID int) *Bag {
	bag := &Bag{}
	bag.UserID = userID
	bag.Cells = make(map[int64]*BagCell)
	return bag
}
