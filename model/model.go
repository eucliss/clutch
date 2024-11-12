package model

import (
	"context"
	"fmt"
	"log"

	"net/url"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/qdrant"

	vqdrant "github.com/tmc/langchaingo/vectorstores/qdrant"
)

type Model struct {
	URL       string `yaml:"url"`
	ModelName string `yaml:"model_name"`
	LLM       *ollama.LLM
	Embedder  *embeddings.EmbedderImpl
	Store     *qdrant.Store
}

// NewModel creates a new Model instance with optional URL and ModelName
func NewModel(host string, modelName string) (Model, error) {
	if host == "" {
		host = "http://localhost:11434"
	}
	if modelName == "" {
		modelName = "nomic-embed-text"
	}
	embedderModel, err := ollama.New(
		ollama.WithServerURL(host),
		ollama.WithModel(modelName),
	)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create Ollama client: %v", err)
	}
	embedder, err := embeddings.NewEmbedder(embedderModel)
	if err != nil {
		return Model{}, fmt.Errorf("failed to create embedder: %v", err)
	}

	llm, err := ollama.New(
		ollama.WithModel("llama3.2"),
		ollama.WithServerURL(host))
	if err != nil {
		return Model{}, fmt.Errorf("failed to create Ollama client: %v", err)
	}

	qdrantUrl := fmt.Sprintf("http://%s:%s", "localhost", "6333")
	qdUrl, err := url.Parse(qdrantUrl)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
	}

	store, err := vqdrant.New(
		vqdrant.WithURL(*qdUrl),
		vqdrant.WithAPIKey(""),
		vqdrant.WithCollectionName("new_collection_testing"),
		vqdrant.WithEmbedder(embedder),
	)
	fmt.Println("Store:", store)
	if err != nil {
		log.Fatalf("Failed to create Qdrant store: %v", err)
	}

	return Model{
		URL:       host,
		ModelName: modelName,
		LLM:       llm,
		Embedder:  embedder,
		Store:     &store,
	}, nil
}

// Add this method to generate embeddings
func (m *Model) GenerateEmbeddings(text string) ([][]float32, error) {

	embeddings, err := m.LLM.CreateEmbedding(context.Background(), []string{text})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %v", err)
	}

	fmt.Println("Length of embeddings:", len(embeddings[0]))
	return embeddings, nil
}

func (m *Model) SetStore(store *qdrant.Store) {
	m.Store = store
}

func (m *Model) RetrieveDocuments(query string) ([]schema.Document, error) {
	fmt.Println("Retrieving documents for query:", query)
	options := []vectorstores.Option{
		vectorstores.WithScoreThreshold(0.10),
	}
	retriever := vectorstores.ToRetriever(m.Store, 20, options...)
	docRetrieved, err := retriever.GetRelevantDocuments(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return docRetrieved, nil
}

// func (m *Model) Ask(ctx context.Context, llm llms.Model, docRetrieved []schema.Document, prompt string) (string, error) {
func (m *Model) Ask(prompt string) (string, error) {

	docRetrieved, err := m.RetrieveDocuments(prompt)
	if err != nil {
		return "", err
	}
	if len(docRetrieved) == 0 {
		return "No relevant documents found.", nil
	}

	history := memory.NewChatMessageHistory()
	for _, doc := range docRetrieved {
		history.AddAIMessage(context.Background(), doc.PageContent)
	}
	conversation := memory.NewConversationBuffer(memory.WithChatHistory(history))

	executor := agents.NewExecutor(
		agents.NewConversationalAgent(m.LLM, nil),
		nil,
		agents.WithMemory(conversation),
	)
	options := []chains.ChainCallOption{
		chains.WithTemperature(0.8),
	}
	res, err := chains.Run(context.Background(), executor, prompt, options...)
	if err != nil {
		return "", err
	}

	return res, nil
}
