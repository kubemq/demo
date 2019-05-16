package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

const mapping = `
{
	"settings": {
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings": {
		"properties": {
			"client_id": {
				"type": "keyword"
			},
			"metadata": {
				"type": "keyword"
			},
			"body": {
				"type": "text"
			}
		}
	}
}`

type Elastic struct {
	client *elastic.Client
}

func NewElasticSearch(url string) (*Elastic, error) {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	info, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	exists, err := client.IndexExists("logs").Do(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		resultDelete, err := client.DeleteIndex("logs").Do(ctx)
		if err != nil {
			return nil, err
		}
		if !resultDelete.Acknowledged {
			return nil, errors.New("index logs was not deleted")
		}

	}

	createIndex, err := client.CreateIndex("logs").BodyString(mapping).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
	}
	log.Println("Elasticsearch created logs index")

	el := &Elastic{
		client: client,
	}
	return el, nil
}

func (el *Elastic) Save(ctx context.Context, msg *Log) error {
	log.Printf("Event Id: %s recevied, saving to elastic.\n", msg.ID)
	_, err := el.client.Index().
		Index("logs").
		Id(msg.ID).
		BodyJson(msg).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
