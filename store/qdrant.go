package store

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantStore struct {
	Client *qdrant.Client
	Host   string
	Port   string
}

func (s *QdrantStore) Initialize() {
	fmt.Println("Initializing Qdrant client")
	// The Go client uses Qdrant's gRPC interface
	port, err := strconv.Atoi(s.Port)
	if err != nil {
		log.Fatalf("Failed to convert port to integer: %v", err)
	}
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: s.Host,
		Port: port,
	})
	if err != nil {
		log.Fatalf("Failed to create Qdrant client: %v", err)
	}
	s.Client = client
}

// func (s *QdrantStore) Initialize() *qdrant.Client {
// 	fmt.Println("Initializing Qdrant client")
// 	// The Go client uses Qdrant's gRPC interface
// 	port, err := strconv.Atoi(s.Port)
// 	if err != nil {
// 		log.Fatalf("Failed to convert port to integer: %v", err)
// 	}
// 	client, err := qdrant.NewClient(&qdrant.Config{
// 		Host: s.Host,
// 		Port: port,
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to create Qdrant client: %v", err)
// 	}
// 	s.Client = client
// 	return client
// }

func (s *QdrantStore) CheckCollectionExists(collectionName string) (bool, error) {
	return s.Client.CollectionExists(context.Background(), collectionName)
}

func (s *QdrantStore) CreateCollection(collectionName string) {
	fmt.Println("Creating collection:", collectionName)
	if exists, err := s.CheckCollectionExists(collectionName); err != nil || exists {
		fmt.Println("Collection already exists")
		return
	}
	s.Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
	})
}

func (s *QdrantStore) InsertDocument(index string, body map[string]interface{}) {
	fmt.Println("Inserting document into collection:", index)
	fmt.Println("Document:", body)
	// s.Client.Upsert(context.Background(), &qdrant.UpsertPoints{
	// 	CollectionName: collectionName,
	// 	Points:         points,
	// })
	return
}

// ... existing code ...
func (s *QdrantStore) CreateIndices(indices ...Index) {
	// For Qdrant, this might be a no-op since collections/indices are typically
	// created during initialization
	return
}

func (s *QdrantStore) Query(collectionName string, query string) (r map[string]interface{}) {
	// For Qdrant, this might be a no-op since collections/indices are typically
	// created during initialization
	return nil
}

func (s *QdrantStore) DeleteIndex(index string) {
	// For Qdrant, this might be a no-op since collections/indices are typically
	// created during initialization
	return
}

func (s *QdrantStore) GetResults(searchResult map[string]interface{}) (res []map[string]interface{}) {
	// For Qdrant, this might be a no-op since collections/indices are typically
	// created during initialization
	return nil
}

func (s *QdrantStore) Test() {
	fmt.Println("Testing Qdrant client")
	// s.Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
	// 	CollectionName: "testing_collection",
	// 	VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
	// 		Size:     4,
	// 		Distance: qdrant.Distance_Cosine,
	// 	}),
	// })
	fmt.Println("Collection created")
	operationInfo, err := s.Client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "testing_collection",
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(1),
				Vectors: qdrant.NewVectors(0.05, 0.61, 0.76, 0.74),
				Payload: qdrant.NewValueMap(map[string]any{"city": "Berlin"}),
			},
			{
				Id:      qdrant.NewIDNum(2),
				Vectors: qdrant.NewVectors(0.19, 0.81, 0.75, 0.11),
				Payload: qdrant.NewValueMap(map[string]any{"city": "London"}),
			},
			{
				Id:      qdrant.NewIDNum(3),
				Vectors: qdrant.NewVectors(0.36, 0.55, 0.47, 0.94),
				Payload: qdrant.NewValueMap(map[string]any{"city": "Moscow"}),
			},
			// Truncated
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(operationInfo)
}

// type StoreConfig struct {
// 	Location string // "db/http_ca.crt"
// 	Address  string // "https://localhost:9200"
// 	username string
// 	password string
// 	cfg      elasticsearch8.Config
// 	caCert   []byte
// 	es       *elasticsearch8.Client
// 	indicies []string
// }

// type Index struct {
// 	Name    string
// 	Mapping string
// }
