package integration

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"boilerplate/config"
	domain "boilerplate/internal/domain/event_stream"
	infrastructure "boilerplate/internal/infrastructure/event_stream"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/redpanda"
)

type TestEvent struct {
	ID      string    `json:"id"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

func setupRedpandaContainer(t *testing.T) (string, func()) {
	ctx := context.Background()

	container, err := redpanda.Run(ctx, "docker.redpanda.com/redpandadata/redpanda:v23.3.3")
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	brokerAddress, err := container.KafkaSeedBroker(ctx)
	if err != nil {
		t.Fatalf("failed to get broker address: %v", err)
	}

	cleanup := func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return brokerAddress, cleanup
}

func TestProducerConsumer(t *testing.T) {
	brokerAddress, cleanup := setupRedpandaContainer(t)
	defer cleanup()

	time.Sleep(2 * time.Second)

	config := config.Config{
		EventStream: infrastructure.Config{
			Brokers: []string{brokerAddress},
			Topic:   "test-topic",
		},
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.Level(config.LogLevel),
	}))

	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	adminConfig := sarama.NewConfig()
	adminConfig.Version = sarama.V2_0_0_0
	adminClient, err := sarama.NewClusterAdmin([]string{brokerAddress}, adminConfig)
	assert.NoError(t, err)
	defer adminClient.Close()

	err = adminClient.CreateTopic(
		config.EventStream.Topic,
		&sarama.TopicDetail{
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		false,
	)
	assert.NoError(t, err)

	producer, err := infrastructure.NewProducer(config.EventStream, logger)
	assert.NoError(t, err)
	defer producer.Close(ctx)

	consumer, err := infrastructure.NewConsumer(config.EventStream, logger)
	assert.NoError(t, err)
	defer consumer.Close(ctx)

	receivedEvents := make(chan domain.Event)
	consumeErrors := make(chan error)

	go func() {
		err := consumer.Consume(ctx, func(ctx context.Context, event domain.Event) error {
			receivedEvents <- event
			return nil
		})
		if err != nil {
			consumeErrors <- err
		}
	}()

	testEvent := TestEvent{
		ID:      "123",
		Message: "test message",
		Time:    time.Now(),
	}
	payload, err := json.Marshal(testEvent)
	assert.NoError(t, err)

	event := domain.Event{
		Type:    "TestEvent",
		Payload: payload,
		Metadata: map[string]string{
			"version": "1.0",
			"test":    "true",
		},
	}

	err = producer.Produce(ctx, event)
	assert.NoError(t, err)

	select {
	case receivedEvent := <-receivedEvents:
		assert.Equal(t, event.Type, receivedEvent.Type)
		assert.Equal(t, event.Payload, receivedEvent.Payload)
		assert.Equal(t, event.Metadata, receivedEvent.Metadata)

		var receivedTestEvent TestEvent
		err = json.Unmarshal(receivedEvent.Payload, &receivedTestEvent)
		assert.NoError(t, err)
		assert.Equal(t, testEvent.ID, receivedTestEvent.ID)
		assert.Equal(t, testEvent.Message, receivedTestEvent.Message)

	case err := <-consumeErrors:
		t.Fatalf("Error consuming message: %v", err)

	case <-ctx.Done():
		t.Fatal("Timeout waiting for message")
	}
}
