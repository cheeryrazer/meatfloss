package gameuser

// GuajiSettlement ...
type GuajiSettlement struct {
	UserID          int    //用户的ID
	MachineLevel    int    // 机器等级
	Speed           int    // 速度
	Quality         int    // 质量
	Luck            int    // 运气
	PositiveOutput  string // 正向事件产出
	Probability1    int    //触发概率1
	OppositeOutput  string // 负向事件产出
	Probability2    int    // 触发概率2
	SettlementToken string //匹配雇员的变化的token
}

// NewGuajiSettlement ...
func NewGuajiSettlement(userID int) (guaji *GuajiSettlement) {
	guaji = &GuajiSettlement{}
	guaji.UserID = userID
	return
}
