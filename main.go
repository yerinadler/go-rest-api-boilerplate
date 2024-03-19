package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

type IEventHandler interface{}

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

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPause *bool) {
	if *isPause {
		client.ResumeAll()
	} else {
		client.PauseAll()
	}

	*isPause = !*isPause
}

func StartConsumption(
	brokers []string,
	topics []string,
	clientId string,
	groupId string,
	handlerFunc func(*sarama.ConsumerMessage) error,
) error {
	keepRunning := true

	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	consumer := KafkaConsumerGroupHandler{
		ready:       make(chan bool),
		handlerFunc: handlerFunc,
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, groupId, config)

	if err != nil {
		log.Panicf("error creating consumer client : %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Panicf("error initiating consumption : %v", err)
			}

			if ctx.Err() != nil {
				return
			}

			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			keepRunning = false
		case <-sigterm:
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}

	cancel()
	wg.Wait()

	if err = client.Close(); err != nil {
		log.Panicf("error closing the client : %v", err)
	}

	return nil
}
