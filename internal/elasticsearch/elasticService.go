package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

func ElasticSave() {
	// Create Elasticsearch client
	esCfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	es, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		fmt.Printf("Error creating Elasticsearch client: %s", err.Error())
		return
	}

	// Create a new product
	product := Product{
		ID:          "1",
		Name:        "Product 1",
		Description: "This is product 1",
		Price:       100,
	}

	// Serialize product to JSON bytes
	jsonProduct, err := json.Marshal(product)
	if err != nil {
		fmt.Printf("Error marshalling product to JSON: %s", err.Error())
		return
	}

	// Create index request
	indexReq := esapi.IndexRequest{
		Index:      "products",
		DocumentID: product.ID,
		Body:       bytes.NewReader(jsonProduct),
	}

	// Send index request
	indexRes, err := indexReq.Do(context.Background(), es)
	if err != nil {
		fmt.Printf("Error sending index request: %s", err.Error())
		return
	}
	defer indexRes.Body.Close()

	// Read response from index request
	var indexResBody map[string]interface{}
	if err := json.NewDecoder(indexRes.Body).Decode(&indexResBody); err != nil {
		fmt.Printf("Error parsing index response: %s", err.Error())
		return
	}
	fmt.Printf("Index response: %v\n", indexResBody)

	// Wait for 1 second to let Elasticsearch index the document
	time.Sleep(1 * time.Second)

	// Search all products
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	jsonSearchBody, err := json.Marshal(searchBody)
	if err != nil {
		fmt.Printf("Error marshalling search body to JSON: %s", err.Error())
		return
	}
	searchReq := esapi.SearchRequest{
		Index: []string{"products"},
		Body:  bytes.NewReader(jsonSearchBody),
	}
	searchRes, err := searchReq.Do(context.Background(), es)
	if err != nil {
		fmt.Printf("Error sending search request: %s", err.Error())
		return
	}
	defer searchRes.Body.Close()

	// Read response from search request
	var searchResBody map[string]interface{}
	if err := json.NewDecoder(searchRes.Body).Decode(&searchResBody); err != nil {
		fmt.Printf("Error parsing search response: %s", err.Error())
		return
	}
	fmt.Printf("Search response: %v\n", searchResBody)
}
