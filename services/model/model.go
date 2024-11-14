package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms/ollama"
)

type Model struct {
	URL               string `yaml:"url"`
	EmbedderURL       string `yaml:"embedder_url"`
	EmbedderModelName string `yaml:"embedder_model_name"`
	ModelName         string `yaml:"model_name"`
	Client            *ollama.LLM
	Embedder          *ollama.LLM
	BasePrompt        string `yaml:"base_prompt"`
}

// NewModel creates a new Model instance with optional URL and ModelName
func NewModel(
	url,
	modelName,
	embedderURL,
	embedderModelName,
	basePrompt string,
) (Model, error) {
	if url == "" {
		url = "http://localhost:11434"
	}
	if embedderModelName == "" {
		embedderModelName = "nomic-embed-text"
	}

	embedder, err := ollama.New(
		ollama.WithServerURL(embedderURL),
		ollama.WithModel(embedderModelName),
	)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create Ollama client: %v", err)
	}

	client, err := ollama.New(
		ollama.WithServerURL(url),
		ollama.WithModel(modelName),
	)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create Ollama client: %v", err)
	}

	return Model{
		URL:        url,
		ModelName:  modelName,
		Client:     client,
		Embedder:   embedder,
		BasePrompt: basePrompt,
	}, nil
}

// Add this method to generate embeddings
func (m *Model) GenerateEmbeddings(text string) ([][]float32, error) {

	embeddings, err := m.Embedder.CreateEmbedding(context.Background(), []string{text})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %v", err)
	}

	fmt.Println("Length of embeddings:", len(embeddings[0]))
	return embeddings, nil
}

func (m *Model) QueryWithContext(query string, ctx string) (string, error) {
	fmt.Println("Base prompt:", m.BasePrompt)
	query = strings.Replace(m.BasePrompt, "{context}", ctx, 1)
	query = strings.Replace(query, "{question}", query, 1)
	return m.Client.Call(context.Background(), query)
}
