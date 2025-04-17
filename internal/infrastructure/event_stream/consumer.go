package eventstream

import (
	"context"
	"fmt"
	"sync"

	eventstream "boilerplate/internal/domain/event_stream"
	"boilerplate/internal/domain/shared"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	config   sarama.Config
	brokers  []string
	topic    string
	logger   shared.Logger
}

func NewConsumer(config Config, logger shared.Logger) (Consumer, error) {
	saramaConfig := sarama.NewConfig()

	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	saramaConfig.Net.TLS.Enable = config.TLSEnabled
	saramaConfig.Net.SASL.Enable = config.SASEnabled
	saramaConfig.Version = sarama.V2_0_0_0

	consumer, err := sarama.NewConsumer(config.Brokers, saramaConfig)
	if err != nil {
		return Consumer{}, fmt.Errorf("failed to create consumer: %w", err)
	}

	return Consumer{
		consumer: consumer,
		config:   *saramaConfig,
		brokers:  config.Brokers,
		topic:    config.Topic,
		logger:   logger,
	}, nil
}

func (c Consumer) Consume(ctx context.Context, handler eventstream.EventHandler) error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	var wg sync.WaitGroup
	done := make(chan struct{})

	for _, partition := range partitions {
		consumePartition, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return fmt.Errorf("failed to start consumer for partition %d: %w", partition, err)
		}

		wg.Add(1)
		go func(consumePartition sarama.PartitionConsumer) {
			defer wg.Done()
			defer consumePartition.AsyncClose()

			for {
				select {

				case msg := <-consumePartition.Messages():
					event := eventstream.Event{
						Type:     string(msg.Key),
						Payload:  msg.Value,
						Metadata: make(map[string]string),
					}

					for _, header := range msg.Headers {
						event.Metadata[string(header.Key)] = string(header.Value)
					}

					if err := handler(ctx, event); err != nil {
						c.logger.ErrorContext(ctx, "Error handling event", "error", err)

						continue
					}

				case err := <-consumePartition.Errors():
					c.logger.ErrorContext(ctx, "Error consuming message", "error", err)

					continue

				case <-ctx.Done():
					return

				case <-done:
					return

				}
			}
		}(consumePartition)
	}

	go func() {
		<-ctx.Done()
		close(done)
	}()

	wg.Wait()

	return nil
}

func (c Consumer) Close(ctx context.Context) error {
	if err := c.consumer.Close(); err != nil {
		c.logger.ErrorContext(ctx, "Failed to close consumer", "error", err)

		return fmt.Errorf("failed to close consumer: %w", err)
	}

	return nil
}
