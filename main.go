package main

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"log"
)

type Product struct {
	Name        string   `json:"name,omitempty"`
	Price       float64  `json:"price,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func NewElasticSearchClient() *elasticsearch.Client {
	es, err := elasticsearch.NewDefaultClient()

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	_, err = es.Ping()
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %s", err)
	}
	return es
}

func CreateNewIndex(es *elasticsearch.Client, indexName string) (*esapi.Response, error) {
	return es.Indices.Create(indexName)
}

func IndexingDocument[T any](es *elasticsearch.Client, indexName string, document []T) (*esapi.Response, error) {
	var buf bytes.Buffer
	for _, doc := range document {
		meta := map[string]interface{}{
			"index": map[string]interface{}{},
		}
		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return nil, err
		}
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return nil, err
		}
	}
	return es.Bulk(
		bytes.NewReader(buf.Bytes()),
		es.Bulk.WithIndex(indexName),
		es.Bulk.WithPretty())
}

func SearchingDocument(es *elasticsearch.Client, indexName string, query map[string]interface{}) (*esapi.Response, error) {
	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	return es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(bytes.NewReader(data)),
		es.Search.WithPretty())
}

func UpdateDocument(es *elasticsearch.Client, indexName string, id string, document interface{}) (*esapi.Response, error) {
	data := map[string]interface{}{
		"doc": document,
	}
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return es.Update(indexName, id, bytes.NewReader(d), es.Update.WithPretty())
}

func DeleteDocument(es *elasticsearch.Client, indexName string, id string) (*esapi.Response, error) {
	return es.Delete(indexName, id, es.Delete.WithPretty())
}

func main() {
	es := NewElasticSearchClient()

	// Create new index
	res, err := CreateNewIndex(es, "products")
	if err != nil {
		log.Printf("Error creating index: %s", err)
	} else {
		log.Println(res)
	}

	// Indexing document
	doc := []Product{
		Product{
			Name:        "Dell Laptop",
			Price:       420.0,
			Description: "Laptop product",
			Tags:        []string{"Electronic"},
		},
		Product{
			Name:        "Iphone",
			Price:       500.0,
			Description: "Iphone product",
			Tags:        []string{"Electronic", "Phone"},
		},
	}
	res, err = IndexingDocument(es, "products", doc)
	if err != nil {
		log.Printf("Error indexing document: %s", err)
	} else {
		log.Println(res)
	}

	// Query document
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}
	res, err = SearchingDocument(es, "products", query)
	if err != nil {
		log.Printf("Error searching documents: %s", err)
	} else {
		log.Println(res)
	}

	// Update document
	updatedDoc := Product{
		Price: 700,
	}
	res, err = UpdateDocument(es, "products", "Hb1VSJYBCAGLpkSMvwiX", updatedDoc)
	if err != nil {
		log.Printf("Error updating product: %s", err)
	} else {
		log.Println(res)
	}

	// Delete document
	res, err = DeleteDocument(es, "products", "Hb1VSJYBCAGLpkSMvwiX")
	if err != nil {
		log.Printf("Error deleting products: %s", err)
	} else {
		log.Println(res)
	}
}
