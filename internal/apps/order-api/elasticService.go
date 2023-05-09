package order_api

import (
	"OrderUserProject/internal/configs"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/labstack/gommon/log"
)

type ElasticService struct {
	Config        *configs.Config
	ElasticClient *elasticsearch.Client
}

func NewElasticService(config *configs.Config) *ElasticService {
	// client with default config
	cfg := elasticsearch.Config{
		Addresses: []string{
			config.Elasticsearch.Addresses["Address 1"],
		},
	}

	elasticClient, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Errorf("Error creating the client: ", err)
	}

	elasticService := &ElasticService{Config: config, ElasticClient: elasticClient}
	return elasticService
}

func (e *ElasticService) GetFromElasticsearch(req OrderGetRequest) ([]interface{}, error) {

	searchBody := make(map[string]interface{})
	query := make(map[string]interface{})

	// Creating query for exact filters
	if len(req.ExactFilters) > 0 {
		boolQuery := make(map[string]interface{})
		mustClauses := make([]map[string]interface{}, 0)

		for field, values := range req.ExactFilters {
			if len(values) > 0 {
				mustClause := make(map[string]interface{})
				mustClause["terms"] = map[string]interface{}{
					field: values,
				}
				mustClauses = append(mustClauses, mustClause)
			}
		}

		boolQuery["must"] = mustClauses
		query["bool"] = boolQuery
	}

	// Creating query for match
	if len(req.Match) > 0 {

	}

	searchBody["query"] = query

	if len(req.Sort) > 0 {
		for field, value := range req.Sort {
			if value == -1 {
				searchBody["sort"] = map[string]interface{}{
					field: "desc",
				}
			} else if value == 1 {
				searchBody["sort"] = map[string]interface{}{
					field: "asc",
				}
			}
		}
	}

	if len(req.Fields) > 0 {
		var idCheck bool
		for _, value := range req.Fields {
			if value == "id" {
				idCheck = true
				break
			} else {
				idCheck = false
			}
		}

		if !idCheck {
			req.Fields = append(req.Fields, "id")
		}

		searchBody["_source"] = req.Fields
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(searchBody); err != nil {
		fmt.Println("Error encoding the query: ", err)
		return nil, err
	}

	res, err := e.ElasticClient.Search(
		e.ElasticClient.Search.WithIndex(e.Config.Elasticsearch.IndexName["OrderSave"]),
		e.ElasticClient.Search.WithBody(buf),
	)
	if err != nil {
		fmt.Println("Error executing the search: ", err)
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		fmt.Println("Error executing the decode: ", err)
		return nil, err
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		fmt.Println("Error executing the decode: ", err)
		return nil, err
	}

	var orders []interface{}

	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	for _, hit := range hits {

		// Casting with type assertion
		source, ok := hit.(map[string]interface{})["_source"]
		if !ok {
			fmt.Println("Source not found in the hit", err)
			return nil, err
		}
		orders = append(orders, source)
	}

	return orders, nil
}
