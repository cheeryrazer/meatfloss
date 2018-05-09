package gameredis

import (
	"encoding/json"
	"fmt"
	"meatfloss/common"
	"meatfloss/gameuser"
	"strconv"

	"meatfloss/config"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

var (
	// redisClient ...
	redisClient *redis.Client
)

// GetRunningTask ...
func GetRunningTask(userID int) (taskInfo string, err error) {
	taskInfo, err = redisClient.HGet("runningTask", strconv.Itoa(userID)).Result()
	if err == redis.Nil {
		err = nil
	}
	return
}

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

// PersistUser ...
func PersistUser(userID int, user *gameuser.User) (err error) {
	key := fmt.Sprintf("user:%d", userID)
	fields := make(map[string]interface{})

	if user.Profile != nil {
		data, err := json.Marshal(user.Profile)
		if err == nil {
			fields["profile"] = string(data)
		}
	}

	if user.Bag != nil {
		data, err := json.Marshal(user.Bag)
		if err == nil {
			fields["bag"] = string(data)
		}
	}

	if user.TaskBox != nil {
		data, err := json.Marshal(user.TaskBox)
		if err == nil {
			fields["taskbox"] = string(data)
		}
	}

	if user.NewsBox != nil {
		data, err := json.Marshal(user.NewsBox)
		if err == nil {
			fields["newsbox"] = string(data)
		}
	}

	if user.EventBox != nil {
		data, err := json.Marshal(user.EventBox)
		if err == nil {
			glog.Info(string(data))
			fields["eventbox"] = string(data)
		}
	}

	_, err = redisClient.HMSet(key, fields).Result()
	if err != nil {
		glog.Warning("redisClient.HMSet failed, error: %s", err)
	}

	return
}

// Profile  *Profile
// Bag      *Bag
// TaskBox  *TaskBox
// NewsBox  *NewsBox
// EventBox *EventBox

// SaveBagInfo ..
func SaveBagInfo(userID int, bag *common.Bag) (err error) {
	bagKey := "role:bag"
	userIDStr := strconv.Itoa(userID)
	bagData, _ := json.Marshal(bag)
	bagStr := string(bagData)
	_, err = redisClient.HSet(bagKey, userIDStr, bagStr).Result()
	return
}

// LoadUser ...
func LoadUser(userID int) *gameuser.User {
	key := fmt.Sprintf("user:%d", userID)
	result, err := redisClient.HMGet(key, []string{
		"profile", // 0
		"bag",     // 1
		"taskbox", // 2
		"newsbox", // 3
		"eventbox"}...).Result()
	_ = err
	_ = result
	if err != nil {
		glog.Errorf("LoadUser from redis failed!")
		return nil
	}

	user := &gameuser.User{}
	user.UserID = userID

	// profile
	if result[0] != nil {
		data, ok := result[0].(string)
		if ok && data != "" {
			obj := &gameuser.Profile{}
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.Profile = obj
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}

	// bag
	if result[1] != nil {
		data, ok := result[1].(string)
		if ok && data != "" {
			obj := &common.Bag{}
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.Bag = obj
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}

	// taskbox
	if result[2] != nil {
		data, ok := result[2].(string)
		if ok && data != "" {
			obj := &gameuser.TaskBox{}
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.TaskBox = obj
				taskNum := len(user.TaskBox.Tasks)
				fmt.Println(data)
				fmt.Println("userID: ", obj.UserID)
				fmt.Println("taskNum: ", taskNum)

				_ = taskNum
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}

	// newsbox
	if result[3] != nil {
		data, ok := result[3].(string)
		if ok && data != "" {
			obj := &gameuser.NewsBox{}
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.NewsBox = obj
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}

	// eventbox
	if result[4] != nil {
		data, ok := result[4].(string)
		if ok && data != "" {
			obj := gameuser.NewEventBox(userID)
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.EventBox = obj
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}

	if user.Profile == nil {
		user.Profile = gameuser.NewProfile(userID)
	}

	if user.Bag == nil {
		user.Bag = common.NewBag(userID)
	}

	if user.TaskBox == nil {
		user.TaskBox = gameuser.NewTaskBox(userID)
	}

	if user.NewsBox == nil {
		user.NewsBox = gameuser.NewNewsBox(userID)
	}

	if user.EventBox == nil {
		user.EventBox = gameuser.NewEventBox(userID)
	}

	return user
}

// structures
