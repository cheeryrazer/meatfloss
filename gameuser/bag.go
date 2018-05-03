package gameuser

// BagCell ...
type BagCell struct {
	GoodsID  string `json:"goodsId"`
	Count    int    `json:"count"`
	UniqueID int64  `json:"uniqueId"` // 注, js只支持48位整形
}

// Bag ...
type Bag struct {
	Cells map[int64]*BagCell
}

// NewBag ...
func NewBag() *Bag {
	bag := &Bag{}
	bag.Cells = make(map[int64]*BagCell)
	return bag
}
