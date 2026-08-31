package api

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"

	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func (m *MessageApi) GetStreamMsg(c *gin.Context) {
	a2r.Call(c, msg.MsgClient.GetStreamMsg, m.Client)
}

func (m *MessageApi) AppendStreamMsg(c *gin.Context) {
	a2r.Call(c, msg.MsgClient.AppendStreamMsg, m.Client)
}

func (m *MessageApi) PutStreamMsg(c *gin.Context) {
	var (
		conversationID string
		clientMsgID    string
	)
	{
		operationID := c.GetHeader(constant.OperationID)
		if operationID == "" {
			operationID = c.Query(constant.OperationID)
		}
		if operationID == "" {
			m.putErr(c, errs.ErrArgs.WrapMsg("operationID is empty"))
			return
		}
		c.Set(constant.OperationID, operationID)
		conversationID = c.Query("conversationID")
		if conversationID == "" {
			conversationID = c.GetHeader("conversationID")
		}
		if conversationID == "" {
			m.putErr(c, errs.ErrArgs.WrapMsg("conversationID is empty"))
			return
		}
		clientMsgID = c.Query("clientMsgID")
		if clientMsgID == "" {
			clientMsgID = c.GetHeader("clientMsgID")
		}
		if clientMsgID == "" {
			m.putErr(c, errs.ErrArgs.WrapMsg("clientMsgID is empty"))
			return
		}
		token := c.GetHeader("token")
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			m.putErr(c, errs.ErrTokenInvalid.WrapMsg("token is empty"))
			return
		}
		resp, err := m.authClient.ParseToken(c, token)
		if err != nil {
			m.putErr(c, err)
			return
		}
		c.Set(constant.OpUserPlatform, constant.PlatformIDToName(int(resp.PlatformID)))
		c.Set(constant.OpUserID, resp.UserID)
	}
	done := make(chan struct{})
	streamCh := make(chan string, 8)

	go func() {
		defer func() {
			close(streamCh)
			c.Request.Body.Close()
		}()
		buf := make([]byte, 256)
		body := NewUTF8Reader(c.Request.Body)
		for i := 1; ; i++ {
			n, err := body.Read(buf)
			if n > 0 {
				select {
				case streamCh <- string(buf[:n]):
				case <-done:
					return
				}
			}
			if err != nil {
				if err == io.EOF {
					log.ZDebug(c, "read request body stream msg done", "clientMsgID", clientMsgID)
				} else {
					log.ZError(c, "read request body stream msg failed", err, "clientMsgID", clientMsgID, "error", err)
				}
				return
			}
			if n < 10 {
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()

	var (
		packet   []string
		end      bool
		index    int
		errCount int
		lastErr  error
	)
	defer func() {
		close(done)
		if lastErr == nil {
			apiresp.GinSuccess(c, nil)
		} else {
			m.putErr(c, lastErr)
		}
	}()
	doAppend := func() {
		if end == false && len(packet) == 0 {
			return
		}
		ctx, cancel := context.WithTimeout(c, time.Second*10)
		defer cancel()
		req := &msg.AppendStreamMsgReq{
			ConversationID: conversationID,
			ClientMsgID:    clientMsgID,
			StartIndex:     int64(index),
			Packets:        packet,
			End:            end,
		}
		_, lastErr = m.Client.AppendStreamMsg(ctx, req)
		if lastErr == nil {
			log.ZDebug(ctx, "AppendStreamMsg ok", "clientMsgID", clientMsgID)
			index += len(packet)
			packet = packet[:0]
			errCount = 0
			return
		}
		errCount++
		if errs.ErrRecordNotFound.Is(lastErr) {
			log.ZWarn(c, "msg not found", nil, "clientMsgID", clientMsgID)
			return
		} else if errs.ErrNoPermission.Is(lastErr) {
			log.ZError(c, "msg permission error", nil, "clientMsgID", clientMsgID)
			return
		} else {
			log.ZError(c, "append stream msg failed", lastErr, "clientMsgID", clientMsgID, "errCount", errCount)
			time.Sleep(time.Millisecond * 50 * time.Duration(errCount))
		}
	}
	for errCount < 10 {
		select {
		case s, ok := <-streamCh:
			if ok {
				packet = append(packet, s)
			}
			if !ok {
				end = true
			}
			doAppend()
			if end == true && lastErr == nil {
				return
			}
		}
	}
}

func NewUTF8Reader(r io.Reader) io.Reader {
	return &UTF8Reader{
		r: bufio.NewReaderSize(r, 512),
	}
}

type UTF8Reader struct {
	r   *bufio.Reader
	buf bytes.Buffer
}

func (r *UTF8Reader) Read(b []byte) (int, error) {
	for {
		n, err := r.r.Read(b)
		if err != nil {
			return 0, err
		}
		r.buf.Write(b[:n])
		data := r.buf.Bytes()
		minIndex := min(len(b), len(data))
		if minIndex == 0 {
			continue
		}
		for i := minIndex; i > 0; i-- {
			if utf8.Valid(data[:i]) {
				n, err := r.buf.Read(b[:i])
				if err != nil {
					return 0, err
				}
				if n != i {
					return 0, fmt.Errorf("invalid UTF-8 encoding")
				}
				return n, nil
			}
		}
	}
}

func (m *MessageApi) putErr(c *gin.Context, err error) {
	c.JSON(http.StatusOK, apiresp.ParseError(err))
}
