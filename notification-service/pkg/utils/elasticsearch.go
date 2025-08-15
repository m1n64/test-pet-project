package utils

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
	"time"
)

type ElasticClient struct {
	Client *elastic.Client
	Index  string
}

func NewElasticClient(url string, index string) *ElasticClient {
	for i := 0; i < 10; i++ {
		client, err := elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetSniff(false),
		)
		if err == nil {
			_, _, pingErr := client.Ping(url).Do(context.Background())
			if pingErr == nil {
				return &ElasticClient{
					Client: client,
					Index:  index,
				}
			}
		}

		log.Println("Elastic search is not ready yet")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Elastic search is not ready")
	return nil
}
