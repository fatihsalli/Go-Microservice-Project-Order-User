package elastic

import (
	"OrderUserProject/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/labstack/gommon/log"
	"strconv"
)

func SaveOrderElastic(order models.Order) {
	// Elasticsearch client
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Build the request body.
	data, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Error marshaling document: %s", err)
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "order-duplicate",
		DocumentID: strconv.Itoa(1),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), 1)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

}

func ReadFromElastic(orderID string) {
	// Elasticsearch client
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Read from elasticsearch
	res, err := es.Get("orders", orderID)
	if err != nil {
		log.Fatalf("Error getting order: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error getting order: %s", res.Status())
	}

	var result models.Order
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	fmt.Printf("Order retrieved: %+v", result)
}
