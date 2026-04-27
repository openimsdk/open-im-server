package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"redpacket/config"
	"redpacket/internal/chain"
	"redpacket/internal/handler"
	"redpacket/internal/model"
	"redpacket/internal/repository"
	"redpacket/internal/service"
	"redpacket/router"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfgFile := ""
	if len(os.Args) > 1 {
		cfgFile = os.Args[1]
	}
	config.Load(cfgFile)
	cfg := &config.Cfg

	// Connect to database
	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&model.RedPacket{},
		&model.RedPacketClaim{},
		&model.RedPacketClaimAuth{},
		&model.RedPacketRefund{},
	); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	// Create blockchain client
	chainClient, err := chain.NewClient(
		cfg.Chain.RPCURL,
		cfg.Chain.ContractAddress,
		cfg.Chain.ChainID,
		cfg.Chain.SignerPrivateKey,
		cfg.Chain.ConfigAdminPrivateKey,
	)
	if err != nil {
		log.Printf("Warning: failed to create chain client: %v (continuing with mock mode)", err)
		// Continue without blockchain for now - can be configured later
	}

	// Create repository and service
	repo := repository.New(db)
	rpSvc := service.NewRedPacketService(repo, chainClient, cfg.Chain.SignerPrivateKey)

	// Create TRON client if configured
	var tronClient *chain.TronClient
	if cfg.Tron.FullNodeURL != "" {
		abiJSON, err := chain.ExtractABIFromEmbeddedArtifact()
		if err != nil {
			log.Printf("Warning: failed to load ABI for TRON: %v", err)
		} else {
			tronClient, err = chain.NewTronClient(
				cfg.Tron.FullNodeURL,
				cfg.Tron.ContractBase58,
				cfg.Tron.OwnerBase58,
				cfg.Tron.PrivateKeyHex,
				abiJSON,
				cfg.Tron.FeeLimit,
			)
			if err != nil {
				log.Printf("Warning: failed to create TRON client: %v", err)
				tronClient = nil
			} else {
				log.Println("✅ TRON client initialized successfully")
			}
		}
	}

	// Create admin service and handler
	adminSvc := service.NewAdminService(chainClient, tronClient)
	adminHandler := handler.NewAdminHandler(adminSvc)

	// Create user handler
	rpHandler := handler.NewRedPacketHandler(rpSvc)

	// Setup router
	r := gin.Default()
	router.Setup(r, rpHandler, adminHandler)

	// Start blockchain indexers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ETH Indexer
	if chainClient != nil {
		ethIndexer := chain.NewIndexer(chainClient, repo, cfg.Indexer.PollInterval, 0)
		ethIndexer.Start(ctx)
		log.Println("📡 ETH Blockchain event indexer started")
	}

	// TRON Indexer (Production-grade)
	if tronClient != nil {
		tronIndexer := chain.NewTronIndexer(tronClient, repo, cfg.Indexer.PollInterval, 0)
		tronIndexer.Start(ctx)
		log.Println("📡 TRON Blockchain event indexer started (Production mode)")
	}

	// Start HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	go func() {
		log.Printf("🚀 RedPacket service listening on :%d", cfg.Server.Port)
		log.Printf("📋 Health check: http://localhost:%d/health", cfg.Server.Port)
		log.Printf("📋 API docs: see backend-api.md")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced shutdown: %v", err)
	}
	log.Println("server stopped")
}

func openDB(cfg *config.Config) (*gorm.DB, error) {
	switch cfg.DB.Driver {
	case "mysql":
		return gorm.Open(mysql.Open(cfg.DB.DSN), &gorm.Config{})
	case "sqlite", "":
		return gorm.Open(sqlite.Open(cfg.DB.DSN), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported db.driver: %s", cfg.DB.Driver)
	}
}
