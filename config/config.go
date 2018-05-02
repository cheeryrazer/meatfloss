package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Note, do not call glog functions here, since it is not initialized.
var (
	// Conf is the unique instance of Config
	Conf        *Configuration
	rtConfMutex sync.RWMutex
)

// Get configuration.
func Get() *Configuration {
	return Conf
}

// Configuration  ...
type Configuration struct {
	RedisServer struct {
		Host string `json:"host"`
		Port int    `json:"port"`
		Db   int    `json:"db"`
	} `json:"redisServer"`
	Log struct {
		LogDir          string `json:"logDir"`
		AlsoLogToStdErr bool   `json:"alsologtostderr"`
	} `json:"log"`
	MySQLServer []struct {
		Host     string `json:"host"`
		Password string `json:"password"`
		Port     int    `json:"port"`
		User     string `json:"user"`
	} `json:"mysqlServer"`
	IsTestEnvioriment bool `json:"IsTestEnvioriment"`
}

// LoadConfig xxx
func LoadConfig() (err error) {
	configPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if !strings.HasSuffix(configPath, string(os.PathSeparator)) {
		configPath += string(os.PathSeparator)
	}
	configPath += "server.json"
	file, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer file.Close()
	conf := &Configuration{}
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(conf); err != nil {
		return
	}
	if len(conf.MySQLServer) < 1 {
		return errors.New("MySQL Server config not found.")
	}
	if len(conf.Log.LogDir) == 0 {
		conf.Log.LogDir = "./log"
	}
	fmt.Printf("+%v", *conf)
	Conf = conf
	return
}
