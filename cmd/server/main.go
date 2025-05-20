package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rag-chat-bot/internal/api"
	"rag-chat-bot/internal/config"
	"rag-chat-bot/internal/database"
	"rag-chat-bot/internal/embedding"
	"rag-chat-bot/internal/rag"
	"rag-chat-bot/internal/service"
	"syscall"
	"time"
)

func main() {
	// Load konfigurasi
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Inisialisasi koneksi database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Inisialisasi OpenAI API client
	openaiClient := embedding.NewOpenAIEmbedding(cfg)

	// Inisialisasi komponen RAG
	ragProcessor := rag.NewProcessor(db, openaiClient)
	ragRetriever := rag.NewRetriever(db, openaiClient, 5) // Ambil 5 dokumen teratas

	// Inisialisasi service
	chatService := service.NewChatService(db, ragRetriever)

	// Inisialisasi handler dan router
	handler := api.NewHandler(chatService, ragProcessor)
	router := handler.SetupRouter()

	// Konfigurasi server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	// Inisialisasi channel untuk shutdown signal
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Jalankan server dalam goroutine
	go func() {
		log.Printf("Server starting on port %d", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
		done <- true
	}()

	// Tunggu signal untuk shutdown
	<-quit
	log.Println("Server is shutting down...")

	// Beri waktu 30 detik untuk request yang sedang berjalan
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Matikan server dengan graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	// Tunggu hingga server benar-benar berhenti
	<-done
	log.Println("Server stopped")
}
