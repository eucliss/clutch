package store

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
)

type ElasticStore struct {
	Location string // "db/http_ca.crt"
	Address  string // "https://localhost:9200"
	username string
	password string
	cfg      elasticsearch8.Config
	caCert   []byte
	es       *elasticsearch8.Client
	indicies []string
}

type Index struct {
	Name    string
	Mapping string
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (c *ElasticStore) GetIndices() ([]string, error) {
	res, err := c.es.Indices.Get(
		[]string{"_all"},
		c.es.Indices.Get.WithHuman(),
		c.es.Indices.Get.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting indices: %w", err)
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	indices := make([]string, 0, len(result))
	for index := range result {
		indices = append(indices, index)
	}

	return indices, nil
}

func (c *ElasticStore) IndexExists(indexName string) (bool, error) {
	if contains(c.indicies, indexName) {
		return true, nil
	}
	_, err := c.es.Indices.Exists([]string{indexName})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *ElasticStore) Username() string {
	return c.username
}

func (c *ElasticStore) SetUsername(u string) {
	c.username = u
}

func (c *ElasticStore) SetPassword(p string) {
	c.password = p
}

func (c *ElasticStore) Cert() []byte {
	return c.caCert
}

func (c *ElasticStore) Initialize() {
	// Gather the CA certificate
	c.gatherCert(c.Location)

	// ElasticStoreure Elasticsearch
	c.ConfigureES()

	// Create Elasticsearch client
	c.createClient()
}

func (c *ElasticStore) gatherCert(location string) ([]byte, error) {
	caCert, err := os.ReadFile(location) // Replace with the path to your CA certificate
	if err != nil {
		log.Fatalf("Error reading CA certificate: %s", err)
	}
	c.caCert = caCert
	return caCert, err
}

func (c *ElasticStore) ConfigureES() {
	// Elasticsearch ElasticStoreuration
	cfg := elasticsearch8.Config{
		Addresses: []string{
			c.Address, // If running Go app locally
			// "http://elasticsearch:9200", // If running Go app inside Docker
		},
		// If you set up authentication, include Username and Password
		Username: c.username, // Default user for Elasticsearch Docker image
		Password: c.password, // Replace with actual password if set
		CACert:   c.caCert,
	}
	c.cfg = cfg
}

func (c *ElasticStore) createClient() (*elasticsearch8.Client, error) {
	// Create Elasticsearch client
	es, err := elasticsearch8.NewClient(c.cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	c.es = es
	return es, err
}

func (c *ElasticStore) CreateIndices(indicies ...Index) {
	for _, id := range indicies {
		_, err := c.es.Indices.Create(
			id.Name,
			c.es.Indices.Create.WithBody(strings.NewReader(id.Mapping)),
		)
		if err != nil {
			log.Fatalf("Error creating index: %s", err)
		}
		c.indicies = append(c.indicies, id.Name)
	}
}

func (c *ElasticStore) InsertDocument(index string, body map[string]interface{}) {
	doc := body
	docJSON, _ := json.Marshal(doc)
	res, err := c.es.Index(index, bytes.NewReader(docJSON))
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	c.refreshIndex(index)
	defer res.Body.Close()
	fmt.Println("Document indexed successfully")
}

func (c *ElasticStore) refreshIndex(index string) {
	_, err := c.es.Indices.Refresh(c.es.Indices.Refresh.WithIndex(index))
	if err != nil {
		log.Fatalf("Error refreshing index: %s", err)
	}
}

func (c *ElasticStore) Query(index string, query string) (r map[string]interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := c.es.Search(
		c.es.Search.WithContext(context.Background()),
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(strings.NewReader(query)),
		c.es.Search.WithTrackTotalHits(true),
		c.es.Search.WithPretty(),
	)
	fmt.Println("res:", res)
	if err != nil {
		log.Fatalf("Error searching index: %s", err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	defer c.PrintResults(r)
	return
}

func (c *ElasticStore) GetResults(searchResult map[string]interface{}) []map[string]interface{} {
	count := searchResult["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	res := make([]map[string]interface{}, int(count))
	index := 0
	for _, hit := range searchResult["hits"].(map[string]interface{})["hits"].([]interface{}) {
		// Parse the attributes/fields of the document
		doc := hit.(map[string]interface{})
		v := doc["_source"].(map[string]interface{})
		res[index] = v
		index++
	}
	return res
}

func (c *ElasticStore) PrintResults(searchResult map[string]interface{}) {
	for _, hit := range searchResult["hits"].(map[string]interface{})["hits"].([]interface{}) {

		// Parse the attributes/fields of the document
		doc := hit.(map[string]interface{})

		// The "_source" data is another map interface nested inside of doc
		source := doc["_source"]
		fmt.Println("doc _source:", reflect.TypeOf(source))

		// Get the document's _id and print it out along with _source data
		docID := doc["_id"]
		fmt.Println("docID:", docID)
		fmt.Println("_source:", source)
	} // end of response iteration
}

func (c *ElasticStore) DeleteByQuery(indexName string, query string) {
	// Perform the DeleteByQuery request
	res, err := c.es.DeleteByQuery(
		[]string{indexName},                                  // The index from which to delete documents
		strings.NewReader(query),                             // The query to match all documents
		c.es.DeleteByQuery.WithContext(context.Background()), // Add context for timeout handling
	)
	if err != nil {
		log.Fatalf("Error deleting documents: %s", err)
	}
	defer res.Body.Close()
}

func (c *ElasticStore) DeleteIndex(indexName string) {
	// Perform the DeleteIndex request
	res, err := c.es.Indices.Delete(
		[]string{indexName}, // Index names
	)
	if err != nil {
		log.Fatalf("Error deleting index: %s", err)
	}
	defer res.Body.Close()
}
