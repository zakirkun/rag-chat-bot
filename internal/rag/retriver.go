package rag

import (
	"context"
	"fmt"
	"log"
	"rag-chat-bot/internal/database"
	"rag-chat-bot/internal/embedding"
	"rag-chat-bot/internal/model"
	"strings"
)

// Retriever adalah komponen untuk mengambil dokumen yang relevan dalam sistem RAG
type Retriever struct {
	db           *database.PostgresDB
	embeddingAPI *embedding.OpenAIEmbedding
	maxResults   int
}

// NewRetriever membuat instance Retriever baru
func NewRetriever(db *database.PostgresDB, embeddingAPI *embedding.OpenAIEmbedding, maxResults int) *Retriever {
	if maxResults <= 0 {
		maxResults = 5 // Default value
	}

	return &Retriever{
		db:           db,
		embeddingAPI: embeddingAPI,
		maxResults:   maxResults,
	}
}

// RetrieveRelevantDocuments mengambil dokumen yang relevan berdasarkan query
func (r *Retriever) RetrieveRelevantDocuments(ctx context.Context, query string) ([]*model.DocumentWithScore, error) {
	// Generate embedding untuk query
	queryEmbedding, err := r.embeddingAPI.CreateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error creating query embedding: %w", err)
	}

	// Pastikan embedding memiliki dimensi yang benar
	if len(queryEmbedding) != 1536 {
		return nil, fmt.Errorf("invalid embedding dimension: expected 1536, got %d", len(queryEmbedding))
	}

	// Cari dokumen yang serupa berdasarkan embedding
	docs, err := r.db.FindSimilarDocuments(ctx, queryEmbedding, r.maxResults)
	if err != nil {
		return nil, fmt.Errorf("error finding similar documents: %w", err)
	}

	return docs, nil
}

// BuildPromptWithContext membangun prompt untuk model LLM dengan dokumen yang relevan sebagai konteks
func (r *Retriever) BuildPromptWithContext(ctx context.Context, userQuery string) (string, error) {
	relevantDocs, err := r.RetrieveRelevantDocuments(ctx, userQuery)
	if err != nil {
		return "", fmt.Errorf("error retrieving relevant documents: %w", err)
	}

	if len(relevantDocs) == 0 {
		return "", fmt.Errorf("no relevant documents found")
	}

	// Bangun prompt dengan konteks dari dokumen yang relevan
	var contextBuilder strings.Builder
	contextBuilder.WriteString("INFORMASI KONTEKS:\n\n")

	for i, doc := range relevantDocs {
		contextBuilder.WriteString(fmt.Sprintf("[Dokumen %d] (Relevansi: %.2f%%)\n", i+1, doc.Score*100))
		contextBuilder.WriteString(fmt.Sprintf("Judul: %s\n", doc.Title))
		contextBuilder.WriteString(fmt.Sprintf("Konten: %s\n\n", doc.Content))
	}

	// Bangun prompt lengkap
	promptTemplate := `Anda adalah asisten AI yang membantu pengguna dengan informasi berdasarkan dokumen yang tersedia. 
Tugas Anda adalah memberikan jawaban yang akurat, informatif, dan relevan berdasarkan informasi yang diberikan.

%s

PERTANYAAN/PERMINTAAN PENGGUNA:
%s

PANDUAN JAWABAN:
1. Berikan jawaban yang akurat dan relevan berdasarkan informasi yang tersedia
2. Jika informasi tidak cukup, jelaskan keterbatasan dan sarankan apa yang mungkin bisa membantu
3. Gunakan bahasa yang jelas dan mudah dipahami
4. Jika ada informasi yang bertentangan, jelaskan perbedaannya
5. Berikan sumber informasi yang digunakan (dokumen mana yang menjadi referensi)

Jawaban Anda harus:
- Langsung menjawab pertanyaan/permintaan pengguna
- Menggunakan informasi dari dokumen yang relevan
- Jelas dan terstruktur
- Jujur tentang keterbatasan informasi yang tersedia

Silakan berikan jawaban Anda:`

	prompt := fmt.Sprintf(promptTemplate, contextBuilder.String(), userQuery)
	return prompt, nil
}

// GenerateResponseFromContext menghasilkan respons dengan mengambil konteks yang relevan dan mengirimnya ke model LLM
func (r *Retriever) GenerateResponseFromContext(ctx context.Context, userQuery string, conversationHistory []model.ChatMessage) (string, error) {
	// Dapatkan prompt dengan konteks yang relevan
	contextPrompt, err := r.BuildPromptWithContext(ctx, userQuery)
	if err != nil {
		log.Printf("Error building prompt: %v", err)
		return "Saya tidak dapat menemukan informasi yang relevan untuk pertanyaan Anda. Bisakah Anda memberikan lebih banyak detail atau menanyakan hal lain?", nil
	}

	// Siapkan pesan untuk model LLM
	var messages []embedding.ChatCompletionMessage

	// Tambahkan sistem prompt
	systemPrompt := embedding.ChatCompletionMessage{
		Role:    "system",
		Content: "Anda adalah asisten AI yang membantu pengguna dengan informasi berdasarkan dokumen yang tersedia.",
	}
	messages = append(messages, systemPrompt)

	// Konversi history percakapan (jika ada)
	for _, msg := range conversationHistory {
		messages = append(messages, embedding.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Tambahkan konteks dan query
	contextMessage := embedding.ChatCompletionMessage{
		Role:    "user",
		Content: contextPrompt,
	}
	messages = append(messages, contextMessage)

	// Dapatkan respons dari model LLM
	response, err := r.embeddingAPI.ChatCompletion(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	return response, nil
}
