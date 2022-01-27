package main

import (
	"context"
	"fmt"
	"log"

	"strings"

	kafka "github.com/segmentio/kafka-go"
)

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1, // 1B
		MaxBytes: 10e6, // 10MB
	})
}

func main() {

	kafkaURL := "localhost:9092"
	topic := "topicTest-3part"
	groupID := "testConsumer"

	reader := getKafkaReader(kafkaURL, topic, groupID)

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