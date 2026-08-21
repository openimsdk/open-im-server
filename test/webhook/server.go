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

package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// RequestRecord stores information about a received webhook request.
type RequestRecord struct {
	Timestamp time.Time
	Command   string
	Body      []byte
	Headers   http.Header
}

// Server is a mock webhook HTTP server for OpenIM callbacks.
type Server struct {
	Engine        *gin.Engine
	handlers      map[string]gin.HandlerFunc
	mu            sync.RWMutex
	httpServer    *http.Server
	recordHistory bool
	history       []RequestRecord
	requestHook   func(command string, reqBody []byte, c *gin.Context)
}

// NewServer creates a new mock webhook server instance.
func NewServer() *Server {
	engine := gin.New()
	engine.Use(gin.Recovery())

	s := &Server{
		Engine:   engine,
		handlers: make(map[string]gin.HandlerFunc),
	}

	engine.POST("/callbackExample/:command", func(c *gin.Context) {
		s.dispatch(c, c.Param("command"))
	})
	engine.POST("/:command", func(c *gin.Context) {
		s.dispatch(c, c.Param("command"))
	})

	s.RegisterDefaultHandlers()
	return s
}

// SetRecordHistory enables or disables recording incoming request history.
func (s *Server) SetRecordHistory(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recordHistory = enabled
}

// OnRequest sets a callback function invoked whenever any webhook request arrives.
func (s *Server) OnRequest(hook func(command string, reqBody []byte, c *gin.Context)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestHook = hook
}

// GetHistory returns a copy of all recorded requests.
func (s *Server) GetHistory() []RequestRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]RequestRecord, len(s.history))
	copy(result, s.history)
	return result
}

// ClearHistory clears all recorded request history.
func (s *Server) ClearHistory() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = nil
}

// Start runs the webhook HTTP server on the given address (e.g. ":10006").
func (s *Server) Start(addr string) error {
	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.Engine,
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the running HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) dispatch(c *gin.Context, command string) {
	s.mu.RLock()
	h, ok := s.handlers[command]
	s.mu.RUnlock()
	if ok {
		h(c)
		return
	}
	c.JSON(http.StatusOK, cbapi.CommonCallbackResp{})
}

// ServeHTTP conforms to the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Engine.ServeHTTP(w, r)
}

// RegisterCallback registers a typed handler for a callback command.
func RegisterCallback[Req any, Resp any](s *Server, command string, fn func(c *gin.Context, req *Req) (*Resp, error)) {
	handler := func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		s.mu.Lock()
		if s.recordHistory {
			s.history = append(s.history, RequestRecord{
				Timestamp: time.Now(),
				Command:   command,
				Body:      body,
				Headers:   c.Request.Header.Clone(),
			})
		}
		hook := s.requestHook
		s.mu.Unlock()

		if hook != nil {
			hook(command, body, c)
		}

		var req Req
		if len(body) > 0 {
			if err := json.Unmarshal(body, &req); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("unmarshal error: %v", err))
				return
			}
		}
		resp, err := fn(c, &req)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if resp != nil {
			c.JSON(http.StatusOK, resp)
		} else {
			c.JSON(http.StatusOK, cbapi.CommonCallbackResp{})
		}
	}

	s.mu.Lock()
	s.handlers[command] = handler
	s.mu.Unlock()
}
