package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"loan-mgt/g-cram/internal/config"
	"loan-mgt/g-cram/internal/db"
	"loan-mgt/g-cram/internal/server"
	"loan-mgt/g-cram/internal/server/ws"
	"loan-mgt/g-cram/internal/service"

	"github.com/gin-gonic/gin"
)

func loadConfig() *config.Config {
	return config.New()
}

func newDB(cfg *config.Config) *db.Store {
	store := db.NewStore(cfg)
	return store
}

func newAMQPConnection(cfg *config.Config) (*service.AMQPConnection, error) {
	return service.NewAMQPConnection(cfg, "downloader")
}

func newWebSocketManager() *ws.WebSocketManager {
	return ws.NewWebSocketManager()
}

func startListener(ws *ws.WebSocketManager, store *db.Store, cfg *config.Config) {
	go service.StartUploadDoneListener(ws, store, cfg)
}

func newRouter(store *db.Store, amqpConn *service.AMQPConnection, ws *ws.WebSocketManager, cfg *config.Config) *gin.Engine {
	return server.NewRouter(store, amqpConn, ws, cfg)
}

func startServer(cfg *config.Config, router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	go func() {
		log.Printf("Starting server on port %d", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exited gracefully")
}

func main() {
	cfg := loadConfig()
	store := newDB(cfg)
	amqpConn, err := newAMQPConnection(cfg)
	if err != nil {
		panic(err)
	}
	ws := newWebSocketManager()
	startListener(ws, store, cfg)

	router := newRouter(store, amqpConn, ws, cfg)
	startServer(cfg, router)
}
