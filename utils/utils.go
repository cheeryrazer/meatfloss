package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
)

var (
	rwMutex  sync.RWMutex
	uniqueID uint64
)

// GetUniqueID ...
func GetUniqueID() uint64 {
	return atomic.AddUint64(&uniqueID, 1)
}

// ProbabilityHit ...
// 0 - 99
func ProbabilityHit(probability int) bool {
	if probability <= 0 {
		return false
	}

	if probability >= 100 {
		return true
	}

	probability--
	// random value in [0, 99]
	random := rand.Intn(99)
	if random <= probability {
		return true
	}
	return false
}

// PrintJSON ...
func PrintJSON(obj interface{}) {
	data, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		fmt.Printf("json.MarshalIndent failed, error: %s", err)
		return
	}
	text := string(data)
	fmt.Println(text)
}
