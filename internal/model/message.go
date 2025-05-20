package model

import (
	"time"
)

// Message merepresentasikan pesan dalam percakapan
type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	Role           string    `json:"role"` // "user" atau "assistant"
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

// ChatMessage adalah format pesan yang digunakan untuk berkomunikasi dengan OpenAI API
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest adalah struktur permintaan chat
type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// ChatResponse adalah struktur respons chat
type ChatResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	SessionID string `json:"session_id"`
}

// CreateDocumentRequest adalah struktur permintaan untuk membuat dokumen baru
type CreateDocumentRequest struct {
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CreateDocumentResponse adalah struktur respons saat membuat dokumen baru
type CreateDocumentResponse struct {
	Success bool `json:"success"`
	DocID   int  `json:"doc_id"`
}
