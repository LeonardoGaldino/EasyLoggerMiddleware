package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LeonardoGaldino/EasyLoggerMiddleware/src/internal/utils"

	"github.com/olivere/elastic/v7"
)

func elasticSearchDispatcher(addr, entry string) {
	client := utils.KeepRetryingAfter(func() (interface{}, error) {
		return elastic.NewSimpleClient(elastic.SetURL(addr))

	}, time.Second).(*elastic.Client)

	indexExists := utils.KeepRetryingAfter(func() (interface{}, error) {
		return client.IndexExists("logs").Do(context.Background())

	}, time.Second).(bool)

	if !indexExists {
		utils.KeepRetryingAfter(func() (interface{}, error) {
			return client.CreateIndex("logs").Do(context.Background())
		}, time.Second)
	}

	result := utils.KeepRetryingAfter(func() (interface{}, error) {
		indexService := elastic.NewIndexService(client)
		jsonMessage := fmt.Sprintf("{\"logs\": \"%s\"}", entry)
		indexService.Index("logs").BodyJson(jsonMessage)
		res, err := indexService.Do(context.Background())
		if err != nil {
			return nil, err
		}
		if res.Result != "created" {
			return nil, errors.New("Not created")
		}
		return res, nil
	}, time.Second).(*elastic.IndexResponse)
	fmt.Printf("%+v\n", result)
}
