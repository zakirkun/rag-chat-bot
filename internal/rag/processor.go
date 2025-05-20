package rag

import (
	"context"
	"fmt"
	"rag-chat-bot/internal/database"
	"rag-chat-bot/internal/embedding"
	"rag-chat-bot/internal/model"
)

// Processor adalah komponen untuk memproses dokumen dalam sistem RAG
type Processor struct {
	db           *database.PostgresDB
	embeddingAPI *embedding.OpenAIEmbedding
}

// NewProcessor membuat instance Processor baru
func NewProcessor(db *database.PostgresDB, embeddingAPI *embedding.OpenAIEmbedding) *Processor {
	return &Processor{
		db:           db,
		embeddingAPI: embeddingAPI,
	}
}

// ProcessDocument memproses dokumen dan menyimpannya dengan embedding-nya
func (p *Processor) ProcessDocument(ctx context.Context, doc *model.Document) (int, error) {
	// Simpan dokumen ke database
	docID, err := p.db.SaveDocument(ctx, doc)
	if err != nil {
		return 0, fmt.Errorf("error saving document: %w", err)
	}

	// Generate embedding untuk dokumen
	embedding, err := p.embeddingAPI.CreateEmbedding(ctx, doc.Content)
	if err != nil {
		return 0, fmt.Errorf("error creating embedding: %w", err)
	}

	// Simpan embedding
	if err := p.db.SaveEmbedding(ctx, docID, embedding); err != nil {
		return 0, fmt.Errorf("error saving embedding: %w", err)
	}

	return docID, nil
}
