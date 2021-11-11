package timed_task

type TimeTask struct {
	delMgoChatChan chan bool
}

var timeTask TimeTask

func GetInstance() *TimeTask {
	if timeTask.delMgoChatChan == nil {
		timeTask.delMgoChatChan = make(chan bool)
		go func() {
			timeTask.delMgoChatChan <- true
		}()
	}
	return &timeTask
}

func (t *TimeTask) Run() {
	for {
		select {
		case <-t.delMgoChatChan:
			t.timedDeleteUserChat()
		}
	}
}
