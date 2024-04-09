package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	logger   *logrus.Logger
	tracer   trace.Tracer
}

func NewKafkaProducer(brokers []string, logger *logrus.Logger, tracer trace.Tracer) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	wrapped := otelsarama.WrapSyncProducer(config, producer)
	if err != nil {
		return nil, err
	}

	kafkaProducer := &KafkaProducer{wrapped, logger, tracer}

	return kafkaProducer, nil
}

func (p *KafkaProducer) Publish(ctx context.Context, topic string, value string, key string) error {
	_, span := p.tracer.Start(ctx, "publish to kafka")
	defer span.End()

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
		Key:   sarama.StringEncoder(key),
	}

	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(producerMessage))

	partition, offset, err := p.producer.SendMessage(producerMessage)
	if err != nil {
		return err
	}

	p.logger.WithContext(ctx).WithFields(logrus.Fields{
		"topic":     topic,
		"value":     value,
		"key":       key,
		"partition": partition,
		"offset":    offset,
	}).Info("published to the Kafka topic")

	return nil
}
