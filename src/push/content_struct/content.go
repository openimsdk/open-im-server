/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/27 11:24).
 */
package content_struct

import "encoding/json"

type Content struct {
	IsDisplay int32  `json:"isDisplay"`
	ID        string `json:"id"`
	Text      string `json:"text"`
}

func NewContentStructString(isDisplay int32, ID string, text string) string {
	c := Content{IsDisplay: isDisplay, ID: ID, Text: text}
	return c.contentToString()
}
func (c *Content) contentToString() string {
	data, _ := json.Marshal(c)
	dataString := string(data)
	return dataString
}
