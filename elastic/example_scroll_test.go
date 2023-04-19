package elastic

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/bytedance/sonic"
	"github.com/elastic/go-elasticsearch/v7"
)

type ModelMsg struct {
	ScrollID string `json:"_scroll_id"`
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Hits     struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore any `json:"max_score"`
		Hits     []struct {
			Index   string   `json:"_index"`
			Type    string   `json:"_type"`
			ID      string   `json:"_id"`
			Score   any      `json:"_score"`
			Ignored []string `json:"_ignored"`
			Source  struct {
				Qos          int    `json:"qos"`
				Level        string `json:"level"`
				Env          string `json:"env"`
				Time         string `json:"time"`
				Version      string `json:"@version"`
				Timestamp    string `json:"@timestamp"`
				Module       string `json:"module"`
				Message      string `json:"message"`
				Source       string `json:"source"`
				MqttTag      string `json:"mqtt_tag"`
				Type         string `json:"type"`
				Caller       string `json:"caller"`
				ServerIP     string `json:"server_ip"`
				Version0     string `json:"version"`
				ReceiveTopic string `json:"receive_topic"`
			} `json:"_source"`
			Sort []int `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

var query = `
{
    "query": {
        "bool": {
            "filter": [
                {
                    "bool": {
                        "filter": [
                            {
                                "multi_match": {
                                    "lenient": true,
                                    "query": "ReceiverMessage",
                                    "type": "best_fields"
                                }
                            },
                            {
                                "multi_match": {
                                    "lenient": true,
                                    "query": "300019",
                                    "type": "best_fields"
                                }
                            }
                        ]
                    }
                }
            ]
        }
    }
}
`

var queryWithUuid = `
{
    "query": {
        "bool": {
            "filter": [
                {
                    "bool": {
                        "filter": [
                            {
                                "multi_match": {
                                    "lenient": true,
                                    "query": "ReceiverMessage",
                                    "type": "best_fields"
                                }
                            },
                            {
                                "multi_match": {
                                    "lenient": true,
                                    "query": "300019",
                                    "type": "best_fields"
                                }
                            },
                            {
                                "multi_match": {
                                    "lenient": true,
                                    "query": "3301000000189651",
                                    "type": "best_fields"
                                }
                            }
                        ]
                    }
                }
            ]
        }
    }
}
`

type IotMessage struct {
	MsgId      string      `json:"msg_id,omitempty"`
	Code       interface{} `json:"code,omitempty"`
	Sync       string      `json:"sync,omitempty"`
	Command    string      `json:"command,omitempty"`
	DeviceId   string      `json:"device_id,omitempty"`
	Topic      string      `json:"topic,omitempty"`
	Timestamp  string      `json:"timestamp,omitempty"`
	Type       string      `json:"type,omitempty"`
	DeviceType string      `json:"device_type,omitempty"`
	EventStart string      `json:"event_start,omitempty"`
	AuthKey    string      `json:"auth_key,omitempty"`
	Times      string      `json:"times,omitempty"`
	EventTime  string      `json:"event_time,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	MqttTag    string      `json:"mqtt_tag,omitempty"`
}

func TestElasticScroll(t *testing.T) {
	log.SetFlags(0)

	var (
		batchNum int
		scrollID string
	)

	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"},
		Username:  "elastic",
		Password:  "123456",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}

	// Perform the initial search request to get
	// the first batch of data and the scroll ID
	//
	log.Println("开始查询...")
	log.Printf("第 %d 批查询数据\n", batchNum)
	res, _ := es.Search(
		es.Search.WithBody(strings.NewReader(queryWithUuid)),
		es.Search.WithIndex("test-2023.12.29"),
		es.Search.WithSort("@timestamp"),
		es.Search.WithSize(3),
		es.Search.WithScroll(time.Minute),
	)

	// Handle the first batch of data and extract the scrollID
	//
	body := read(res.Body)
	err = res.Body.Close()
	if err != nil {
		log.Fatalf("===> Error: %s", err)
	}

	var msg ModelMsg
	err = sonic.Unmarshal([]byte(body), &msg)
	if err != nil {
		log.Fatalf("===> Error: %s", err)
	}

	scrollID = msg.ScrollID

	if len(msg.Hits.Hits) == 0 {
		log.Println("完成游标遍历...")
		return
	}

	for _, hit := range msg.Hits.Hits {
		var iotMsg IotMessage
		err = sonic.Unmarshal([]byte(hit.Source.Message), &iotMsg)
		if err != nil {
			log.Fatalf("===> Error: %s", err)
		}

		fmt.Printf("======> %+v\n", iotMsg)
	}

	log.Println()
	log.Println()

	// Perform the scroll requests in sequence

	for {
		log.Println(strings.Repeat("-", 80))
		batchNum++

		log.Printf("第 %d 批查询数据\n", batchNum)

		// Perform the scroll request and pass the scrollID and scroll duration
		//
		res, err := es.Scroll(es.Scroll.WithScrollID(scrollID), es.Scroll.WithScroll(time.Minute))
		if err != nil {
			log.Fatalf("===> Error: %s", err)
		}
		if res.IsError() {
			log.Fatalf("===> Error response: %s", res)
		}

		data := read(res.Body)
		err = res.Body.Close()
		if err != nil {
			log.Fatalf("===> Error: %s", err)
		}

		var msg ModelMsg
		err = sonic.Unmarshal([]byte(data), &msg)
		if err != nil {
			log.Fatalf("===> Error: %s", err)
		}

		// Extract the scrollID from response
		scrollID = msg.ScrollID

		// Break out of the loop when there are no results
		if len(msg.Hits.Hits) == 0 {
			log.Println("完成游标遍历...")
			break
		}

		for _, hit := range msg.Hits.Hits {
			var iotMsg IotMessage
			err = sonic.Unmarshal([]byte(hit.Source.Message), &iotMsg)
			if err != nil {
				log.Fatalf("===> Error: %s", err)
			}

			fmt.Printf("======> %+v\n", iotMsg)
		}

		log.Println()
	}
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}
