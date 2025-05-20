-- Aktifkan ekstensi pgvector
CREATE EXTENSION IF NOT EXISTS vector;

-- Tabel untuk menyimpan dokumen sumber
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabel untuk menyimpan embedding vektor dari dokumen
CREATE TABLE document_embeddings (
    id SERIAL PRIMARY KEY,
    document_id INTEGER REFERENCES documents(id) ON DELETE CASCADE,
    embedding VECTOR(1536), -- Menggunakan dimensi 1536 untuk OpenAI embeddings
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indeks untuk pencarian vektor dengan metode HNSW (efisien untuk pencarian knn)
CREATE INDEX ON document_embeddings USING hnsw (embedding vector_cosine_ops);

-- Tabel untuk menyimpan percakapan
CREATE TABLE conversations (
    id SERIAL PRIMARY KEY,
    session_id TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabel untuk menyimpan pesan dalam percakapan
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER REFERENCES conversations(id) ON DELETE CASCADE,
    role TEXT NOT NULL CHECK (role IN ('user', 'assistant')),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indeks untuk membantu kueri
CREATE INDEX idx_document_embeddings_document_id ON document_embeddings(document_id);
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_conversations_session_id ON conversations(session_id);