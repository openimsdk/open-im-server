// Copyright © 2024 OpenIM. All rights reserved.
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

package encryption

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/internal/rpc/encryption/stores"
	kdisc "github.com/openimsdk/open-im-server/v3/pkg/common/discovery"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
)

type Server struct {
	*Config
	keysManager   *KeysManager
	discoveryConn discovery.Conn
	httpServer    *http.Server
}

func Start(ctx context.Context, cfg *Config) error {
	log.ZInfo(ctx, "encryption server start")

	// Initialize service registry
	client, err := kdisc.NewDiscoveryRegister(&cfg.Discovery, nil)
	if err != nil {
		return err
	}

	// Initialize MongoDB
	mongoClient, err := mongoutil.NewMongoDB(ctx, cfg.MongodbConfig.Build())
	if err != nil {
		return err
	}

	// Get the specific database
	db := mongoClient.GetDB()

	// Initialize stores
	identityStore := stores.NewIdentityStore(db)
	preKeyStore := stores.NewPreKeyStore(db)
	signedPreKeyStore := stores.NewSignedPreKeyStore(db)

	// Initialize managers
	keysManager := NewKeysManager(identityStore, preKeyStore, signedPreKeyStore)

	server := &Server{
		Config:        cfg,
		keysManager:   keysManager,
		discoveryConn: client,
	}

	// Setup HTTP server
	if err := server.setupHTTPServer(); err != nil {
		return err
	}

	// Start HTTP server
	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.ZError(ctx, "HTTP server failed to start", err)
		}
	}()

	log.ZInfo(ctx, "encryption server started successfully", "port", cfg.RpcConfig.Ports[0])

	// Keep the service running
	select {
	case <-ctx.Done():
		return server.shutdown(ctx)
	}
}

func (s *Server) setupHTTPServer() error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())
	router.Use(s.loggingMiddleware())

	// API routes
	api := router.Group("/api/v1/encryption")
	{
		// Key management endpoints
		api.GET("/prekeys/:user_id/:device_id", s.GetPreKeys)
		api.POST("/prekeys/:user_id/:device_id", s.SetPreKeys)
		api.GET("/prekeys/:user_id/:device_id/count", s.GetPreKeyCount)
		api.GET("/identity/:user_id/:device_id", s.GetIdentityKey)

		// Health check
		api.GET("/health", s.HealthCheck)
	}

	// Create HTTP server
	port := fmt.Sprintf(":%d", s.Config.RpcConfig.Ports[0])
	s.httpServer = &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return nil
}

func (s *Server) shutdown(ctx context.Context) error {
	log.ZInfo(ctx, "shutting down encryption server")

	// Shutdown HTTP server
	if s.httpServer != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			log.ZError(ctx, "failed to shutdown HTTP server", err)
		}
	}

	return nil
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] %s %s %d %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
		)
	})
}

// HealthCheck handles GET /api/v1/encryption/health
func (s *Server) HealthCheck(c *gin.Context) {
	c.JSON(200, APIResponse{
		Code:    0,
		Message: "success",
		Data: map[string]interface{}{
			"status":    "healthy",
			"mode":      s.Config.GetEncryptionMode(),
			"timestamp": time.Now().Unix(),
		},
	})
}