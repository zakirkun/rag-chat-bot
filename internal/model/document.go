package model

import (
	"encoding/json"
	"time"
)

// Document merepresentasikan dokumen sumber untuk sistem RAG
type Document struct {
	ID        int                    `json:"id"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata"`
	CreatedAt time.Time              `json:"created_at"`
}

// DocumentWithScore merepresentasikan dokumen dengan skor kesamaan
type DocumentWithScore struct {
	Document
	Score float64 `json:"score"`
}

// ToJSON mengkonversi Document ke JSON string
func (d *Document) ToJSON() (string, error) {
	bytes, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON mengkonversi JSON string ke Document
func DocumentFromJSON(data string) (*Document, error) {
	var doc Document
	err := json.Unmarshal([]byte(data), &doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
