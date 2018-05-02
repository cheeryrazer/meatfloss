package task

import (
	"meatfloss/client"
	"meatfloss/gameredis"
	"meatfloss/message"
	"strconv"
	"time"
)

// StartTaskManager ...
func StartTaskManager() {
	go taskCleaner()
}

func taskCleaner() {
	for {
		cleanOnce()
		time.Sleep(1 * time.Second)
	}

}

func cleanOnce() {
	now := int(time.Now().Unix())
	var finishedTasks []*gameredis.RunningTaskInfo
	tasks := gameredis.GetAllRunningTask()
	for userID, taskInfo := range tasks {
		if now <= taskInfo.Timestamp+taskInfo.PreTime {
			continue
		}
		taskInfo.UserID = userID
		finishedTasks = append(finishedTasks, taskInfo)
	}

	if len(finishedTasks) == 0 {
		return
	}

	// 根据任务触发事件， 然后清除任务
	// 普通时间
	var msgs []message.EventNotify
	var events []*message.EventInfo
	for _, task := range finishedTasks {
		_ = task
		var msg = message.EventNotify{}
		oneEvent := &message.EventInfo{}
		msg.Meta.MessageType = "EventNotify"
		msg.Meta.MessageTypeID = message.MsgTypeEventNotify
		msg.Data.UserID = task.UserID
		oneEvent.Type = "normal"
		oneEvent.Title = "任务完成"
		oneEvent.Content = "任务已经完成了"
		oneEvent.Time = time.Now().Format("2006-01-02 15:04:05")
		oneEvent.UserID = task.UserID
		oneEvent.EventID = task.TaskID
		taskID := gameredis.GetUniqueID()
		oneEvent.GenID = strconv.FormatInt(taskID, 10)
		msg.Data.Events = append(msg.Data.Events, oneEvent)
		msgs = append(msgs, msg)
		events = append(events, oneEvent)
	}

	gameredis.ClearTasksAndSaveEvents(finishedTasks, events)

	for _, msg := range msgs {
		client.Mgr.SendToClient(msg.Data.UserID, msg)
	}

}
