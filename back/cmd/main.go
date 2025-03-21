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
)

func main() {
	// Load configuration
	cfg := config.New()

	store := db.NewStore(cfg)
	defer store.Close()

	amqpConn, err := service.NewAMQPConnection(config.New(), "cramer")
	if err != nil {
		panic(err)
	}
	defer amqpConn.Conn.Close()

	ws := ws.NewWebSocketManager()
	defer ws.Close()

	go service.StartListener(ws, store, cfg)

	// Set up router
	router := server.NewRouter(store, amqpConn, ws, cfg)

	// Configure HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// Start server in a goroutine
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
