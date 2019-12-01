package dispatcher

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

func elasticSearchDispatcher(addr, entry string) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(addr))
	if err != nil {
		panic(err)
	}
	r, err := client.IndexExists("logs").Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !r {
		_, err := client.CreateIndex("logs").Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	indexService := elastic.NewIndexService(client)
	jsonMessage := fmt.Sprintf("{\"logs\": \"%s\"}", entry)
	indexService.Index("logs").BodyJson(jsonMessage)
	r2, err := indexService.Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", r2)
}
