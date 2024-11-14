package store

import (
	"clutch/common"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantStore struct {
	Client           *qdrant.Client
	Host             string
	Port             string
	CollectionSizes  map[string]int
	CollectionCounts map[string]int
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

func (s *QdrantStore) CheckCollectionExists(collectionName string) (bool, error) {
	return s.Client.CollectionExists(context.Background(), collectionName)
}

func (s *QdrantStore) CreateCollection(collectionName string, size int) {
	fmt.Println("Creating collection:", collectionName)
	if exists, err := s.CheckCollectionExists(collectionName); err != nil || exists {
		fmt.Println("Collection already exists")
		return
	}
	if s.CollectionCounts == nil {
		s.CollectionCounts = make(map[string]int)
		s.CollectionSizes = make(map[string]int)
	}
	s.Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(size),
			Distance: qdrant.Distance_Cosine,
		}),
	})
	s.CollectionSizes[collectionName] = size
	s.CollectionCounts[collectionName] = 0
}

func (s *QdrantStore) InsertDocument(index string, body map[string]interface{}) {
	fmt.Println("Inserting document into collection:", index)
	fmt.Println("Document:", body)

	// Create flattened map
	flattenedBody := make(map[string]interface{})
	common.FlattenMap(flattenedBody, "", body)

	// Convert document to string for embedding
	jsonStr, err := json.Marshal(flattenedBody)
	if err != nil {
		log.Printf("Error marshaling document: %v", err)
		return
	}

	cfg := common.GetConfig()
	model := cfg.Model
	// Generate embeddings
	embeddings, err := model.GenerateEmbeddings(string(jsonStr))
	if err != nil {
		log.Printf("Error generating embeddings: %v", err)
		return
	}

	var id uint64
	if exists, err := s.CheckCollectionExists(index); err != nil || !exists {
		s.CreateCollection(index, len(embeddings[0]))
		id = uint64(s.CollectionCounts[index])
	} else {
		id, err = s.Client.Count(context.Background(), &qdrant.CountPoints{
			CollectionName: index,
		})
		if err != nil {
			log.Printf("Error getting collection size: %v", err)
			return
		}
	}

	// Insert document with embeddings
	res, err := s.Client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: index,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(id),
				Vectors: qdrant.NewVectors(embeddings[0]...),
				Payload: qdrant.NewValueMap(body),
			},
		},
	})
	if err != nil || res.Status != qdrant.UpdateStatus_Acknowledged {
		log.Printf("Error upserting document: %v", err)
		return
	}
	fmt.Println("---------- Done Inserting into Qdrant ----------")
	return
}

func (s *QdrantStore) Query(collectionName string, query string) (r map[string]interface{}) {
	// Filters
	// ScoreThreshold
	// Embedd Query
	cfg := common.GetConfig()
	model := cfg.Model
	// Generate vector
	vector, err := model.GenerateEmbeddings(query)
	if err != nil {
		log.Printf("Error generating embeddings: %v", err)
		return
	}
	// fmt.Println("Vector:", vector)
	// search points

	filter := qdrant.Filter{
		Must: []*qdrant.Condition{
			qdrant.NewMatch("location", "field_1"),
		},
	}

	res, err := s.Client.QueryBatch(context.Background(), &qdrant.QueryBatchPoints{
		CollectionName: collectionName,
		QueryPoints: []*qdrant.QueryPoints{
			{
				CollectionName: collectionName,
				Filter:         &filter,
			},
		},
	})
	if err != nil {
		log.Printf("Error querying collection: %v", err)
		return nil
	}

	limit := uint64(3)
	results, err := s.Client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(vector[0]...),
		// Filter: &qdrant.Filter{
		// 	Should: []*qdrant.Condition{
		// 		qdrant.NewMatch("machine_id", "4"),
		// 	},
		// },
		Limit: &limit,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Results: ", results)

	return map[string]interface{}{
		"results": res,
	}
}

func (s *QdrantStore) DeleteIndex(index string) {
	// Delete collection
	s.Client.DeleteCollection(context.Background(), index)
	return
}

func (s *QdrantStore) GetResults(searchResult map[string]interface{}) (res []map[string]interface{}) {
	// For Qdrant, this might be a no-op since collections/indices are typically
	// created during initialization
	return nil
}

// func (s *QdrantStore) Test() {
// 	fmt.Println("Testing Qdrant client")
// 	// s.Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
// 	// 	CollectionName: "testing_collection",
// 	// 	VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
// 	// 		Size:     4,
// 	// 		Distance: qdrant.Distance_Cosine,
// 	// 	}),
// 	// })
// 	fmt.Println("Collection created")
// 	operationInfo, err := s.Client.Upsert(context.Background(), &qdrant.UpsertPoints{
// 		CollectionName: "testing_collection",
// 		Points: []*qdrant.PointStruct{
// 			{
// 				Id:      qdrant.NewIDNum(1),
// 				Vectors: qdrant.NewVectors(0.05, 0.61, 0.76, 0.74),
// 				Payload: qdrant.NewValueMap(map[string]any{"city": "Berlin"}),
// 			},
// 			{
// 				Id:      qdrant.NewIDNum(2),
// 				Vectors: qdrant.NewVectors(0.19, 0.81, 0.75, 0.11),
// 				Payload: qdrant.NewValueMap(map[string]any{"city": "London"}),
// 			},
// 			{
// 				Id:      qdrant.NewIDNum(3),
// 				Vectors: qdrant.NewVectors(0.36, 0.55, 0.47, 0.94),
// 				Payload: qdrant.NewValueMap(map[string]any{"city": "Moscow"}),
// 			},
// 			// Truncated
// 		},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(operationInfo)
// }
