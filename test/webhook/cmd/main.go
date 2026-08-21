// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/test/webhook"
	"github.com/openimsdk/protocol/constant"
)

func main() {
	port := flag.Int("port", 10006, "Webhook server listening port")
	sensitiveWord := flag.String("mask-word", "xxx", "Keyword to mask in callbackBeforeMsgModifyCommand")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	srv := webhook.NewServer()

	// 1. Log every incoming webhook request in real-time
	srv.OnRequest(func(command string, reqBody []byte, c *gin.Context) {
		log.Printf("[Webhook] %s %s | Command: %s | Body: %s",
			c.Request.Method, c.Request.URL.Path, command, string(reqBody))
	})

	// 2. Example: Custom hook for message content modification (sensitive word filtering)
	if *sensitiveWord != "" {
		webhook.RegisterCallback(srv, cbapi.CallbackBeforeMsgModifyCommand, func(c *gin.Context, req *cbapi.CallbackMsgModifyCommandReq) (*cbapi.CallbackMsgModifyCommandResp, error) {
			var resp cbapi.CallbackMsgModifyCommandResp
			if req.ContentType != constant.Text {
				return &resp, nil
			}
			var textElem struct {
				Content string `json:"content"`
			}
			if err := json.Unmarshal([]byte(req.Content), &textElem); err != nil {
				return nil, err
			}
			if strings.Contains(textElem.Content, *sensitiveWord) {
				textElem.Content = strings.ReplaceAll(textElem.Content, *sensitiveWord, strings.Repeat("*", len(*sensitiveWord)))
				masked, err := json.Marshal(&textElem)
				if err != nil {
					return nil, err
				}
				s := string(masked)
				resp.Content = &s
				log.Printf("[Webhook Mask] Masked content: %s -> %s", req.Content, s)
			}
			return &resp, nil
		})
	}

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("=================================================================")
	log.Printf(" OpenIM Mock Webhook Server is listening on %s", addr)
	log.Printf(" Callback URL in config/webhooks.yml:")
	log.Printf("   url: 'http://127.0.0.1:%d/callbackExample'", *port)
	log.Printf(" Supported callbacks: 48 registered OpenIM callback commands")
	log.Printf(" Press Ctrl+C to stop")
	log.Printf("=================================================================")

	go func() {
		if err := srv.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Webhook Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Webhook Server exited cleanly.")
}
