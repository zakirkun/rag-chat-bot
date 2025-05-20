package api

import (
	"encoding/json"
	"log"
	"net/http"
	"rag-chat-bot/internal/model"
	"rag-chat-bot/internal/rag"
	"rag-chat-bot/internal/service"
)

// Handler mengelola permintaan API
type Handler struct {
	chatService *service.ChatService
	processor   *rag.Processor
}

// NewHandler membuat instance Handler baru
func NewHandler(chatService *service.ChatService, processor *rag.Processor) *Handler {
	return &Handler{
		chatService: chatService,
		processor:   processor,
	}
}

// HandleChat menangani permintaan chat
func (h *Handler) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi permintaan
	if req.Message == "" {
		http.Error(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}

	if req.SessionID == "" {
		// Generate session ID sederhana jika tidak disediakan
		req.SessionID = generateSessionID()
	}

	// Proses pesan
	resp, err := h.chatService.ProcessUserMessage(r.Context(), &req)
	if err != nil {
		log.Printf("Error processing message: %v", err)
		http.Error(w, "Error processing message", http.StatusInternalServerError)
		return
	}

	// Kirim respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleAddDocument menangani penambahan dokumen baru
func (h *Handler) HandleAddDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi permintaan
	if req.Title == "" || req.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Buat dan proses dokumen
	doc := &model.Document{
		Title:    req.Title,
		Content:  req.Content,
		Metadata: req.Metadata,
	}

	docID, err := h.processor.ProcessDocument(r.Context(), doc)
	if err != nil {
		log.Printf("Error processing document: %v", err)
		http.Error(w, "Error processing document", http.StatusInternalServerError)
		return
	}

	// Kirim respons
	resp := model.CreateDocumentResponse{
		Success: true,
		DocID:   docID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// HandleGetConversation menangani pengambilan riwayat percakapan
func (h *Handler) HandleGetConversation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	if sessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// Ambil riwayat percakapan
	messages, err := h.chatService.GetConversationHistory(r.Context(), sessionID)
	if err != nil {
		log.Printf("Error getting conversation history: %v", err)
		http.Error(w, "Error getting conversation history", http.StatusInternalServerError)
		return
	}

	// Kirim respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    true,
		"session_id": sessionID,
		"messages":   messages,
	})
}

// generateSessionID menghasilkan ID sesi sederhana
func generateSessionID() string {
	// Dalam implementasi sebenarnya, gunakan UUID atau ID yang lebih aman
	// contoh: github.com/google/uuid
	return "sess_" + randString(16)
}

// randString menghasilkan string acak dengan panjang tertentu
func randString(n int) string {
	// Implementasi sederhana, untuk produksi gunakan crypto/rand
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[i%len(letters)]
	}
	return string(b)
}
