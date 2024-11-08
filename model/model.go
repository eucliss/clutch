package model

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms/ollama"
)

type Model struct {
	URL       string `yaml:"url"`
	ModelName string `yaml:"model_name"`
	Client    *ollama.LLM
}

// NewModel creates a new Model instance with optional URL and ModelName
func NewModel(url, modelName string) (Model, error) {
	if url == "" {
		url = "http://localhost:11434"
	}
	if modelName == "" {
		modelName = "nomic-embed-text"
	}

	client, err := ollama.New(
		ollama.WithServerURL(url),
		ollama.WithModel(modelName),
	)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create Ollama client: %v", err)
	}

	return Model{
		URL:       url,
		ModelName: modelName,
		Client:    client,
	}, nil
}

// Add this method to generate embeddings
func (m *Model) GenerateEmbeddings(text string) ([][]float32, error) {

	embeddings, err := m.Client.CreateEmbedding(context.Background(), []string{text})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %v", err)
	}

	fmt.Println("Length of embeddings:", len(embeddings[0]))
	return embeddings, nil
}
