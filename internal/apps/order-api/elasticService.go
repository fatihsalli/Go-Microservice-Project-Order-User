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
	// Get config for generic endpoint
	config := configs.GetGenericEndpointConfig("elasticsearch")

	searchBody := make(map[string]interface{})
	query := make(map[string]interface{})
	boolQuery := make(map[string]interface{})
	mustClauses := make([]map[string]interface{}, 0)
	mustNotClauses := make([]map[string]interface{}, 0)

	// Creating query for exact filters
	if len(req.ExactFilters) > 0 {
		for field, values := range req.ExactFilters {
			if len(values) > 0 {
				mustClause := make(map[string]interface{})
				mustClause["terms"] = map[string]interface{}{
					config.ExactFilterArea[field]: values,
				}
				mustClauses = append(mustClauses, mustClause)
			}
		}
		boolQuery["must"] = mustClauses
		query["bool"] = boolQuery
	}

	// Creating query for match
	if len(req.Match) > 0 {
		for _, model := range req.Match {
			if config.MatchFilterParameter[model.Parameter] == "eq" {
				// => "match":{"total":1800}
				mustClause := make(map[string]interface{})
				mustClause["match"] = map[string]interface{}{
					config.ExactFilterArea[model.MatchField]: model.Value,
				}
				mustClauses = append(mustClauses, mustClause)
				// => "must": [{"match": {"total": 1800}}]
				boolQuery["must"] = mustClauses
				// =>  "bool": {"must": [{"match": {"total": 1800}}]}
				query["bool"] = boolQuery
			} else if config.MatchFilterParameter[model.Parameter] == "ne" {
				// => "match":{"total":1800}
				mustNotClause := make(map[string]interface{})
				mustNotClause["match"] = map[string]interface{}{
					config.ExactFilterArea[model.MatchField]: model.Value,
				}
				mustNotClauses = append(mustNotClauses, mustNotClause)
				// => "must_not": [{"match": {"total": 1800}}]
				boolQuery["must_not"] = mustNotClauses
				// =>  "bool": {"must_not": [{"match": {"total": 1800}}]}
				query["bool"] = boolQuery
			} else if config.MatchFilterParameter[model.Parameter] == "gt" || config.MatchFilterParameter[model.Parameter] == "gte" ||
				config.MatchFilterParameter[model.Parameter] == "lt" || config.MatchFilterParameter[model.Parameter] == "lte" {
				// => "lt":2000
				parameterQuery := make(map[string]interface{})
				parameterQuery[config.MatchFilterParameter[model.Parameter]] = model.Value
				// => "total":{"lt":2000}
				rangeQuery := make(map[string]interface{})
				rangeQuery[model.MatchField] = parameterQuery
				// => "range": {"total":{"lt": 2000}}
				mustClause := make(map[string]interface{})
				mustClause["range"] = rangeQuery
				mustClauses = append(mustClauses, mustClause)
				// => "must": [{"range": {"total":{"lt": 2000}}}]
				boolQuery["must"] = mustClauses
				// =>  "bool": {"must": [{"range": {"total":{"lt": 2000}}}]}
				query["bool"] = boolQuery
			} else if config.MatchFilterParameter[model.Parameter] == "in" {
				// => "terms":{"total":[1800,2000,2200]}
				mustClause := make(map[string]interface{})
				mustClause["terms"] = map[string]interface{}{
					config.ExactFilterArea[model.MatchField]: model.Value,
				}
				mustClauses = append(mustClauses, mustClause)
				// => "must": ["terms":{"total":[1800,2000,2200]}]
				boolQuery["must"] = mustClauses
				// =>  "bool": {"must": ["terms":{"total":[1800,2000,2200]}]}
				query["bool"] = boolQuery
			} else if config.MatchFilterParameter[model.Parameter] == "nin" {
				// => "terms":{"total":[1900,2000,2200]}
				mustNotClause := make(map[string]interface{})
				mustNotClause["terms"] = map[string]interface{}{
					config.ExactFilterArea[model.MatchField]: model.Value,
				}
				mustNotClauses = append(mustNotClauses, mustNotClause)
				// => "must_not": ["terms":{"total":[1900,2000,2200]}]
				boolQuery["must_not"] = mustNotClauses
				// =>  "bool": {"must_not": ["terms":{"total":[1900,2000,2200]}]}
				query["bool"] = boolQuery
			} else if config.MatchFilterParameter[model.Parameter] == "exists" {
				// => "exists":{"field":"total"}
				mustClause := make(map[string]interface{})
				mustClause["exists"] = map[string]interface{}{
					"field": model.Value,
				}
				mustClauses = append(mustClauses, mustClause)
				// => "must": ["exists":{"field":"total"}]
				boolQuery["must"] = mustClauses
				// =>  "bool": {"must": ["exists":{"field":"total"}]}
				query["bool"] = boolQuery
			} else if config.MatchFilterParameter[model.Parameter] == "regex" {
				// => "regexp":{"product.name": ".*a.*"}
				mustClause := make(map[string]interface{})
				mustClause["regexp"] = map[string]interface{}{
					model.MatchField: model.Value,
				}
				mustClauses = append(mustClauses, mustClause)
				// => "must": ["regexp":{"product.name": ".*a.*"}]
				boolQuery["must"] = mustClauses
				// =>  "bool": {"must": ["regexp":{"product.name": ".*a.*"}]}
				query["bool"] = boolQuery
			}
		}
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
