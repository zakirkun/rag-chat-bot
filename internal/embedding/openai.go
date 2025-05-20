package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"rag-chat-bot/internal/config"
)

// OpenAIEmbedding adalah klien untuk membuat embedding menggunakan OpenAI API
type OpenAIEmbedding struct {
	apiKey         string
	embeddingModel string
	chatModel      string
}

// EmbeddingRequest adalah struktur untuk permintaan embedding ke OpenAI API
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// EmbeddingResponse adalah struktur untuk respons embedding dari OpenAI API
type EmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// NewOpenAIEmbedding membuat klien baru untuk OpenAI embedding
func NewOpenAIEmbedding(cfg *config.Config) *OpenAIEmbedding {
	return &OpenAIEmbedding{
		apiKey:         cfg.OpenAIAPIKey,
		embeddingModel: cfg.OpenAIEmbeddingModel,
		chatModel:      cfg.OpenAIChatModel,
	}
}

// CreateEmbedding membuat embedding vektor dari teks
func (o *OpenAIEmbedding) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Siapkan permintaan
	reqBody := EmbeddingRequest{
		Model: o.embeddingModel,
		Input: []string{text},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Kirim permintaan ke OpenAI API
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Baca respons
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %s", string(body))
	}

	// Parse respons
	var embeddingResp EmbeddingResponse
	if err := json.Unmarshal(body, &embeddingResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(embeddingResp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return embeddingResp.Data[0].Embedding, nil
}

// ChatCompletionRequest adalah struktur untuk permintaan chat completion ke OpenAI API
type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
}

// ChatCompletionMessage adalah format pesan untuk chat completion
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse adalah struktur untuk respons chat completion dari OpenAI API
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ChatCompletion membuat chat completion dengan OpenAI API
func (o *OpenAIEmbedding) ChatCompletion(ctx context.Context, messages []ChatCompletionMessage) (string, error) {
	// Siapkan permintaan
	reqBody := ChatCompletionRequest{
		Model:    o.chatModel,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	// Kirim permintaan ke OpenAI API
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Baca respons
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}

	// Parse respons
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no completion returned")
	}

	return chatResp.Choices[0].Message.Content, nil
}
