package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

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
	config.Version = sarama.V2_5_0_0
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
	log.Println("Sarama consumer up and running!...")

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
