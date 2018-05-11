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

// NewBagWithInitialData ...
func NewBagWithInitialData(userID int) *Bag {
	bag := &Bag{}
	bag.UserID = userID
	bag.Cells = make(map[int64]*BagCell)

	{
		cell := &BagCell{}
		cell.Count = 5
		cell.GoodsID = "wp0001"
		cell.UniqueID = 100000 + 1
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &BagCell{}
		cell.Count = 10
		cell.GoodsID = "wp0002"
		cell.UniqueID = 100000 + 2
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &BagCell{}
		cell.Count = 3
		cell.GoodsID = "wp0003"
		cell.UniqueID = 100000 + 3
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &BagCell{}
		cell.Count = 1
		cell.GoodsID = "wp0004"
		cell.UniqueID = 100000 + 4
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &BagCell{}
		cell.Count = 5
		cell.GoodsID = "fs0001"
		cell.UniqueID = 200000 + 1
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &BagCell{}
		cell.Count = 10
		cell.GoodsID = "fs0002"
		cell.UniqueID = 200000 + 2
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &BagCell{}
		cell.Count = 3
		cell.GoodsID = "fs0003"
		cell.UniqueID = 200000 + 3
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &BagCell{}
		cell.Count = 1
		cell.GoodsID = "fs0004"
		cell.UniqueID = 200000 + 4
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &BagCell{}
		cell.Count = 1
		cell.GoodsID = "jj0001"
		cell.UniqueID = 300000 + 1
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &BagCell{}
		cell.Count = 1
		cell.GoodsID = "jj0002"
		cell.UniqueID = 300000 + 2
		bag.Cells[cell.UniqueID] = cell
	}

	return bag
}

// TaskInfo ...
type TaskInfo struct {
	TaskID    string
	Timestamp int
	PreTime   int
	UserID    int
	ID        int64
	NPCID     string
	Time      string
}
