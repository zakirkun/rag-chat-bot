# 🤖 RAG Chat Bot

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://www.postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-required-2496ED?style=flat&logo=docker)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Retrieval-Augmented Generation (RAG) Chat Bot that provides accurate responses based on available documents. This application uses PostgreSQL with pgvector extension for storing and searching vector embeddings, and OpenAI for generating embeddings and responses.

## ✨ Features

- 🔍 Vector-based document search
- 💾 Conversation context storage and retrieval
- 🤝 OpenAI API integration for embeddings and chat
- 🎯 Efficient RAG system for accurate responses
- 📝 Document metadata support
- ⚡ HNSW index for similarity search

## 🛠️ System Requirements

- Go 1.21 or newer
- PostgreSQL 15 or newer with pgvector extension
- Docker and Docker Compose (optional, for deployment)
- OpenAI API key

## 🚀 Installation

1. Clone repository:
```bash
git clone https://github.com/zakirkun/rag-chat-bot.git
cd rag-chat-bot
```

2. Copy configuration file:
```bash
cp .env.example .env
```

3. Edit `.env` file and adjust configuration:
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

4. Run PostgreSQL database using Docker Compose:
```bash
docker-compose up -d
```

5. Install dependencies:
```bash
go mod download
```

6. Run the application:
```bash
go run main.go
```

## 📊 Database Structure

### Documents Table
- `id`: Unique document ID
- `title`: Document title
- `content`: Document content
- `metadata`: Document metadata in JSONB format
- `created_at`: Document creation timestamp

### Document Embeddings Table
- `id`: Unique embedding ID
- `document_id`: Reference to document
- `embedding`: Vector embedding (1536 dimensions)
- `created_at`: Embedding creation timestamp

### Conversations Table
- `id`: Unique conversation ID
- `session_id`: User session ID
- `created_at`: Conversation creation timestamp

### Messages Table
- `id`: Unique message ID
- `conversation_id`: Reference to conversation
- `role`: Sender role (user/assistant)
- `content`: Message content
- `created_at`: Message timestamp

## 🔌 API Usage

### Chat Endpoint
```http
POST /api/chat
Content-Type: application/json

{
    "session_id": "unique-session-id",
    "message": "User question or message"
}
```

### Document Upload Endpoint
```http
POST /api/documents
Content-Type: application/json

{
    "title": "Document Title",
    "content": "Document content",
    "metadata": {
        "source": "Document source",
        "category": "Document category"
    }
}
```

## 🏗️ Architecture

The application uses a modular architecture with main components:

1. **Database Layer**
   - Handles database operations
   - Manages vector embeddings
   - Stores conversation history

2. **RAG Layer**
   - Retrieves relevant documents
   - Builds context for LLM
   - Manages OpenAI interactions

3. **API Layer**
   - Handles HTTP requests
   - Manages user sessions
   - Provides API endpoints

## 👨‍💻 Development

### Running Tests
```bash
go test ./...
```

### Running Linter
```bash
golangci-lint run
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.