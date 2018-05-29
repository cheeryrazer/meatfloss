package gameredis

import (
	"encoding/json"
	"fmt"
	"meatfloss/common"
	"meatfloss/gameuser"
	"meatfloss/message"
	"strconv"

	"meatfloss/config"

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
	
			fields["eventbox"] = string(data)
		}
	}

	if user.Layout != nil {
		data, err := json.Marshal(user.Layout)
		if err == nil {
	
			fields["layout"] = string(data)
		}
	}

	if user.LoginTime != nil {
		data, err := json.Marshal(user.LoginTime)
		if err == nil {
	
			fields["logintime"] = string(data)
		}
	}

	if user.GuajiOutputBox != nil {
		data, err := json.Marshal(user.GuajiOutputBox)
		if err == nil {
		
			fields["guajioutputbox"] = string(data)
		}
	}

	if user.ClickOutputBox != nil {
		data, err := json.Marshal(user.ClickOutputBox)
		if err == nil {
	
			fields["clickoutputbox"] = string(data)
		}
	}

	if user.GuajiSettlement != nil {
		data, err := json.Marshal(user.GuajiSettlement)
		if err == nil {
			fields["guajisettlement"] = string(data)
		}
	}

	if user.GuajiProfile != nil {
		data, err := json.Marshal(user.GuajiProfile)
		if err == nil {
			fields["guajiprofile"] = string(data)
		}
	}

	_, err = redisClient.HMSet(key, fields).Result()
	if err != nil {
		glog.Warning("redisClient.HMSet failed, error: %s", err)
	}

	return
}

// Profile   *Profile
// Bag       *Bag
// TaskBox   *TaskBox
// NewsBox   *NewsBox
// EventBox  *EventBox
// Layout    *Layout
// LoginTime *LoginTime
// GuajiOutputBox  *GuajiOutputBox
// GuajiSettlement  *GuajiSettlement
// GuajiProfile  *GuajiProfile

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
		"profile",         // 0
		"bag",             // 1
		"taskbox",         // 2
		"newsbox",         // 3
		"eventbox",        // 4
		"layout",          // 5
		"logintime",       // 6
		"guajioutputbox",  // 7
		"guajisettlement", // 8
		"guajiprofile",
		"clickoutputbox"}...).Result()
	fmt.Println(result)
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
				// fmt.Println(data)
				// fmt.Println("userID: ", obj.UserID)
				// fmt.Println("taskNum: ", taskNum)

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

	// Layout
	if result[5] != nil {
		data, ok := result[5].(string)
		if ok && data != "" {
			obj := &message.ClientLayout{}
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.Layout = obj
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}
	// LoginTime
	if result[6] != nil {
		data, ok := result[6].(string)
		if ok && data != "" {
			obj := &gameuser.LoginTime{}
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				fmt.Println("测试挂机---------------3333")
				user.LoginTime = obj
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}

	//GuajiOutputBox
	if result[7] != nil {
		data, ok := result[7].(string)
		if ok && data != "" {
			obj := gameuser.NewGuajiOutputBox(userID)
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.GuajiOutputBox = obj
				fmt.Println("测试挂机---------------444")
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}

	//GuajiSettlement
	if result[8] != nil {
		data, ok := result[8].(string)
		if ok && data != "" {
			obj := &gameuser.GuajiSettlement{}
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.GuajiSettlement = obj
				fmt.Println("测试挂机---------------666")
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}

	//GuajiProfile
	if result[9] != nil {
		data, ok := result[9].(string)
		if ok && data != "" {
			obj := &gameuser.GuajiProfile{}
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.GuajiProfile = obj

			} else {
				fmt.Println(err)
				glog.Warning("json.Unmarshal failed")
			}
		}
	}
	//ClickOutputBox
	if result[10] != nil {
		data, ok := result[10].(string)
		if ok && data != "" {
			obj := gameuser.NewClickOutputBox(userID)
			err := json.Unmarshal([]byte(data), obj)
			if err == nil {
				user.ClickOutputBox = obj
				fmt.Println("测试挂机---------------555")
			} else {
				glog.Warning("json.Unmarshal failed")
			}
		}
	}
	fmt.Println(user.GuajiProfile)
	fmt.Println("测试挂机---------------")
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

	if user.Layout == nil {
		user.Layout = message.NewClientLayout()
	}

	if user.LoginTime == nil {
		user.LoginTime = gameuser.NewLoginTime(userID)
	}

	if user.GuajiOutputBox == nil {
		user.GuajiOutputBox = gameuser.NewGuajiOutputBox(userID)
	}

	if user.ClickOutputBox == nil {
		user.ClickOutputBox = gameuser.NewClickOutputBox(userID)
	}

	if user.GuajiSettlement == nil {
		user.GuajiSettlement = gameuser.NewGuajiSettlement(userID)
	}

	if user.GuajiProfile == nil {
		user.GuajiProfile = gameuser.NewGuajiProfile(userID)
	}

	return user
}

// structures
