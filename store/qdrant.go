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

// func convertToQueryResultsToJson(res []*qdrant.BatchResult) (r map[string]interface{}) {
// 	fmt.Println("---------- Converting to JSON ----------")
// 	fmt.Println("Length of results:", len(res[0].GetResult()))
// 	r = make(map[string]interface{})
// 	for _, value := range res[0].GetResult() {
// 		payload := value.GetPayload()
// 		fmt.Println("Payload:", payload)
// 		// Convert payload to json

// 		jsonPayload, err := json.Marshal(payload)
// 		if err != nil {
// 			log.Printf("Error marshalling payload: %v", err)
// 			return nil
// 		}
// 		fmt.Println("JSON Payload:", string(jsonPayload))
// 		fmt.Println("================")

//		}
//		return nil
//	}

// func convertToJson(queryResults []*qdrant.BatchResult) (r map[string]interface{}) {
// 	r = make(map[string]interface{})
// 	if len(queryResults) == 0 || len(queryResults[0].GetResult()) == 0 {
// 		return nil
// 	}
// 	results := make([]map[string]interface{}, 0)
// 	qdrantResults := queryResults[0].GetResult()
// 	for _, value := range qdrantResults {
// 		payload := value.GetPayload()
// 		payloadMap := make(map[string]interface{})
// 	}
// 	return nil
// }

func getQdrantValue(value *qdrant.Value) (r interface{}) {
	switch val := value.GetKind().(type) {
	case *qdrant.Value_StringValue:
		return val.StringValue
	case *qdrant.Value_IntegerValue:
		return val.IntegerValue
	case *qdrant.Value_BoolValue:
		return val.BoolValue
	case *qdrant.Value_StructValue:
		return convertStructValueToJson(val.StructValue)
	}
	return nil
}

func convertStructValueToJson(structValue *qdrant.Struct) (r map[string]interface{}) {
	r = make(map[string]interface{})
	for k, v := range structValue.GetFields() {
		r[k] = getQdrantValue(v)
	}
	return r
}

func convertToQueryResultsToJson(res []*qdrant.BatchResult) (r map[string]interface{}) {
	r = make(map[string]interface{})
	if len(res) == 0 || len(res[0].GetResult()) == 0 {
		return nil
	}

	// Convert each result to a map
	results := make([]map[string]interface{}, 0)
	for _, value := range res[0].GetResult() {
		payload := value.GetPayload()

		// Convert Qdrant's payload to a regular map
		payloadMap := make(map[string]interface{})
		for k, v := range payload {
			// Handle different value types from Qdrant
			payloadMap[k] = getQdrantValue(v)
		}
		results = append(results, payloadMap)
	}

	r["hits"] = results
	fmt.Println("---------- Done Converting to JSON ----------")
	fmt.Println("JSON :", r)
	fmt.Println("length of results:", len(results))
	return r
}

func (s *QdrantStore) Query(collectionName string, query string) (r map[string]interface{}) {

	cfg := common.GetConfig()
	model := cfg.Model
	// Generate vectorp
	vector, err := model.GenerateEmbeddings(query)
	if err != nil {
		log.Printf("Error generating embeddings: %v", err)
		return
	}

	filter := qdrant.Filter{
		Must: []*qdrant.Condition{
			qdrant.NewMatch("location", "field_1"),
		},
	}

	limit := uint64(4)
	scoreThreshold := float32(0.6)

	queryPoints := &qdrant.QueryPoints{
		CollectionName: collectionName,
		Filter:         &filter,
		Query:          qdrant.NewQuery(vector[0]...),
		Limit:          &limit,
		ScoreThreshold: &scoreThreshold,
	}
	queryPoints.WithPayload = &qdrant.WithPayloadSelector{
		SelectorOptions: &qdrant.WithPayloadSelector_Enable{
			Enable: true,
		},
	}

	res, err := s.Client.QueryBatch(context.Background(), &qdrant.QueryBatchPoints{
		CollectionName: collectionName,
		QueryPoints:    []*qdrant.QueryPoints{queryPoints},
	})
	if err != nil {
		log.Printf("Error querying collection: %v", err)
		return nil
	}
	r = convertToQueryResultsToJson(res)
	PrettyPrint(r["hits"].([]map[string]interface{}))
	return r
}

func PrettyPrint(v []map[string]interface{}) {
	fmt.Println("---------- Pretty Print ----------")
	for _, mp := range v {
		fmt.Println("----------")
		for k, v := range mp {
			fmt.Println(k, v)
		}
	}
	fmt.Println("---------- Done Pretty Print ----------")
}

func (s *QdrantStore) DeleteIndex(index string) {
	// Delete collection
	s.Client.DeleteCollection(context.Background(), index)
	return
}

func (s *QdrantStore) GetResults(searchResult map[string]interface{}) (res []map[string]interface{}) {
	return nil
	// getPoints := &qdrant.GetPoints{
	// 	CollectionName: "clutch_testing_events",
	// 	Ids:            []*qdrant.PointId{qdrant.NewIDNum(3)},
	// }
	// getPoints.WithPayload = &qdrant.WithPayloadSelector{
	// 	SelectorOptions: &qdrant.WithPayloadSelector_Enable{
	// 		Enable: true,
	// 	},
	// }
	// // Get point 1 in collection
	// point, err := s.Client.Get(context.Background(), getPoints)
	// if err != nil {
	// 	log.Printf("Error getting point: %v", err)
	// 	return nil
	// }
	// fmt.Println("Point:", point)
	// fmt.Println("Payload:", point[0].GetPayload())
	// return nil
	// // Take a json of {id: {num: X}} and get those results from Qdrant
	// collectionName := searchResult["collection"].(string)
	// // convert docs to []map[string]interface{}

	// docs := searchResult["results"].([]*qdrant.BatchResult)
	// fmt.Println("Search results:", docs)
	// ids := []uint64{}
	// ids_points := []*qdrant.PointId{}
	// for _, doc := range docs {
	// 	sc := doc.GetResult()
	// 	for _, point := range sc {
	// 		id := point.GetId()
	// 		id_num := id.GetNum()
	// 		uuid := id.GetUuid()
	// 		// Skip if id_num is 0 and uuid is not empty
	// 		if id_num == 0 && uuid != "" {
	// 			continue
	// 		} else {
	// 			ids = append(ids, id_num)
	// 			ids_points = append(ids_points, id)
	// 		}
	// 	}
	// }
	// fmt.Println("IDS:", ids)

	// points, err := s.Client.Get(context.Background(), &qdrant.GetPoints{
	// 	CollectionName: collectionName,
	// 	Ids:            ids_points,
	// })
	// if err != nil {
	// 	log.Printf("Error getting point: %v", err)
	// 	return nil
	// }
	// fmt.Println("Points:", points)
	// for _, point := range points {
	// 	payload := point.GetPayload()
	// 	fmt.Println("Payload:", payload)
	// 	// convert payload to map[string]interface{}
	// 	payloadMap := make(map[string]interface{})
	// 	for k, v := range payload {
	// 		payloadMap[k] = v
	// 	}
	// 	res = append(res, payloadMap)
	// }
	// return res
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
