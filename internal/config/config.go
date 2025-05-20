package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menyimpan semua konfigurasi aplikasi
type Config struct {
	// Server
	ServerPort int

	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// OpenAI
	OpenAIAPIKey         string
	OpenAIEmbeddingModel string
	OpenAIChatModel      string
}

// LoadConfig memuat konfigurasi dari variabel lingkungan
func LoadConfig() (*Config, error) {
	// Coba muat .env file jika ada
	_ = godotenv.Load()

	config := &Config{}

	// Server config
	serverPort, err := strconv.Atoi(getEnvOrDefault("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}
	config.ServerPort = serverPort

	// Database config
	config.DBHost = getEnvOrDefault("DB_HOST", "localhost")
	dbPort, err := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}
	config.DBPort = dbPort
	config.DBUser = getEnvOrDefault("DB_USER", "postgres")
	config.DBPassword = getEnvOrDefault("DB_PASSWORD", "postgres")
	config.DBName = getEnvOrDefault("DB_NAME", "ragchatbot")
	config.DBSSLMode = getEnvOrDefault("DB_SSL_MODE", "disable")

	// OpenAI config
	config.OpenAIAPIKey = getEnvOrDefault("OPENAI_API_KEY", "")
	if config.OpenAIAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY is required")
	}
	config.OpenAIEmbeddingModel = getEnvOrDefault("OPENAI_EMBEDDING_MODEL", "text-embedding-ada-002")
	config.OpenAIChatModel = getEnvOrDefault("OPENAI_CHAT_MODEL", "gpt-3.5-turbo")

	return config, nil
}

// GetPostgresConnectionString mengembalikan string koneksi untuk PostgreSQL
func (c *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// Helper untuk mendapatkan variabel lingkungan dengan nilai default
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
