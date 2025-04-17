package eventstream

import (
	"context"
	"fmt"

	eventstream "boilerplate/internal/domain/event_stream"
	"boilerplate/internal/domain/shared"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
	logger   shared.Logger
}

func NewProducer(config Config, logger shared.Logger) (Producer, error) {
	saramaConfig := sarama.NewConfig()

	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	saramaConfig.Net.TLS.Enable = config.TLSEnabled
	saramaConfig.Net.SASL.Enable = config.SASEnabled
	saramaConfig.Version = sarama.V2_0_0_0

	producer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)
	if err != nil {
		return Producer{}, fmt.Errorf("failed to create producer: %w", err)
	}

	return Producer{
		producer: producer,
		topic:    config.Topic,
		logger:   logger,
	}, nil
}

func (p Producer) Produce(ctx context.Context, event eventstream.Event) error {
	headers := make([]sarama.RecordHeader, 0, len(event.Metadata))
	for k, v := range event.Metadata {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	msg := &sarama.ProducerMessage{
		Topic:   p.topic,
		Key:     sarama.StringEncoder(event.Type),
		Value:   sarama.ByteEncoder(event.Payload),
		Headers: headers,
	}

	select {
	case <-ctx.Done():
		return ctx.Err()

	default:
		_, _, err := p.producer.SendMessage(msg)
		if err != nil {
			p.logger.ErrorContext(ctx, "Failed to send message", "error", err)

			return fmt.Errorf("failed to send message: %w", err)
		}

		return nil
	}
}

func (p Producer) Close(ctx context.Context) error {
	if err := p.producer.Close(); err != nil {
		p.logger.ErrorContext(ctx, "Failed to close producer", "error", err)

		return fmt.Errorf("failed to close producer: %w", err)
	}

	return nil
}
