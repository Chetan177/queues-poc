package main

// https://github.com/aws-samples/aws-msk-content-streaming

import (
	"context"
	"math/rand"
	"net"
	"strconv"

	"fmt"

	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

var BrokerURLs = []string{
	
}

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	kafkaConfig := kafka.WriterConfig{
		Brokers:  BrokerURLs,
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
	return kafka.NewWriter(kafkaConfig)
}

func createKafkaTopic(kafkaURL, topic string) {
	conn, err := kafka.Dial("tcp", kafkaURL)
	if err != nil {
		panic(err.Error())
	}
	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 2,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

}

func main() {

	kafkaURL := BrokerURLs[0]
	topic := "topicTest"

	createKafkaTopic(kafkaURL, topic)

	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()
	fmt.Println("start producing ... !!")

	for i := 0; ; i++ {
		keyval := rand.Intn(3)
		key := fmt.Sprintf("Key-%d", keyval)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("produced", key)
		}
		time.Sleep(1 * time.Second)
	}
}
