package gameuser

// GuajiSettlement ...
type GuajiSettlement struct {
	UserID              int    //用户的ID
	Number              string // 编号
	MachineLevel        int    // 机器等级
	MinLevel            int    // 需要等级
	Speed               int    // 速度
	Quality             int    // 质量
	Luck                int    // 运气
	InitialTemperature  int    // 初始温度
	MaxTemperature      int    // 最高温度
	CDPerDegree         int    // 每度冷却时间（s)
	CD                  int    // 冷却时间
	TemperaturePerClick int    // 每次点击温度
	MachineImage        string // 机器图片
	NumEmployees        int    // 可雇佣数
	PositiveOutput      string // 正向事件产出
	Probability1        int    //触发概率1
	OppositeOutput      string // 负向事件产出
	Probability2        int    // 触发概率2
	ClickOutput         string // 每次点击产出
	CritProbability     int    // 暴击概率
	CritOutput          string // 暴击产出
	Upmaterial          string // 升级材料
	Uptime              int    // 升级时间
	CurrentTemperature 	int	   // 当前温度
}

// NewGuajiSettlement ...
func NewGuajiSettlement(userID int) (guaji *GuajiSettlement) {
	guaji = &GuajiSettlement{}
	guaji.UserID = userID
	return
}
