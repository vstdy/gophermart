package kafka

import (
	"fmt"
)

// Config keeps Kafka params.
type Config struct {
	BrokersAddresses          []string `mapstructure:"kafka_brokers_addresses"`
	AccrualsTopicName         string   `mapstructure:"kafka_accruals_topic_name"`
	AccrualsConsumerGroupName string   `mapstructure:"kafka_accruals_consumer_group_name"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if len(config.BrokersAddresses) == 0 {
		return fmt.Errorf("kafka_broker_address field: empty")
	}

	if config.AccrualsTopicName == "" {
		return fmt.Errorf("kafka_accruals_topic_name field: empty")
	}

	if config.AccrualsConsumerGroupName == "" {
		return fmt.Errorf("kafka_accruals_consumer_group_name field: empty")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		BrokersAddresses:          []string{"127.0.0.1:9093", "127.0.0.1:9094", "127.0.0.1:9095"},
		AccrualsTopicName:         "accruals",
		AccrualsConsumerGroupName: "accruals_ntf",
	}
}
