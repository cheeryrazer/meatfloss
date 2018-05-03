package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // import mysql driver
	"github.com/golang/glog"
)

var (
	// ErrNoRecord no record found
	ErrNoRecord = errors.New("no record found")
)

var db *sql.DB

// Initialize database.
func Initialize(host string, port int, user, password string) error {
	//  "root:helloworld@(127.0.0.1:3306)/utan1"
	dsn := fmt.Sprintf("%v:%v@(%v:%v)/?charset=utf8", user, password, host, port)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err == nil {
		if e := db.Ping(); e != nil {
			glog.Info("db ping failed: ", e)
		} else {
			glog.Info("db ping ok!: ")
		}
	}
	db.SetMaxOpenConns(256)
	return err
}

// GetUserID ...
func GetUserID(phone string) (userID int, err error) {
	sql := fmt.Sprintf("SELECT id from meatfloss.tbl_account where account = '%s'", phone)
	glog.Info(sql)
	rows, err := db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, ErrNoRecord
	}

	err = rows.Scan(&userID)
	if err != nil {
		return
	}
	return
}

// CreateAccount ...
func CreateAccount(account string) (userID int, err error) {
	sql := fmt.Sprintf("INSERT INTO meatfloss.tbl_account(account) values ('%s')", account)
	glog.Info(sql)
	res, err := db.Exec(sql)
	if err != nil {
		glog.Errorf("failed to execute %s", sql)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		glog.Errorf("Error: %s", err.Error())
		return
	}
	println("LastInsertId:", id)
	userID = int(id)
	return

	// sql := fmt.Sprintf("INSERT INTO meatfloss.tbl_roles(role_id, user_id, name, name, type) values ('%s')", account)
	// glog.Info(sql)
	// res, err := db.Exec(sql)
	// if err != nil {
	// 	glog.Errorf("failed to execute %s", sql)
	// 	return
	// }
}

// GetAllUserIDs ...
func GetAllUserIDs() (userIDs []int, err error) {
	sql := fmt.Sprintf("SELECT id from meatfloss.tbl_account")
	glog.Info(sql)
	rows, err := db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		err = rows.Scan(&userID)
		if err != nil {
			return
		}
		userIDs = append(userIDs, userID)
	}

	return
}

// Goods ...
type Goods struct {
	ID                 string `json:"id"`                   // id
	Type               int    `json:"type"`                 // type
	OrderID            int    `json:"order_id"`             // order_id
	Consumable         int    `json:"consumable"`           // consumable
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
	ExperienceGain     int    `json:"experience_gain"`      // experience_gain
	FriendlyDegreeGain int    `json:"friendly_degree_gain"` // friendly_degree_gain
	AllowPileup        int    `json:"allow_pileup"`         // allow_pileup
}

// LoadGoodsConf ...
func LoadGoodsConf() ([]*Goods, error) {
	const sqlstr = `SELECT ` +
		`id, type, order_id, consumable, image_name, image_effect, name, description, min_level, designer_min_level, can_be_sold, price_for_sail, intelligence_gain, stamina_gain, experience_gain, friendly_degree_gain, allow_pileup` +
		` FROM meatfloss.tbl_goods `

	q, err := db.Query(sqlstr)
	if err != nil {
		glog.Errorf("db.Query failed, error: %s", err)
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Goods
	for q.Next() {
		tg := Goods{}

		// scan
		err = q.Scan(&tg.ID, &tg.Type, &tg.OrderID, &tg.Consumable, &tg.ImageName, &tg.ImageEffect, &tg.Name, &tg.Description, &tg.MinLevel, &tg.DesignerMinLevel, &tg.CanBeSold, &tg.PriceForSail, &tg.IntelligenceGain, &tg.StaminaGain, &tg.ExperienceGain, &tg.FriendlyDegreeGain, &tg.AllowPileup)
		if err != nil {
			return nil, err
		}

		res = append(res, &tg)
	}

	return res, nil
}

// NPC ...
type NPC struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Gender           int    `json:"gender"`
	Spine            string `json:"spine"`
	DecorationType   int    `json:"decorationType"`
	GuestDecoration  int    `json:"guestDecoration"`
	Intimacy         int    `json:"intimacy"`
	Intelligence     int    `json:"intelligence"`
	Stamina          int    `json:"stamina"`
	GuestProbability int    `json:"guestProbability"`
}

// LoadNPCConf ...
func LoadNPCConf() (npcList []NPC, err error) {
	sql := fmt.Sprintf("select id, name, description, gender, spine, decoration_type, guest_decoration, intimacy, intelligence, stamina, guest_probability from meatfloss.tbl_npc")
	rows, err := db.Query(sql)
	if err != nil {
		glog.Errorf("failed to execute %s", sql)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var npc = NPC{}
		err = rows.Scan(
			&npc.ID,
			&npc.Name,
			&npc.Description,
			&npc.Gender,
			&npc.Spine,
			&npc.DecorationType,
			&npc.GuestDecoration,
			&npc.Intimacy,
			&npc.Intelligence,
			&npc.Stamina,
			&npc.GuestProbability)
		if err != nil {
			return nil, err
		}
		npcList = append(npcList, npc)
	}
	return
}

// Task ...
type Task struct {
	ID               string `json:"id"`
	Stars            int    `json:"stars"`
	Type             int    `json:"type"`
	MinLevel         int    `json:"min_level"`
	Npc              string `json:"npc"`
	Intelligence     int    `json:"intelligence"`
	Stamina          int    `json:"stamina"`
	FriendlyDegree   int    `json:"friendly_degree"`
	DailyTriggerNum  int    `json:"daily_trigger_num"`
	TotalTriggerNum  int    `json:"total_trigger_num"`
	Probability      int    `json:"probability"`
	AssociationGroup int    `json:"association_group"`
	TriggerOrder     int    `json:"trigger_order"`
	IntimacyNpc      int    `json:"intimacy_npc"`
	IntimacyGain     int    `json:"intimacy_gain"`
	Image            string `json:"image"`
	Description      string `json:"description"`
	Choice1          string `json:"choice1"`
	Choice2          string `json:"choice2"`
	Choice3          string `json:"choice3"`
	Reward1          string `json:"reward1"`
	Exp1             int    `json:"exp1"`
	Reward2          string `json:"reward2"`
	Exp2             int    `json:"exp2"`
	Reward3          string `json:"reward3"`
	Exp3             int    `json:"exp3"`
	PreTime          int    `json:"pre_time"`
	PostTime         int    `json:"post_time"`
}

// LoadTaskConf ...
func LoadTaskConf() (objList []Task, err error) {
	sql := fmt.Sprintf("select id, stars, type, min_level, npc, intelligence, stamina, friendly_degree, daily_trigger_num, total_trigger_num, probability, association_group, trigger_order, intimacy_npc, intimacy_gain, image, description, choice1, choice2, choice3, reward1, exp1, reward2, exp2, reward3, exp3, pre_time, post_time from meatfloss.tbl_task")
	rows, err := db.Query(sql)
	if err != nil {
		glog.Errorf("failed to execute %s", sql)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var obj = Task{}
		err = rows.Scan(
			&obj.ID,
			&obj.Stars,
			&obj.Type,
			&obj.MinLevel,
			&obj.Npc,
			&obj.Intelligence,
			&obj.Stamina,
			&obj.FriendlyDegree,
			&obj.DailyTriggerNum,
			&obj.TotalTriggerNum,
			&obj.Probability,
			&obj.AssociationGroup,
			&obj.TriggerOrder,
			&obj.IntimacyNpc,
			&obj.IntimacyGain,
			&obj.Image,
			&obj.Description,
			&obj.Choice1,
			&obj.Choice2,
			&obj.Choice3,
			&obj.Reward1,
			&obj.Exp1,
			&obj.Reward2,
			&obj.Exp2,
			&obj.Reward3,
			&obj.Exp3,
			&obj.PreTime,
			&obj.PostTime)
		if err != nil {
			return nil, err
		}
		objList = append(objList, obj)
	}
	return
}

// RandomEvent represents a row from 'meatfloss.tbl_event'.
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
}

// LoadRandomEventConf ...
func LoadRandomEventConf() ([]*RandomEvent, error) {
	const sqlstr = `SELECT ` +
		`id, stars, type, time, min_level, intelligence, stamina, friendly_degree, daily_trigger_num, total_trigger_num, probability, image, description, choice1, choice2, choice3, reward1, exp1, reward2, exp2, reward3, exp3, pre_time, post_time` +
		` FROM meatfloss.tbl_event `

	q, err := db.Query(sqlstr)
	if err != nil {
		fmt.Printf("sql: %s\n", sqlstr)

		fmt.Printf("err: %s\n", err)
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*RandomEvent
	for q.Next() {
		te := RandomEvent{}

		// scan
		err = q.Scan(&te.ID, &te.Stars, &te.Type, &te.Time, &te.MinLevel, &te.Intelligence, &te.Stamina, &te.FriendlyDegree, &te.DailyTriggerNum, &te.TotalTriggerNum, &te.Probability, &te.Image, &te.Description, &te.Choice1, &te.Choice2, &te.Choice3, &te.Reward1, &te.Exp1, &te.Reward2, &te.Exp2, &te.Reward3, &te.Exp3, &te.PreTime, &te.PostTime)
		if err != nil {
			return nil, err
		}

		res = append(res, &te)
	}

	return res, nil
}

// NPCGuest ...
type NPCGuest struct {
	ID                string `json:"id"`                  // id
	AssociationNpc    string `json:"association_npc"`     // association_npc
	NpcName           string `json:"npc_name"`            // npc_name
	IntimacyLevel     int    `json:"intimacy_level"`      // intimacy_level
	NpcDuration       int    `json:"npc_duration"`        // npc_duration
	Dialog1           string `json:"dialog1"`             // dialog1
	Dialog2           string `json:"dialog2"`             // dialog2
	Dialog3           string `json:"dialog3"`             // dialog3
	Reward            string `json:"reward"`              // reward
	MaxRewardTimes    int    `json:"max_reward_times"`    // max_reward_times
	Gift              string `json:"gift"`                // gift
	IntimacyGain      int    `json:"intimacy_gain"`       // intimacy_gain
	MaxIntimacyDaily  int    `json:"max_intimacy_daily"`  // max_intimacy_daily
	NpcPeriod         string `json:"npc_period"`          // npc_period
	AutoProbability   int    `json:"auto_probability"`    // auto_probability
	QuestionLibrary   string `json:"question_library"`    // question_library
	MaxQuestionsDaily int    `json:"max_questions_daily"` // max_questions_daily
}

// LoadNPCGuestConf ...
func LoadNPCGuestConf() ([]*NPCGuest, error) {
	const sqlstr = `SELECT ` +
		`id, association_npc, npc_name, intimacy_level, npc_duration, dialog1, dialog2, dialog3, reward, max_reward_times, gift, intimacy_gain, max_intimacy_daily, npc_period, auto_probability, question_library, max_questions_daily` +
		` FROM meatfloss.tbl_guest `

	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*NPCGuest
	for q.Next() {
		tg := NPCGuest{}

		// scan
		err = q.Scan(&tg.ID, &tg.AssociationNpc, &tg.NpcName, &tg.IntimacyLevel, &tg.NpcDuration, &tg.Dialog1, &tg.Dialog2, &tg.Dialog3, &tg.Reward, &tg.MaxRewardTimes, &tg.Gift, &tg.IntimacyGain, &tg.MaxIntimacyDaily, &tg.NpcPeriod, &tg.AutoProbability, &tg.QuestionLibrary, &tg.MaxQuestionsDaily)
		if err != nil {
			return nil, err
		}

		res = append(res, &tg)
	}

	return res, nil
}

// Apparel represents a row from 'meatfloss.tbl_apparel'.
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
}

// LoadApparelConf ...
func LoadApparelConf() ([]*Apparel, error) {
	const sqlstr = `SELECT ` +
		`id, type, order_id, image_name, image_effect, name, description, min_level, designer_min_level, can_be_sold, price_for_sail, intelligence_gain, stamina_gain, friendly_degree_gain, stars, allow_pileup` +
		` FROM meatfloss.tbl_apparel `

	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Apparel
	for q.Next() {
		ta := Apparel{}

		// scan
		err = q.Scan(&ta.ID, &ta.Type, &ta.OrderID, &ta.ImageName, &ta.ImageEffect, &ta.Name, &ta.Description, &ta.MinLevel, &ta.DesignerMinLevel, &ta.CanBeSold, &ta.PriceForSail, &ta.IntelligenceGain, &ta.StaminaGain, &ta.FriendlyDegreeGain, &ta.Stars, &ta.AllowPileup)
		if err != nil {
			return nil, err
		}

		res = append(res, &ta)
	}

	return res, nil
}

// Furniture represents a row from 'meatfloss.tbl_furniture'.
type Furniture struct {
	ID               string `json:"id"`                 // id
	Type             int    `json:"type"`               // type
	OrderID          int    `json:"order_id"`           // order_id
	ImageName        string `json:"image_name"`         // image_name
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
}

// LoadFurnitureConf ...
func LoadFurnitureConf() ([]*Furniture, error) {
	const sqlstr = `SELECT ` +
		`id, type, order_id, image_name, image_effect, name, description, min_level, designer_min_level, can_be_sold, dismantling, fashion_gain, warmth_gain, cool_gain, lovely_gain, motion_gain, stars, allow_pileup` +
		` FROM meatfloss.tbl_furniture `

	q, err := db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	var res []*Furniture
	for q.Next() {
		tf := Furniture{}

		// scan
		err = q.Scan(&tf.ID, &tf.Type, &tf.OrderID, &tf.ImageName, &tf.ImageEffect, &tf.Name, &tf.Description, &tf.MinLevel, &tf.DesignerMinLevel, &tf.CanBeSold, &tf.Dismantling, &tf.FashionGain, &tf.WarmthGain, &tf.CoolGain, &tf.LovelyGain, &tf.MotionGain, &tf.Stars, &tf.AllowPileup)
		if err != nil {
			return nil, err
		}

		res = append(res, &tf)
	}

	return res, nil
}
