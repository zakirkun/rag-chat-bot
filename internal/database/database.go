package database

import (
	"context"
	"database/sql"
	"fmt"
	"rag-chat-bot/internal/config"
	"rag-chat-bot/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // Driver PostgreSQL
)

// PostgresDB menyediakan akses ke database PostgreSQL
type PostgresDB struct {
	pool *pgxpool.Pool
}

// NewPostgresDB membuat koneksi baru ke database PostgreSQL
func NewPostgresDB(cfg *config.Config) (*PostgresDB, error) {
	connString := cfg.GetPostgresConnectionString()

	// Membuat pool koneksi
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	// Melakukan koneksi ke database
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Memastikan koneksi berhasil
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &PostgresDB{pool: pool}, nil
}

// Close menutup koneksi database
func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

// SaveDocument menyimpan dokumen baru ke database
func (db *PostgresDB) SaveDocument(ctx context.Context, doc *model.Document) (int, error) {
	var docID int

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Menyimpan dokumen
	err = tx.QueryRow(ctx,
		"INSERT INTO documents (title, content, metadata) VALUES ($1, $2, $3) RETURNING id",
		doc.Title, doc.Content, doc.Metadata).Scan(&docID)
	if err != nil {
		return 0, fmt.Errorf("error inserting document: %w", err)
	}

	// Commit transaksi
	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return docID, nil
}

// SaveEmbedding menyimpan embedding untuk dokumen ke database
func (db *PostgresDB) SaveEmbedding(ctx context.Context, docID int, embedding []float32) error {
	// Convert float32 slice to string representation of array
	vectorStr := "["
	for i, v := range embedding {
		if i > 0 {
			vectorStr += ","
		}
		vectorStr += fmt.Sprintf("%f", v)
	}
	vectorStr += "]"

	_, err := db.pool.Exec(ctx,
		"INSERT INTO document_embeddings (document_id, embedding) VALUES ($1, $2::vector)",
		docID, vectorStr)
	if err != nil {
		return fmt.Errorf("error inserting embedding: %w", err)
	}

	return nil
}

// FindSimilarDocuments mencari dokumen yang serupa berdasarkan embedding kueri
func (db *PostgresDB) FindSimilarDocuments(ctx context.Context, queryEmbedding []float32, limit int) ([]*model.DocumentWithScore, error) {
	// Convert float32 slice to string representation of array
	vectorStr := "["
	for i, v := range queryEmbedding {
		if i > 0 {
			vectorStr += ","
		}
		vectorStr += fmt.Sprintf("%f", v)
	}
	vectorStr += "]"

	rows, err := db.pool.Query(ctx, `
		SELECT d.id, d.title, d.content, d.metadata, 
		       1 - (e.embedding <=> $1::vector) as similarity_score
		FROM document_embeddings e
		JOIN documents d ON e.document_id = d.id
		ORDER BY e.embedding <=> $1::vector
		LIMIT $2
	`, vectorStr, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying similar documents: %w", err)
	}
	defer rows.Close()

	var results []*model.DocumentWithScore

	for rows.Next() {
		var doc model.DocumentWithScore
		var metadataJSON []byte

		if err := rows.Scan(&doc.ID, &doc.Title, &doc.Content, &metadataJSON, &doc.Score); err != nil {
			return nil, fmt.Errorf("error scanning document row: %w", err)
		}

		// Parse metadata jika diperlukan
		doc.Metadata = make(map[string]interface{})
		// Jika ingin memproses metadata, tambahkan kode di sini

		results = append(results, &doc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// SaveConversation menyimpan percakapan baru dan mengembalikan ID-nya
func (db *PostgresDB) SaveConversation(ctx context.Context, sessionID string) (int, error) {
	var conversationID int
	err := db.pool.QueryRow(ctx,
		"INSERT INTO conversations (session_id) VALUES ($1) RETURNING id",
		sessionID).Scan(&conversationID)
	if err != nil {
		return 0, fmt.Errorf("error creating conversation: %w", err)
	}
	return conversationID, nil
}

// SaveMessage menyimpan pesan dalam percakapan
func (db *PostgresDB) SaveMessage(ctx context.Context, msg *model.Message) error {
	_, err := db.pool.Exec(ctx,
		"INSERT INTO messages (conversation_id, role, content) VALUES ($1, $2, $3)",
		msg.ConversationID, msg.Role, msg.Content)
	if err != nil {
		return fmt.Errorf("error saving message: %w", err)
	}
	return nil
}

// GetConversationMessages mengambil semua pesan dalam percakapan
func (db *PostgresDB) GetConversationMessages(ctx context.Context, conversationID int) ([]*model.Message, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT id, conversation_id, role, content, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`, conversationID)
	if err != nil {
		return nil, fmt.Errorf("error querying conversation messages: %w", err)
	}
	defer rows.Close()

	var messages []*model.Message

	for rows.Next() {
		var msg model.Message

		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning message row: %w", err)
		}

		messages = append(messages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return messages, nil
}

// GetConversationBySessionID menemukan ID percakapan berdasarkan session ID
func (db *PostgresDB) GetConversationBySessionID(ctx context.Context, sessionID string) (int, error) {
	var conversationID int
	err := db.pool.QueryRow(ctx, `
		SELECT id FROM conversations
		WHERE session_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, sessionID).Scan(&conversationID)

	if err != nil {
		if err == sql.ErrNoRows || err == pgx.ErrNoRows {
			// Jika tidak ada percakapan dengan session ID ini, buat baru
			return db.SaveConversation(ctx, sessionID)
		}
		return 0, fmt.Errorf("error finding conversation: %w", err)
	}

	return conversationID, nil
}
