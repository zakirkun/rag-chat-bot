# RAG Chat Bot

RAG (Retrieval-Augmented Generation) Chat Bot adalah aplikasi chatbot yang menggunakan teknologi RAG untuk memberikan jawaban yang akurat berdasarkan dokumen yang tersedia. Aplikasi ini menggunakan PostgreSQL dengan ekstensi pgvector untuk menyimpan dan mencari embedding vektor, serta OpenAI untuk menghasilkan embedding dan respons.

## Fitur

- Pencarian dokumen berbasis embedding vektor
- Penyimpanan dan pengambilan konteks percakapan
- Integrasi dengan OpenAI API untuk embedding dan chat
- Sistem RAG yang efisien untuk jawaban yang akurat
- Dukungan untuk metadata dokumen
- Pencarian similarity dengan HNSW index

## Persyaratan Sistem

- Go 1.21 atau lebih baru
- PostgreSQL 15 atau lebih baru dengan ekstensi pgvector
- Docker dan Docker Compose (opsional, untuk deployment)
- OpenAI API key

## Instalasi

1. Clone repository:
```bash
git clone https://github.com/zakirkun/rag-chat-bot.git
cd rag-chat-bot
```

2. Salin file konfigurasi:
```bash
cp .env.example .env
```

3. Edit file `.env` dan sesuaikan konfigurasi:
```env
# Server configuration
SERVER_PORT=8080

# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ragchatbot
DB_SSL_MODE=disable

# OpenAI API configuration
OPENAI_API_KEY=your-api-key
OPENAI_EMBEDDING_MODEL=text-embedding-ada-002
OPENAI_CHAT_MODEL=gpt-4
```

4. Jalankan database PostgreSQL menggunakan Docker Compose:
```bash
docker-compose up -d
```

5. Install dependencies:
```bash
go mod download
```

6. Jalankan aplikasi:
```bash
go run main.go
```

## Struktur Database

### Tabel Documents
- `id`: ID unik dokumen
- `title`: Judul dokumen
- `content`: Isi dokumen
- `metadata`: Metadata dokumen dalam format JSONB
- `created_at`: Waktu pembuatan dokumen

### Tabel Document Embeddings
- `id`: ID unik embedding
- `document_id`: Referensi ke dokumen
- `embedding`: Vektor embedding (1536 dimensi)
- `created_at`: Waktu pembuatan embedding

### Tabel Conversations
- `id`: ID unik percakapan
- `session_id`: ID sesi pengguna
- `created_at`: Waktu pembuatan percakapan

### Tabel Messages
- `id`: ID unik pesan
- `conversation_id`: Referensi ke percakapan
- `role`: Peran pengirim (user/assistant)
- `content`: Isi pesan
- `created_at`: Waktu pengiriman pesan

## Penggunaan API

### Endpoint Chat
```http
POST /api/chat
Content-Type: application/json

{
    "session_id": "unique-session-id",
    "message": "Pertanyaan atau pesan pengguna"
}
```

### Endpoint Upload Dokumen
```http
POST /api/documents
Content-Type: application/json

{
    "title": "Judul Dokumen",
    "content": "Isi dokumen",
    "metadata": {
        "source": "Sumber dokumen",
        "category": "Kategori dokumen"
    }
}
```

## Arsitektur

Aplikasi menggunakan arsitektur modular dengan komponen utama:

1. **Database Layer**
   - Menangani operasi database
   - Mengelola embedding vektor
   - Menyimpan riwayat percakapan

2. **RAG Layer**
   - Mengambil dokumen relevan
   - Membangun konteks untuk LLM
   - Mengelola interaksi dengan OpenAI

3. **API Layer**
   - Menangani request HTTP
   - Mengelola sesi pengguna
   - Menyediakan endpoint API

## Pengembangan

### Menjalankan Tests
```bash
go test ./...
```

### Menjalankan Linter
```bash
golangci-lint run
```

## Kontribusi

1. Fork repository
2. Buat branch fitur (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## Lisensi

Distribusikan di bawah lisensi MIT. Lihat `LICENSE` untuk informasi lebih lanjut.

## Kontak

Link Project: [https://github.com/zakirkun/rag-chat-bot](https://github.com/zakirkun/rag-chat-bot) 