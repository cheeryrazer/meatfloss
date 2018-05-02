package guest

import (
	"assistant_game_server/gameconf"
	"assistant_game_server/utils"
)

// var (
// 	userIDChan chan int
// )

// func init() {
// 	userIDChan = make(chan int, 1024*64)
// }

// // AddUserID ...
// func AddUserID(userID int) {
// 	userIDChan <- userID
// }

// // Start starts guest worker goroutine
// func Start() {
// 	go worker()
// }

// func worker() {
// 	for {
// 		userID := <-userIDChan
// 		dealWithUserID(userID)
// 		fmt.Println("UserID: ", userID)
// 	}
// }

// func dealWithUserID(userID int) {

// }

func calculateNPCList() []string {
	var list []string
	for _, v := range gameconf.AllNPCGuests {
		if utils.ProbabilityHit(v.AutoProbability) {
			list = append(list, v.ID)
		}
	}
	return list
}
