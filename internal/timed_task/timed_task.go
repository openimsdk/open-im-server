package timed_task

import (
	"Open_IM/pkg/common/db"
	"time"
)

func (t *TimeTask) timedDeleteUserChat() {
	now := time.Now()
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	tm := time.NewTimer(next.Sub(now))

	<-tm.C

	count, _ := db.DB.MgoUserCount()
	for i := 0; i < count; i++ {
		time.Sleep(10 * time.Millisecond)
		uid, _ := db.DB.MgoSkipUID(i)
		db.DB.DelUserChatMongo2(uid)
	}

	go func() {
		t.delMgoChatChan <- true
	}()
}
