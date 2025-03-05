package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/protocol/constant"
)

func main() {
	g := gin.Default()
	g.POST("/callbackExample/callbackBeforeMsgModifyCommand", toGin(handlerMsg))
	if err := g.Run(":10006"); err != nil {
		panic(err)
	}
}

func toGin[R any](fn func(c *gin.Context, req *R)) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Printf("HTTP %s %s %s\n", c.Request.Method, c.Request.URL, body)
		var req R
		if err := json.Unmarshal(body, &req); err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		fn(c, &req)
	}
}

func handlerMsg(c *gin.Context, req *cbapi.CallbackMsgModifyCommandReq) {
	var resp cbapi.CallbackMsgModifyCommandResp
	if req.ContentType != constant.Text {
		c.JSON(http.StatusOK, &resp)
		return
	}
	var textElem struct {
		Content string `json:"content"`
	}
	if err := json.Unmarshal([]byte(req.Content), &textElem); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	const word = "xxx"
	if strings.Contains(textElem.Content, word) {
		textElem.Content = strings.ReplaceAll(textElem.Content, word, strings.Repeat("*", len(word)))
		content, err := json.Marshal(&textElem)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		tmp := string(content)
		resp.Content = &tmp
	}
	c.JSON(http.StatusOK, &resp)
}
