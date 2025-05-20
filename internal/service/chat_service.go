package service

import (
	"context"
	"fmt"
	"log"
	"rag-chat-bot/internal/database"
	"rag-chat-bot/internal/model"
	"rag-chat-bot/internal/rag"
)

// ChatService mengelola layanan percakapan
type ChatService struct {
	db        *database.PostgresDB
	retriever *rag.Retriever
}

// NewChatService membuat instance ChatService baru
func NewChatService(db *database.PostgresDB, retriever *rag.Retriever) *ChatService {
	return &ChatService{
		db:        db,
		retriever: retriever,
	}
}

// ProcessUserMessage memproses pesan pengguna dan menghasilkan respons
func (s *ChatService) ProcessUserMessage(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error) {
	// Dapatkan atau buat percakapan baru berdasarkan session ID
	conversationID, err := s.db.GetConversationBySessionID(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("error getting conversation: %w", err)
	}

	// Simpan pesan pengguna
	userMsg := &model.Message{
		ConversationID: conversationID,
		Role:           "user",
		Content:        req.Message,
	}

	if err := s.db.SaveMessage(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("error saving user message: %w", err)
	}

	// Ambil riwayat percakapan
	messages, err := s.db.GetConversationMessages(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error getting conversation history: %w", err)
	}

	// Konversi pesan ke format yang diperlukan oleh retriever
	var chatMessages []model.ChatMessage
	for _, msg := range messages {
		chatMessages = append(chatMessages, model.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Generate respons menggunakan RAG retriever
	response, err := s.retriever.GenerateResponseFromContext(ctx, req.Message, chatMessages)
	if err != nil {
		log.Printf("Error generating response: %v, will return generic response", err)
		response = "Maaf, saya mengalami kesulitan dalam memproses permintaan Anda. Silakan coba lagi atau tanyakan dengan cara yang berbeda."
	}

	// Simpan respons asisten
	assistantMsg := &model.Message{
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        response,
	}

	if err := s.db.SaveMessage(ctx, assistantMsg); err != nil {
		log.Printf("Error saving assistant message: %v", err)
		// Lanjutkan meskipun ada kesalahan penyimpanan
	}

	// Kembalikan respons
	return &model.ChatResponse{
		Success:   true,
		Message:   response,
		SessionID: req.SessionID,
	}, nil
}

// GetConversationHistory mengambil riwayat percakapan
func (s *ChatService) GetConversationHistory(ctx context.Context, sessionID string) ([]*model.Message, error) {
	// Dapatkan ID percakapan dari session ID
	conversationID, err := s.db.GetConversationBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error getting conversation: %w", err)
	}

	// Ambil semua pesan dalam percakapan
	messages, err := s.db.GetConversationMessages(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error getting conversation messages: %w", err)
	}

	return messages, nil
}
