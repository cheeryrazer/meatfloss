package gameconf

import (
	"encoding/json"
	"fmt"
	"meatfloss/db"
	"meatfloss/message"
	"meatfloss/utils"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

var (
	// AllSuperGoods ...
	AllSuperGoods map[string]*SuperGoods
	// AllGoods ...
	AllGoods map[string]*Goods
	// AllApparels ...
	AllApparels map[string]*Apparel
	// AllFurniture ... note, furniture is uncountable noun
	AllFurniture map[string]*Furniture
	// AllNPCs ...
	AllNPCs map[string]*NPC
	// AllTasks ...
	AllTasks map[string]*Task
	// AllRandomEvents ...
	AllRandomEvents map[string]*RandomEvent
	// AllNPCGuests ...
	AllNPCGuests map[string]*NPCGuest
	// AllGuajis ...
	AllGuajis map[int]*Guaji
	// AllEmployees ...
	AllEmployees map[string]*Employee
	// AllHierarchical ...
	AllHierarchical map[int]*Hierarchical
	// AllLattice ...
	AllLattice map[int]*Lattice
	// AllConfige ...
	AllConfige map[int]*Confige
)

func init() {
	AllSuperGoods = make(map[string]*SuperGoods)
	AllGoods = make(map[string]*Goods)
	AllApparels = make(map[string]*Apparel)
	AllFurniture = make(map[string]*Furniture)
	AllNPCs = make(map[string]*NPC)
	AllTasks = make(map[string]*Task)
	AllRandomEvents = make(map[string]*RandomEvent)
	AllNPCGuests = make(map[string]*NPCGuest)
	AllGuajis = make(map[int]*Guaji)
	AllEmployees = make(map[string]*Employee)
	AllHierarchical = make(map[int]*Hierarchical)
	AllLattice = make(map[int]*Lattice)
	AllConfige = make(map[int]*Confige)
}

// SuperGoods ...
type SuperGoods struct {
	ID          string
	Name        string
	UniqueID    int64
	AllowPileup int
}

// Goods ...
type Goods struct {
	ID                 string
	Type               int
	OrderID            int
	Consumable         int
	ImageName          string
	ImageEffect        int
	Name               string
	Description        string
	MinLevel           int
	DesignerMinLevel   int
	CanBeSold          int
	PriceForSail       int
	IntelligenceGain   int
	StaminaGain        int
	ExperienceGain     int
	FriendlyDegreeGain int
	AllowPileup        int
	UniqueID           int64
}

// LoadFromDatabase ...
func LoadFromDatabase() (err error) {
	loadGoods()
	loadApparel()

	loadNPCs()
	loadTasks()
	loadRandomEvents()
	loadNPCGuest()
	loadGuaji()
	loadEmployee()
	loadHierarchical()
	loadLattice()
	loadFurniture()
	loadConfige()
	// 这个放在最后
	loadSuperGoods()
	return
}

func loadSuperGoods() {
	for _, v := range AllGoods {
		obj := &SuperGoods{}
		obj.ID = v.ID
		obj.Name = v.Name
		obj.UniqueID = v.UniqueID
		obj.AllowPileup = v.AllowPileup
		AllSuperGoods[obj.ID] = obj
	}

	for _, v := range AllApparels {
		aa := &SuperGoods{}
		aa.ID = v.ID
		aa.Name = v.Name
		aa.UniqueID = v.UniqueID
		aa.AllowPileup = v.AllowPileup
		AllSuperGoods[aa.ID] = aa
	}

	for _, v := range AllFurniture {
		obj := &SuperGoods{}
		obj.ID = v.ID
		obj.Name = v.Name
		obj.UniqueID = v.UniqueID
		obj.AllowPileup = v.AllowPileup
		fmt.Println("%v", obj)
		AllSuperGoods[obj.ID] = obj
	}
}

func loadGoods() (err error) {
	dbGoods, err := db.LoadGoodsConf()
	if err != nil {
		return
	}
	_ = dbGoods
	for _, row := range dbGoods {
		goods := &Goods{
			ID:                 row.ID,
			Type:               row.Type,
			OrderID:            row.OrderID,
			Consumable:         row.Consumable,
			ImageName:          row.ImageName,
			ImageEffect:        row.ImageEffect,
			Name:               row.Name,
			Description:        row.Description,
			MinLevel:           row.MinLevel,
			DesignerMinLevel:   row.DesignerMinLevel,
			CanBeSold:          row.CanBeSold,
			PriceForSail:       row.PriceForSail,
			IntelligenceGain:   row.IntelligenceGain,
			StaminaGain:        row.StaminaGain,
			ExperienceGain:     row.ExperienceGain,
			FriendlyDegreeGain: row.FriendlyDegreeGain,
			AllowPileup:        row.AllowPileup}

		temp := strings.Replace(goods.ID, "wp", "", -1)
		goods.UniqueID, _ = strconv.ParseInt(temp, 10, 64)
		goods.UniqueID += 100000
		AllGoods[goods.ID] = goods
	}
	bin, err := json.MarshalIndent(AllGoods, "", " ")
	if err != nil {
		glog.Error("json.Marshal(config.Get()) failed")
	}
	text := string(bin)
	_ = text
	// fmt.Println(text)
	return
}

// NPC ...
type NPC struct {
	ID               string
	Name             string
	Description      string
	Gender           int
	Spine            string
	DecorationType   int
	GuestDecoration  int
	Intimacy         int
	Intelligence     int
	Stamina          int
	GuestProbability int
}

func loadNPCs() (err error) {
	dbNPCs, err := db.LoadNPCConf()
	if err != nil {
		return
	}
	for _, row := range dbNPCs {
		obj := &NPC{
			row.ID,
			row.Name,
			row.Description,
			row.Gender,
			row.Spine,
			row.DecorationType,
			row.GuestDecoration,
			row.Intimacy,
			row.Intelligence,
			row.Stamina,
			row.GuestProbability}
		AllNPCs[obj.ID] = obj
	}
	bin, err := json.MarshalIndent(AllNPCs, "", " ")
	if err != nil {
		glog.Error("json.Marshal(config.Get()) failed")
	}
	text := string(bin)
	_ = text
	// fmt.Println(text)
	return
}

// Task ...
type Task struct {
	ID               string
	Stars            int
	Type             int
	MinLevel         int
	Npc              string
	Intelligence     int
	Stamina          int
	FriendlyDegree   int
	DailyTriggerNum  int
	TotalTriggerNum  int
	Probability      int
	AssociationGroup int
	TriggerOrder     int
	IntimacyNpc      int
	IntimacyGain     int
	Image            string
	Description      string
	Choice1          string
	Choice2          string
	Choice3          string
	Reward1          string
	Exp1             int
	Reward2          string
	Exp2             int
	Reward3          string
	Exp3             int
	PreTime          int
	PostTime         int
	Choices          []string
	Rewards          []message.Reward
}

func loadTasks() (err error) {
	dbObjs, err := db.LoadTaskConf()
	if err != nil {
		return
	}
	for _, row := range dbObjs {
		obj := &Task{
			row.ID,
			row.Stars,
			row.Type,
			row.MinLevel,
			row.Npc,
			row.Intelligence,
			row.Stamina,
			row.FriendlyDegree,
			row.DailyTriggerNum,
			row.TotalTriggerNum,
			row.Probability,
			row.AssociationGroup,
			row.TriggerOrder,
			row.IntimacyNpc,
			row.IntimacyGain,
			row.Image,
			row.Description,
			row.Choice1,
			row.Choice2,
			row.Choice3,
			row.Reward1,
			row.Exp1,
			row.Reward2,
			row.Exp2,
			row.Reward3,
			row.Exp3,
			row.PreTime,
			row.PostTime,

			make([]string, 3),
			make([]message.Reward, 3)}

		rewardStrList := []string{row.Reward1, row.Reward2, row.Reward3}
		expList := []int{row.Exp1, row.Exp2, row.Exp3}
		for i, str := range rewardStrList {
			_ = str
			// for examples, str = wp0001;1000|wp0002;1000
			oneReward := message.Reward{}
			ones := strings.Split(str, "|")
			for _, one := range ones {
				// for example, one =  wp0001;1000
				twos := strings.Split(one, ";")
				if len(twos) == 2 {
					goodsID := strings.TrimSpace(twos[0])
					if goodsID != "" {
						goodsNum, err := strconv.Atoi(twos[1])
						if err == nil {
							sw := message.SingleReward{goodsID, goodsNum}
							oneReward.List = append(oneReward.List, sw)
						}
					}
				}
			}
			oneReward.Exp = expList[i]
			obj.Rewards[i] = oneReward
		}
		obj.Choices[0] = obj.Choice1
		obj.Choices[1] = obj.Choice2
		obj.Choices[2] = obj.Choice3

		AllTasks[obj.ID] = obj
	}
	// bin, err := json.MarshalIndent(AllTasks, "", " ")
	// if err != nil {
	// 	glog.Error("json.Marshal(config.Get()) failed")
	// }
	// text := string(bin)
	// _ = text
	// fmt.Println(text)
	return
}

// ------------------------------------

// GetNPC ...
func GetNPC(id string) *NPC {
	npc, _ := AllNPCs[id]
	return npc
}

// GetTasksByNPC ...
func GetTasksByNPC(id string) (tasks []*Task) {
	for _, task := range AllTasks {
		if task.Npc == id {
			tasks = append(tasks, task)
		}
	}
	return
}

// RandomEvent represents a row from 'db2_utan_meatfloss.tbl_event'.
type RandomEvent struct {
	ID              string `json:"id"`                // id
	Stars           int    `json:"stars"`             // stars
	Type            int    `json:"type"`              // type
	Time            string `json:"time"`              // time
	MinLevel        int    `json:"min_level"`         // min_level
	Intelligence    int    `json:"intelligence"`      // intelligence
	Stamina         int    `json:"stamina"`           // stamina
	FriendlyDegree  int    `json:"friendly_degree"`   // friendly_degree
	DailyTriggerNum int    `json:"daily_trigger_num"` // daily_trigger_num
	TotalTriggerNum int    `json:"total_trigger_num"` // total_trigger_num
	Probability     int    `json:"probability"`       // probability
	Image           string `json:"image"`             // image
	Description     string `json:"description"`       // description
	Choice1         string `json:"choice1"`           // choice1
	Choice2         string `json:"choice2"`           // choice2
	Choice3         string `json:"choice3"`           // choice3
	Reward1         string `json:"reward1"`           // reward1
	Exp1            int    `json:"exp1"`              // exp1
	Reward2         string `json:"reward2"`           // reward2
	Exp2            int    `json:"exp2"`              // exp2
	Reward3         string `json:"reward3"`           // reward3
	Exp3            int    `json:"exp3"`              // exp3
	PreTime         int    `json:"pre_time"`          // pre_time
	PostTime        int    `json:"post_time"`         // post_time

	Choices []string
	Rewards []message.Reward
}

func loadRandomEvents() (err error) {
	dbObjs, err := db.LoadRandomEventConf()
	if err != nil {
		return
	}
	for _, row := range dbObjs {
		obj := &RandomEvent{
			row.ID,
			row.Stars,
			row.Type,
			row.Time,
			row.MinLevel,
			row.Intelligence,
			row.Stamina,
			row.FriendlyDegree,
			row.DailyTriggerNum,
			row.TotalTriggerNum,
			row.Probability,
			row.Image,
			row.Description,
			row.Choice1,
			row.Choice2,
			row.Choice3,
			row.Reward1,
			row.Exp1,
			row.Reward2,
			row.Exp2,
			row.Reward3,
			row.Exp3,
			row.PreTime,
			row.PostTime,
			make([]string, 3),
			make([]message.Reward, 3)}

		rewardStrList := []string{row.Reward1, row.Reward2, row.Reward3}
		expList := []int{row.Exp1, row.Exp2, row.Exp3}
		for i, str := range rewardStrList {
			_ = str
			// for examples, str = wp0001;1000|wp0002;1000
			oneReward := message.Reward{}
			ones := strings.Split(str, "|")
			for _, one := range ones {
				// for example, one =  wp0001;1000
				twos := strings.Split(one, ";")
				if len(twos) == 2 {
					goodsID := strings.TrimSpace(twos[0])
					if goodsID != "" {
						goodsNum, err := strconv.Atoi(twos[1])
						if err == nil {
							sw := message.SingleReward{GoodsID: goodsID, GoodsNum: goodsNum}
							oneReward.List = append(oneReward.List, sw)
						}
					}
				}
			}
			oneReward.Exp = expList[i]
			obj.Rewards[i] = oneReward
		}
		obj.Choices[0] = obj.Choice1
		obj.Choices[1] = obj.Choice2
		obj.Choices[2] = obj.Choice3

		AllRandomEvents[obj.ID] = obj
	}
	bin, err := json.MarshalIndent(AllRandomEvents, "", " ")
	if err != nil {
		glog.Error("json.Marshal(config.Get()) failed")
	}
	text := string(bin)
	_ = text
	// fmt.Println(text)
	return
}

// Apparel ...
type Apparel struct {
	ID                 string `json:"id"`                   // id
	Type               int    `json:"type"`                 // type
	OrderID            int    `json:"order_id"`             // order_id
	ImageName          string `json:"image_name"`           // image_name
	ImageEffect        int    `json:"image_effect"`         // image_effect
	Name               string `json:"name"`                 // name
	Description        string `json:"description"`          // description
	MinLevel           int    `json:"min_level"`            // min_level
	DesignerMinLevel   int    `json:"designer_min_level"`   // designer_min_level
	CanBeSold          int    `json:"can_be_sold"`          // can_be_sold
	PriceForSail       int    `json:"price_for_sail"`       // price_for_sail
	IntelligenceGain   int    `json:"intelligence_gain"`    // intelligence_gain
	StaminaGain        int    `json:"stamina_gain"`         // stamina_gain
	FriendlyDegreeGain int    `json:"friendly_degree_gain"` // friendly_degree_gain
	Stars              int    `json:"stars"`                // stars
	AllowPileup        int    `json:"allow_pileup"`         // allow_pileup

	UniqueID int64
}

func loadApparel() (err error) {
	dbApparel, err := db.LoadApparelConf()
	if err != nil {
		return
	}
	for _, row := range dbApparel {
		apparel := &Apparel{
			ID:                 row.ID,
			Type:               row.Type,
			OrderID:            row.OrderID,
			ImageName:          row.ImageName,
			ImageEffect:        row.ImageEffect,
			Name:               row.Name,
			Description:        row.Description,
			MinLevel:           row.MinLevel,
			DesignerMinLevel:   row.DesignerMinLevel,
			CanBeSold:          row.CanBeSold,
			PriceForSail:       row.PriceForSail,
			IntelligenceGain:   row.IntelligenceGain,
			StaminaGain:        row.StaminaGain,
			FriendlyDegreeGain: row.FriendlyDegreeGain,
			Stars:              row.Stars,
			AllowPileup:        row.AllowPileup}

		temp := strings.Replace(apparel.ID, "fs", "", -1)
		apparel.UniqueID, _ = strconv.ParseInt(temp, 10, 64)
		apparel.UniqueID += 200000
		AllApparels[apparel.ID] = apparel
	}
	utils.PrintJSON(AllApparels)
	return
}

// Furniture ...
type Furniture struct {
	ID               string `json:"id"`                 // id
	Type             int    `json:"type"`               // type
	OrderID          int    `json:"order_id"`           // order_id
	ImageName        string `json:"image_name"`         // image_name
	Icon             string `json:"icon"`               // iconc
	ImageEffect      int    `json:"image_effect"`       // image_effect
	Name             string `json:"name"`               // name
	Description      string `json:"description"`        // description
	MinLevel         int    `json:"min_level"`          // min_level
	DesignerMinLevel int    `json:"designer_min_level"` // designer_min_level
	CanBeSold        int8   `json:"can_be_sold"`        // can_be_sold
	Dismantling      string `json:"dismantling"`        // dismantling
	FashionGain      int    `json:"fashion_gain"`       // fashion_gain
	WarmthGain       int    `json:"warmth_gain"`        // warmth_gain
	CoolGain         int    `json:"cool_gain"`          // cool_gain
	LovelyGain       int    `json:"lovely_gain"`        // lovely_gain
	MotionGain       int    `json:"motion_gain"`        // motion_gain
	Stars            int    `json:"stars"`              // stars
	AllowPileup      int    `json:"allow_pileup"`       // allow_pileup
	MakeTime         int    `json:"maketime"`           // maketime
	MaterialNeed     string `json:"materialneed"`       // materialneed
	NeedMaterial     []message.Guaji
	UniqueID         int64
}

// loadFurniture ...
func loadFurniture() (err error) {
	dbFurniture, err := db.LoadFurnitureConf()
	if err != nil {
		return
	}
	for _, row := range dbFurniture {
		furniture := &Furniture{
			ID:               row.ID,
			Type:             row.Type,
			OrderID:          row.OrderID,
			ImageName:        row.ImageName,
			ImageEffect:      row.ImageEffect,
			Name:             row.Name,
			Description:      row.Description,
			MinLevel:         row.MinLevel,
			DesignerMinLevel: row.DesignerMinLevel,
			CanBeSold:        row.CanBeSold,
			Dismantling:      row.Dismantling,
			FashionGain:      row.FashionGain,
			WarmthGain:       row.WarmthGain,
			CoolGain:         row.CoolGain,
			LovelyGain:       row.LovelyGain,
			MotionGain:       row.MotionGain,
			Stars:            row.Stars,
			Icon:             row.Icon,
			MaterialNeed:     row.MaterialNeed,
			MakeTime:         row.MakeTime,
			AllowPileup:      row.AllowPileup,
			NeedMaterial:     make([]message.Guaji, 1)}

		temp := strings.Replace(furniture.ID, "jj", "", -1)
		furniture.UniqueID, _ = strconv.ParseInt(temp, 10, 64)
		furniture.UniqueID += 300000

		//expList := []string{row.MaterialNeed}
		guajiStrListSingle := []string{row.MaterialNeed}
		for i, str := range guajiStrListSingle {
			_ = str
			// for examples, str = wp0001;1000|wp0002;1000
			oneGuaji := message.Guaji{}
			ones := strings.Split(str, "|")
			for _, one := range ones {
				// for example, one =  wp0001;1000
				twos := strings.Split(one, ";")
				if len(twos) == 2 {
					goodsID := strings.TrimSpace(twos[0])
					if goodsID != "" {
						goodsNum, err := strconv.Atoi(twos[1])
						if err == nil {
							sw := message.SingleGuaji{GoodsID: goodsID, GoodsNum: goodsNum}
							oneGuaji.List = append(oneGuaji.List, sw)
						}
					}
				}
			}
			furniture.NeedMaterial[i] = oneGuaji
		}
		AllFurniture[furniture.ID] = furniture
	}
	utils.PrintJSON(AllFurniture)
	return
}

// NPCGuest ...
type NPCGuest struct {
	ID                string
	AssociationNpc    string
	NpcName           string
	IntimacyLevel     int
	NpcDuration       int
	Dialog1           string
	Dialog2           string
	Dialog3           string
	Reward            string
	MaxRewardTimes    int
	Gift              string
	IntimacyGain      int
	MaxIntimacyDaily  int
	NpcPeriod         string
	AutoProbability   int
	QuestionLibrary   string
	MaxQuestionsDaily int
}

func loadNPCGuest() (err error) {
	dbNPCGuest, err := db.LoadNPCGuestConf()
	if err != nil {
		return
	}
	for _, row := range dbNPCGuest {
		guest := &NPCGuest{
			ID:                row.ID,
			AssociationNpc:    row.AssociationNpc,
			NpcName:           row.NpcName,
			IntimacyLevel:     row.IntimacyLevel,
			NpcDuration:       row.NpcDuration,
			Dialog1:           row.Dialog1,
			Dialog2:           row.Dialog2,
			Dialog3:           row.Dialog3,
			Reward:            row.Reward,
			MaxRewardTimes:    row.MaxRewardTimes,
			Gift:              row.Gift,
			IntimacyGain:      row.IntimacyGain,
			MaxIntimacyDaily:  row.MaxIntimacyDaily,
			NpcPeriod:         row.NpcPeriod,
			AutoProbability:   row.AutoProbability,
			QuestionLibrary:   row.QuestionLibrary,
			MaxQuestionsDaily: row.MaxQuestionsDaily}
		AllNPCGuests[guest.ID] = guest
	}
	utils.PrintJSON(AllNPCGuests)
	return
}

// Guaji ...
type Guaji struct {
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
	Outputs             []message.Reward
	Guajis              []message.Guaji
}

func loadGuaji() (err error) {

	dbGuaji, err := db.LoadGuajiConf()
	if err != nil {
		return
	}
	for _, row := range dbGuaji {
		obj := &Guaji{
			Number:              row.ID,
			MachineLevel:        row.Jqlv,
			MinLevel:            row.Lv,
			Speed:               row.Sudu,
			Quality:             row.Zhiliang,
			Luck:                row.Yunqi,
			InitialTemperature:  row.Cswd,
			MaxTemperature:      row.Zgwd,
			CDPerDegree:         row.Mmlq,
			CD:                  row.Cd,
			TemperaturePerClick: row.Mcdj,
			MachineImage:        row.Tupian,
			NumEmployees:        row.Glnpc,
			PositiveOutput:      row.Zhengxiang,
			Probability1:        row.Gailv1,
			OppositeOutput:      row.Fuxiang,
			Probability2:        row.Gailv2,
			ClickOutput:         row.Djcc,
			CritProbability:     row.Bjgl,
			CritOutput:          row.Bjcc,
			Upmaterial:          row.Sjcl,
			Uptime:              row.Time,
			Outputs:             make([]message.Reward, 3),
			Guajis:              make([]message.Guaji, 1)}
		guajiStrList := []string{row.Zhengxiang, row.Fuxiang, row.Bjcc}
		expList := []int{row.Gailv1, row.Gailv2, row.Bjgl}
		for i, str := range guajiStrList {
			_ = str
			// for examples, str = wp0001;1000|wp0002;1000
			oneReward := message.Reward{}
			ones := strings.Split(str, "|")
			for _, one := range ones {
				// for example, one =  wp0001;1000
				twos := strings.Split(one, ";")
				if len(twos) == 2 {
					goodsID := strings.TrimSpace(twos[0])
					if goodsID != "" {
						goodsNum, err := strconv.Atoi(twos[1])
						if err == nil {
							sw := message.SingleReward{GoodsID: goodsID, GoodsNum: goodsNum}
							oneReward.List = append(oneReward.List, sw)
						}
					}
				}
			}
			oneReward.Exp = expList[i]
			obj.Outputs[i] = oneReward
		}

		guajiStrListSingle := []string{row.Sjcl}
		for i, str := range guajiStrListSingle {
			_ = str
			// for examples, str = wp0001;1000|wp0002;1000
			oneGuaji := message.Guaji{}
			ones := strings.Split(str, "|")
			for _, one := range ones {
				// for example, one =  wp0001;1000
				twos := strings.Split(one, ";")
				if len(twos) == 2 {
					goodsID := strings.TrimSpace(twos[0])
					if goodsID != "" {
						goodsNum, err := strconv.Atoi(twos[1])
						if err == nil {
							sw := message.SingleGuaji{GoodsID: goodsID, GoodsNum: goodsNum}
							oneGuaji.List = append(oneGuaji.List, sw)
						}
					}
				}
			}
			obj.Guajis[i] = oneGuaji
		}
		AllGuajis[obj.MachineLevel] = obj
	}
	utils.PrintJSON(AllGuajis)
	return
}

// Employee ...
type Employee struct {
	Number       string // 编号
	AvatarImage  string // 头像图片
	EmployeeName string // 雇员名字
	Speed        int    // 速度
	Quality      int    // 质量
	Luck         int    // 运气
	Introdution  string // 介绍
}

func loadEmployee() (err error) {

	dbEmployee, err := db.LoadEmployee()
	if err != nil {
		return
	}
	for _, row := range dbEmployee {
		employee := &Employee{
			Number:       row.ID,
			AvatarImage:  row.Tupian,
			EmployeeName: row.Name,
			Speed:        row.Sudu,
			Quality:      row.Zhiliang,
			Luck:         row.Yunqi,
			Introdution:  row.Jieshao}
		AllEmployees[employee.Number] = employee
	}
	utils.PrintJSON(AllEmployees)
	return
}

// Hierarchical ...
type Hierarchical struct {
	Level               int    //等级
	EssentialExperience int    // 所需经验值
	OpenFunction        string // 开启功能
	OpenFurniture       string // 开启制作家具
	OpenClothing        string // 开启制作服饰
	OpendEvent          string // 开启的事件
	Reward              string // 奖励
}

func loadHierarchical() (err error) {
	dbHierarchical, err := db.LoadHierarchical()
	if err != nil {
		return
	}
	for _, row := range dbHierarchical {
		hierarchical := &Hierarchical{
			Level:               row.Lv,
			EssentialExperience: row.Exp,
			OpenFunction:        row.Gongneng,
			OpenFurniture:       row.Jiaju,
			OpenClothing:        row.Fushi,
			OpendEvent:          row.Shjian,
			Reward:              row.Jiangli}
		AllHierarchical[hierarchical.Level] = hierarchical
	}
	utils.PrintJSON(AllHierarchical)
	return
}

// Lattice ...
type Lattice struct {
	ID          string // 格子的编号
	UnlockPrice int    // 解锁售价（钻石）
	OrderID     int    // 索引
}

func loadLattice() (err error) {
	dbLattice, err := db.LoadLattice()
	if err != nil {
		return
	}
	for _, row := range dbLattice {
		lattice := &Lattice{
			ID:          row.ID,
			UnlockPrice: row.Shoujia,
			OrderID:     row.Orderid}
		AllLattice[lattice.OrderID] = lattice
	}
	utils.PrintJSON(AllLattice)
	return
}

// Confige ...
type Confige struct {
	ID        int
	Gujiatime int64
}

func loadConfige() (err error) {
	dbConfige, err := db.LoadConfige()
	if err != nil {
		return
	}
	for _, row := range dbConfige {
		confige := &Confige{
			ID:        row.ID,
			Gujiatime: row.Gujiatime}
		AllConfige[confige.ID] = confige
	}
	utils.PrintJSON(AllConfige)
	return
}
