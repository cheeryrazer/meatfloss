package gameredis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"assistant_game_server/common"
	"assistant_game_server/config"
	"assistant_game_server/db"
	"assistant_game_server/message"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

var (
	// redisClient ...
	redisClient *redis.Client
)

// Initialize redis.
func Initialize() {
	addr := fmt.Sprint(config.Get().RedisServer.Host, ":", config.Get().RedisServer.Port)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",                          // no password set
		DB:       config.Get().RedisServer.Db, // use default DB
		PoolSize: 64,                          // max connections
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		glog.Error("redisClient.Ping() failed, error: ", err)
	} else {
		glog.Info("redisClient.Ping() ok!")
	}

	return
}

// MarkNewsAdRead ...
func MarkNewsAdRead(userID int, pushID string, articleID string) (err error) {
	key := fmt.Sprintf("newspush:%d", userID)
	value := fmt.Sprintf("%s:%s", pushID, articleID)
	result := redisClient.LRem(key, 1, value)
	err = result.Err()
	return
}

// GetRunningTask ...
func GetRunningTask(userID int) (taskInfo string, err error) {
	taskInfo, err = redisClient.HGet("runningTask", strconv.Itoa(userID)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

// GetAllRunningTask ...
func GetAllRunningTask() (tasks map[int]*RunningTaskInfo) {
	result, err := redisClient.HGetAll("runningTask").Result()
	if err != nil {
		glog.Errorf("RedisClient.HGetAll failed, error: %s", err)
		return
	}

	tasks = make(map[int]*RunningTaskInfo)
	for k, v := range result {
		userID, err := strconv.Atoi(k)
		if err != nil {
			glog.Infof("strconv.Atoi failed, error: %s", err)
			continue
		}
		taskInfo := &RunningTaskInfo{}
		err = json.Unmarshal([]byte(v), taskInfo)
		if err != nil {
			glog.Infof("strconv.Atoi failed, error: %s", err)
			continue
		}
		tasks[userID] = taskInfo
	}
	return
}

// SetRunningTask ...
func SetRunningTask(userID int, taskInfo *RunningTaskInfo) (result bool, err error) {
	infoStr, _ := json.Marshal(taskInfo)
	result, err = redisClient.HSet("runningTask", strconv.Itoa(userID), infoStr).Result()
	if err != nil {
		glog.Infof("RedisClient.HSet failed, error: %s", err)
	}
	return
}

// GetAllNews ..
func GetAllNews(userID int) (newArticleInfos []message.ArticleInfo) {
	newArticleInfos = make([]message.ArticleInfo, 0)
	key := fmt.Sprintf("newspush:%d", userID)
	result, err := redisClient.LRange(key, 0, 1000).Result()
	if err != nil {
		return
	}

	var articleIDs []string
	var pushIDs []string
	for _, val := range result {
		ids := strings.Split(val, ":")
		if len(ids) != 2 {
			continue
		}
		pushID := ids[0]
		if len(pushID) == 0 {
			continue
		}

		articleID := ids[1]
		if len(articleID) == 0 {
			continue
		}
		pushIDs = append(pushIDs, pushID)
		articleIDs = append(articleIDs, articleID)
	}

	articleResult, err := redisClient.HMGet("articleinfo", articleIDs...).Result()
	if err != nil {
		return
	}

	articleInfos := make(map[string]*message.ArticleInfo)
	for _, val := range articleResult {
		_ = val
		b64Text, ok := val.(string)
		if !ok {
			continue
		}
		jsonText, e := base64.StdEncoding.DecodeString(b64Text)
		if e != nil {
			continue
		}
		articleInfo := &message.ArticleInfo{}
		e = json.Unmarshal(jsonText, articleInfo)
		if e != nil {
			continue
		}
		articleInfos[articleInfo.ArticleID] = articleInfo
	}

	for i, pushID := range pushIDs {
		articleID := articleIDs[i]
		articleInfo, ok := articleInfos[articleID]
		if !ok {
			continue
		}

		newArticleInfo := message.ArticleInfo{}
		newArticleInfo.ArticleID = articleInfo.ArticleID
		newArticleInfo.PushID = pushID
		newArticleInfo.PicURL = articleInfo.PicURL
		newArticleInfo.Tags = articleInfo.Tags
		newArticleInfo.Title = articleInfo.Title
		newArticleInfos = append(newArticleInfos, newArticleInfo)
	}
	return
}

// ClearTasksAndSaveEvents ...
func ClearTasksAndSaveEvents(tasks []*RunningTaskInfo, events []*message.EventInfo) {
	pipe := redisClient.Pipeline()
	for _, task := range tasks {
		pipe.HDel("runningTask", strconv.Itoa(task.UserID))
	}

	for _, event := range events {
		str, _ := json.Marshal(event)
		key := fmt.Sprintf("events:%d", event.UserID)
		pipe.HSet(key, event.GenID, str)
	}

	_, err := pipe.Exec()
	if err != nil {
		glog.Errorf("pipe.Exec() failed")
	}
}

// InitRole ...
func InitRole(userID int) (err error) {
	pipe := redisClient.Pipeline()
	key := fmt.Sprintf("role:profile:%d", userID)
	fields := make(map[string]interface{})
	fields["name"] = "肉松"
	fields["gender"] = 1
	fields["spine"] = "1"
	fields["level"] = 1
	fields["exp"] = 0
	fields["intimacy"] = 100
	fields["intells"] = 200
	fields["stamina"] = 60
	_, err = pipe.HMSet(key, fields).Result()
	bagKey := "role:bag"
	userIDStr := strconv.Itoa(userID)

	// Goods     100000
	// Apparel   200000
	// Furniture 300000

	bag := common.NewBag()
	{
		cell := &common.BagCell{}
		cell.Count = 5
		cell.GoodsID = "wp0001"
		cell.UniqueID = 100000 + 1
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &common.BagCell{}
		cell.Count = 10
		cell.GoodsID = "wp0002"
		cell.UniqueID = 100000 + 2
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &common.BagCell{}
		cell.Count = 3
		cell.GoodsID = "wp0003"
		cell.UniqueID = 100000 + 3
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &common.BagCell{}
		cell.Count = 1
		cell.GoodsID = "wp0004"
		cell.UniqueID = 100000 + 4
		bag.Cells[cell.UniqueID] = cell
	}

	// {
	// 	cell := &common.BagCell{}
	// 	cell.Count = 1
	// 	cell.GoodsID = "wp0005"
	// 	cell.UniqueID = 100000 + 5
	// 	bag.Cells[cell.UniqueID] = cell
	// }

	// {
	// 	cell := &common.BagCell{}
	// 	cell.Count = 1
	// 	cell.GoodsID = "wp0006"
	// 	cell.UniqueID = 100000 + 6
	// 	bag.Cells[cell.UniqueID] = cell
	// }

	// {
	// 	cell := &common.BagCell{}
	// 	cell.Count = 1
	// 	cell.GoodsID = "wp0007"
	// 	cell.UniqueID = 100000 + 7
	// 	bag.Cells[cell.UniqueID] = cell
	// }

	// {
	// 	cell := &common.BagCell{}
	// 	cell.Count = 1
	// 	cell.GoodsID = "wp0008"
	// 	cell.UniqueID = 100000 + 8
	// 	bag.Cells[cell.UniqueID] = cell
	// }

	// {
	// 	cell := &common.BagCell{}
	// 	cell.Count = 1
	// 	cell.GoodsID = "wp0009"
	// 	cell.UniqueID = 100000 + 9
	// 	bag.Cells[cell.UniqueID] = cell
	// }

	{
		cell := &common.BagCell{}
		cell.Count = 5
		cell.GoodsID = "fs0001"
		cell.UniqueID = 200000 + 1
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &common.BagCell{}
		cell.Count = 10
		cell.GoodsID = "fs0002"
		cell.UniqueID = 200000 + 2
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &common.BagCell{}
		cell.Count = 3
		cell.GoodsID = "fs0003"
		cell.UniqueID = 200000 + 3
		bag.Cells[cell.UniqueID] = cell
	}

	{
		cell := &common.BagCell{}
		cell.Count = 1
		cell.GoodsID = "fs0004"
		cell.UniqueID = 200000 + 4
		bag.Cells[cell.UniqueID] = cell
	}

	bagData, _ := json.Marshal(bag)
	bagStr := string(bagData)
	_, err = pipe.HSet(bagKey, userIDStr, bagStr).Result()
	_, err = pipe.Exec()
	if err == redis.Nil {
		err = nil
	}
	return
}

// SaveBagInfo ..
func SaveBagInfo(userID int, bag *common.Bag) (err error) {
	bagKey := "role:bag"
	userIDStr := strconv.Itoa(userID)
	bagData, _ := json.Marshal(bag)
	bagStr := string(bagData)
	_, err = redisClient.HSet(bagKey, userIDStr, bagStr).Result()
	return
}

// GetRoleInfo ...
func GetRoleInfo(userID int) (info message.GameBaseInfoNotify, Events []*message.EventInfo, bag *common.Bag, npcGuests []string, err error) {
	// TODO, 改成pipe方式
	Events = make([]*message.EventInfo, 0)
	profileKey := fmt.Sprintf("role:profile:%d", userID)
	profileResult, err := redisClient.HGetAll(profileKey).Result()
	if err != nil {
		return
	}

	{
		info.Data.Profile.Name = profileResult["name"]
		gender := profileResult["gender"]
		info.Data.Profile.Gender, err = strconv.Atoi(gender)
		info.Data.Profile.Spine = profileResult["spine"]
		level := profileResult["level"]
		info.Data.Profile.Level, err = strconv.Atoi(level)
		exp := profileResult["exp"]
		info.Data.Profile.Experience, err = strconv.Atoi(exp)
		intimacy := profileResult["intimacy"]
		info.Data.Profile.Intimacy, err = strconv.Atoi(intimacy)
		intells := profileResult["intells"]
		info.Data.Profile.Intelligence, err = strconv.Atoi(intells)
		stamina := profileResult["stamina"]
		info.Data.Profile.Stamina, err = strconv.Atoi(stamina)
	}

	bagKey := "role:bag"
	userIDStr := strconv.Itoa(userID)
	bagResult, err := redisClient.HGet(bagKey, userIDStr).Result()
	if err != nil {
		return
	}

	{
		bag = common.NewBag()
		err = json.Unmarshal([]byte(bagResult), &bag)
		if err != nil {
			return
		}

		for k, v := range bag.Cells {
			cell := message.CellInfo{}
			cell.GoodsID = v.GoodsID
			cell.Count = v.Count
			cell.UniqueID = k
			info.Data.Bag.Cells = append(info.Data.Bag.Cells, cell)
		}

	}

	eventsKey := fmt.Sprintf("events:%d", userID)
	eventsResult, err := redisClient.HGetAll(eventsKey).Result()
	if err != nil {
		glog.Errorf("RedisClient.HGetAll failed, error: %s", err)
		return
	}
	{
		for k, v := range eventsResult {
			_, err := strconv.Atoi(k)
			if err != nil {
				glog.Infof("strconv.Atoi failed, error: %s", err)
				continue
			}
			eventInfo := &message.EventInfo{}
			err = json.Unmarshal([]byte(v), eventInfo)
			if err != nil {
				glog.Infof("json.Unmarshal failed, error: %s", err)
				continue
			}
			Events = append(Events, eventInfo)
		}
	}

	{

		info.Data.Tasks = make([]message.RunningTask, 0)
		taskInfo, e := GetRunningTask(userID)
		if e != nil {
			return
		}

		if taskInfo != "" {
			task := RunningTaskInfo{}
			err := json.Unmarshal([]byte(taskInfo), &task)
			if err == nil {
				newTask := message.RunningTask{}
				newTask.TaskID = task.TaskID
				newTask.PreTime = task.PreTime
				newTask.CreateAt = time.Now().Format("2006-01-02 15:04:05")
				info.Data.Tasks = append(info.Data.Tasks, newTask)
			} else {
				glog.Warningf("json unmarshal failed, err: %s", err)
			}
		}
	}

	{
		dateStr := time.Now().Format("2006-01-02")
		asGuestKey := fmt.Sprintf("asGuest:%d", userID)
		asGuestResult, e := redisClient.HGet(asGuestKey, dateStr).Result()
		if e != nil && e != redis.Nil {
			glog.Errorf("RedisClient.HGetAll failed, error: %s", err)
			return
		}

		var npcGuests []string
		if e != redis.Nil {
			e := json.Unmarshal([]byte(asGuestResult), &npcGuests)
			if e != nil {
				glog.Errorf("json.Unmarshal failed, error: %s", e)
			}
		}
	}

	return
}

// SetNPCGuestList ...
func SetNPCGuestList(userID int, npcGuests []string) (err error) {
	npcGuestsText, _ := json.Marshal(npcGuests)
	dateStr := time.Now().Format("2006-01-02")
	asGuestKey := fmt.Sprintf("asGuest:%d", userID)
	_, err = redisClient.HSet(asGuestKey, dateStr, npcGuestsText).Result()
	if err != nil {
		glog.Errorf("json.Unmarshal failed, error: %s", err)
		return
	}
	return
}

// GetEvent ...
func GetEvent(userID int, eventGenID string) *message.EventInfo {
	key := fmt.Sprintf("events:%d", userID)
	result, err := redisClient.HGet(key, eventGenID).Result()
	if err != nil {
		glog.Errorf("redisClient.HGet failed, error: %s", err)
		return nil
	}

	eventInfo := &message.EventInfo{}
	err = json.Unmarshal([]byte(result), eventInfo)
	if err != nil {
		glog.Infof("json.Unmarshal failed, error: %s", err)
		return nil
	}
	return eventInfo
}

// DelEvent ...
func DelEvent(userID int, eventGenID string) (err error) {
	key := fmt.Sprintf("events:%d", userID)
	_, err = redisClient.HDel(key, eventGenID).Result()
	if err == redis.Nil {
		return nil
	}
	return err
}

// PushArticle ...
func PushArticle(articleInfo *message.ArticleInfo) {
	var err error
	articleID := strings.TrimSpace(articleInfo.ArticleID)
	id := time.Now().UnixNano()
	_ = id
	var userList []int
	userList, err = db.GetAllUserIDs()
	pipe := redisClient.Pipeline()
	value := fmt.Sprintf("%d:%s", id, articleID)
	for _, userID := range userList {
		key := fmt.Sprintf("newspush:%d", userID)
		pipe.RPush(key, value)
		fmt.Println(userID)
	}

	_, err = pipe.Exec()
	if err != nil {
		return
	}

	jsonData, err := json.Marshal(*articleInfo)
	b64Text := base64.StdEncoding.EncodeToString(jsonData)
	redisClient.HSet("articleinfo", articleID, b64Text)
	fmt.Println(b64Text)
}

// GetUniqueID ...
func GetUniqueID() int64 {
	result, err := redisClient.Incr("meatFlossUniqueID").Result()
	if err != nil {
		return 0
	}
	return result
}

// GetGoodsUniqueID ...
func GetGoodsUniqueID() int64 {
	result, err := redisClient.Incr("meatFlossUniqueID").Result()
	if err != nil {
		return 0
	}
	return result + 10000000
}

// structures

// RunningTaskInfo ...
type RunningTaskInfo struct {
	TaskID    string //
	Timestamp int
	PreTime   int
	UserID    int
	ID        int64
	NPCID     string
}
