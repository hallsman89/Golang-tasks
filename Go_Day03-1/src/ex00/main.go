 package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var (
	res             *esapi.Response
	countSuccessful uint64
	err             error
)

const (
	csvFileName     = "../../materials/data.csv"
	indexName       = "places"
	mappingFileName = "schema.json"
)

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal("Can't create client", err)
	}
	mapping, err := readMappingFromFile()
	if err != nil {
		log.Fatal("Can't read file", err)
	}
	createIndexMapping(es, mapping)
	data, err := parseCsvFile()
	if err != nil {
		log.Fatal("Can't parse csv file", err)
	}
	loadDataIntoElastic(es, data)
}

func readMappingFromFile() (string, error) {
	res, err := os.ReadFile(mappingFileName)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func createIndexMapping(es *elasticsearch.Client, mapping string) {
	if res, err = es.Indices.Delete([]string{indexName},
		es.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
		log.Fatal("Can't delete index", err) // удаляем если уже существует такой индекс
	}
	defer func() { _ = res.Body.Close() }()
	res, err := es.Indices.Create(indexName, es.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		log.Fatal("Can't create index", err)
	}
	if res.IsError() {
		log.Fatal("Can't create index", res)
	}
	defer func() { _ = res.Body.Close() }()
}

func loadDataIntoElastic(es *elasticsearch.Client, data []Data) {
	bi, err := esutil.NewBulkIndexer(
		esutil.BulkIndexerConfig{
			Index:         indexName,
			Client:        es,
			NumWorkers:    8,
			FlushBytes:    10000,
			FlushInterval: 30 * time.Second,
		})
	if err != nil {
		log.Fatal("Error creating the indexer", err)
	}
	for _, person := range data {
		personInfo, err := json.Marshal(person)
		if err != nil {
			log.Fatalf("Cannot encode article %s: %s", person.ID, err)
		}
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: person.ID,
				Body:       bytes.NewReader(personInfo),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Println("ERROR:", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatal("Unexpected error:", err)
		}
	}
	if err := bi.Close(context.Background()); err != nil {
		log.Fatal("Unexpected error:", err)
	}
}
