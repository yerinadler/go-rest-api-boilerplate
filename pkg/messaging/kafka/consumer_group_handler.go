package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaConsumerGroupHandler struct {
	ready       chan bool
	handlerFunc func(*sarama.ConsumerMessage) error
}

func (cg *KafkaConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (cg *KafkaConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	close(cg.ready)
	return nil
}

func (cg *KafkaConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			if err := cg.handlerFunc(message); err != nil {
				return nil
			}

			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
