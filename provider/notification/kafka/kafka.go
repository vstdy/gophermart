package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"

	canonical "github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/provider/notification"
)

var _ notification.Notification = (*Kafka)(nil)

type (
	// Kafka keeps kafka dependencies.
	Kafka struct {
		config              Config
		sConf               *sarama.Config
		syncProducer        sarama.SyncProducer
		asyncProducer       sarama.AsyncProducer
		accNtfConsumerGroup sarama.ConsumerGroup
		accNtfConsumer      *AccrualNtfConsumer
	}

	// KafkaOption defines functional argument for Service constructor.
	KafkaOption func(*Kafka) error
)

// WithConfig sets Config.
func WithConfig(config Config) KafkaOption {
	return func(kfk *Kafka) error {
		kfk.config = config

		return nil
	}
}

// NewKafkaProvider returns a new Kafka instance.
func NewKafkaProvider(opts ...KafkaOption) (*Kafka, error) {
	kafka := &Kafka{
		config: NewDefaultConfig(),
		sConf:  sarama.NewConfig(),
		accNtfConsumer: &AccrualNtfConsumer{
			transactions: make(chan canonical.Transaction),
			ready:        make(chan struct{}),
		},
	}
	for optIdx, opt := range opts {
		if err := opt(kafka); err != nil {
			return nil, fmt.Errorf("applying option [%d]: %w", optIdx, err)
		}
	}

	if err := kafka.config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	if err := newProducer(kafka); err != nil {
		return nil, fmt.Errorf("newProducer: %w", err)
	}

	if err := newAccNtfConsumerGroup(kafka); err != nil {
		return nil, fmt.Errorf("newAccNtfConsumerGroup: %w", err)
	}

	return kafka, nil
}

// Close closes all Kafka dependencies.
func (kfk *Kafka) Close() error {
	if err := kfk.asyncProducer.Close(); err != nil {
		return fmt.Errorf("async producer: %w", err)
	}
	if err := kfk.syncProducer.Close(); err != nil {
		return fmt.Errorf("sync producer: %w", err)
	}
	if err := kfk.accNtfConsumerGroup.Close(); err != nil {
		return fmt.Errorf("consumer group: %w", err)
	}
	close(kfk.accNtfConsumer.transactions)

	return nil
}

// newProducer sets up producers.
func newProducer(kafka *Kafka) error {
	kafka.sConf.Producer.Partitioner = sarama.NewRandomPartitioner
	kafka.sConf.Producer.RequiredAcks = sarama.WaitForAll
	kafka.sConf.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(kafka.config.BrokersAddresses, kafka.sConf)
	if err != nil {
		return fmt.Errorf("creating sync producer: %w", err)
	}
	asyncProducer, err := sarama.NewAsyncProducer(kafka.config.BrokersAddresses, kafka.sConf)
	if err != nil {
		return fmt.Errorf("creating async producer: %w", err)
	}

	go func() {
		for err := range asyncProducer.Errors() {
			log.Warn().Err(err).Msg("Msg async:")
		}
	}()

	go func() {
		for succ := range asyncProducer.Successes() {
			log.Info().Msgf("Msg written async. Partition: %v. Offset: %v",
				succ.Partition, succ.Offset)
		}
	}()

	kafka.syncProducer = syncProducer
	kafka.asyncProducer = asyncProducer

	return nil
}

// newAccNtfConsumerGroup adds an accrual notifications consumer group.
func newAccNtfConsumerGroup(kafka *Kafka) error {
	kafka.sConf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	kafka.sConf.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(
		kafka.config.BrokersAddresses, kafka.config.AccrualsConsumerGroupName, kafka.sConf)
	if err != nil {
		return err
	}

	kafka.accNtfConsumerGroup = consumerGroup

	return nil
}

// accNtfTopics returns a list of all accrual notifications topics.
func (kfk Kafka) accNtfTopics() []string {
	return []string{kfk.config.AccrualsTopicName}
}
