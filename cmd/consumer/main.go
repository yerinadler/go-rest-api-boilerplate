package main

import (
	"log"

	"github.com/IBM/sarama"
	"github.com/example/go-rest-api-revision/config"
	"github.com/example/go-rest-api-revision/pkg/messaging/kafka"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := kafka.StartConsumption(
		cfg.Kafka.Brokers,
		[]string{"test"},
		"go-rest-api-consumer",
		"go-rest-api",
		func(message *sarama.ConsumerMessage) error {
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			return nil
		},
	); err != nil {
		log.Fatal(err)
	}
}
