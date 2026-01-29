package database

import (
	"bytes"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

var ESClient *elasticsearch.Client

func InitElasticsearch() {
	var err error

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
	}

	ESClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		panic("Failed to create Elasticsearch client: " + err.Error())
	}

	// Verify connection
	res, err := ESClient.Info()
	if err != nil {
		panic("Failed to connect to Elasticsearch: " + err.Error())
	}
	defer res.Body.Close()

	if res.IsError() {
		panic("Elasticsearch error: " + res.Status())
	}

	fmt.Println("Elasticsearch connection established")
}

// CreateIndexIfNotExists creates an index if it doesn't exist
func CreateIndexIfNotExists(indexName string, mapping string) error {
	// Check if index exists
	res, err := ESClient.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		// Create index with mapping
		res, err := ESClient.Indices.Create(indexName, ESClient.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))))
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("error creating index: %s", res.Status())
		}
	}

	return nil
}
