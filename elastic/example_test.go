package elastic

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func TestElasticSearch(t *testing.T) {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"},
		Username:  "elastic",
		Password:  "123456",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	var (
		index = "test-*"
		// sort    = `{"sort": [{"timestamp": {"order": "aes"}}]}`
	)

	size := 100

	// 发送搜索请求
	res, err := esapi.SearchRequest{
		Index:  []string{index},
		Body:   strings.NewReader(query),
		Pretty: true,
		Size:   &size,
		Sort:   []string{`{"@timestamp":{"order":"aes","unmapped_type":"boolean"}}`},
	}.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error sending the search request: %s", err)
	}
	defer res.Body.Close()

	// 处理搜索结果
	if res.IsError() {
		log.Fatalf("Error response: %s", res.String())
	}

	body := strings.TrimLeft(res.String(), "[200 OK] ")

	var msg ModelMsg
	err = sonic.Unmarshal([]byte(body), &msg)
	if err != nil {
		log.Fatalf("msg -> json.Unmarshal: %s, body:%s\n", err, res.String())
	}

	for _, v := range msg.Hits.Hits {
		fmt.Printf("====> ts: %s, msg: %s\n", v.Source.Timestamp, v.Source.Message)
	}
}
