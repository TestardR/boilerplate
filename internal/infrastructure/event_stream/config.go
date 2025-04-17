package eventstream

type Config struct {
	Brokers    []string `envconfig:"EVENT_STREAM_BROKERS" default:"localhost:9092"`
	Topic      string   `envconfig:"EVENT_STREAM_TOPIC" default:"test-topic"`
	TLSEnabled bool     `envconfig:"EVENT_STREAM_TLS_ENABLED" default:"false"`
	SASEnabled bool     `envconfig:"EVENT_STREAM_SASL_ENABLED" default:"false"`
}
