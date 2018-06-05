package common

// BagCell ...
type BagCell struct {
	GoodsID  string `json:"goodsId"`
	Count    int    `json:"count"`
	UniqueID int64  `json:"uniqueId"` // 注, js只支持48位整形
}

// Bag ...
type Bag struct {
	UserID      int
	Cells       map[int64]*BagCell
	BagEmployee []*EmployeesInfo
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

// GuajiOutputInfo ...
type GuajiOutputInfo struct {
	UserID int    // 用户的id
	Name   string //用户的名字
	Type   string //产出的类型 z正向 f反向 n没有
	Time   string //产出时间
	Items  string //产出的物品
}

// ClickOutputInfo ...
type ClickOutputInfo struct {
	UserID            int    // 用户的id
	Type              int    // 0.点击产出 1.暴击产出
	Time              int    //产出时间
	GoodID            string //产出物品
	GoodNum           string // 产出数量
	MessageSequenceID int64  //唯一id
}

// EmployeesInfo ...
type EmployeesInfo struct {
	EmployeesID string // 雇员的ID
}

// RequiredInfo ...
type RequiredInfo struct {
	GoodID  string //需要的物品的ID
	GoodNum int    // 需要的物品的数量
}

// EndInfo ...
type EndInfo struct {
	GoodID  string //产出的物品的ID
	GoodNum int    // 产出的物品的数量
}

// Collections ...
type Collections struct {
	GoodID string //收藏的物品的ID
	Img    string //收藏的物品的图片
}

// Lattice ...
type Lattice struct {
	Time       int    // 剩余的时间
	Type       int    //1表示正在制作0表示空闲3表示未开启
	LatticeNum string //数据表对应的格子号
	Required   string //需要的材料
	End        string //制作的物品
	Coin       int    //解锁需要的钻石
	OrderID    int    //索引的ID
}
