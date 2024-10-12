package push

import (
	"github.com/openimsdk/protocol/sdkws"
	"testing"
)

func TestName(t *testing.T) {
	var c ConsumerHandler
	c.readCh = make(chan *sdkws.MarkAsReadTips)

	go c.loopRead()

	go func() {
		for i := 0; ; i++ {
			seq := int64(i + 1)
			if seq%3 == 0 {
				seq = 1
			}
			c.readCh <- &sdkws.MarkAsReadTips{
				ConversationID:   "c100",
				MarkAsReadUserID: "u100",
				HasReadSeq:       seq,
			}
		}
	}()

	select {}
}
