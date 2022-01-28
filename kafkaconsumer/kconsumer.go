package main

import (
	"context"
	"fmt"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

var BrokerURLs = []string{
	
}

func getKafkaReader(topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  BrokerURLs,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1, // 1B
		MaxBytes: 10e6, // 10MB
	})
}

func main() {

	topic := "topicTest"
	groupID := "testConsumer"

	reader := getKafkaReader(topic, groupID)

	defer reader.Close()

	fmt.Println("start consuming ... !!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("message at partition:%v offset:%v	%s \n", m.Partition, m.Offset, string(m.Key))
	}
}