package kafka

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"

	canonical "github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/provider/notification/kafka/model"
)

// AccrualNtfConsumer represents an accrual notifications consumer group consumer.
type AccrualNtfConsumer struct {
	transactions chan canonical.Transaction
	ready        chan struct{}
}

// Setup performs actions at the beginning of a new session, before ConsumeClaim
func (consumer *AccrualNtfConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

// Cleanup performs actions at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *AccrualNtfConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim starts a consumer loop of ConsumerGroupClaim's Messages.
func (consumer *AccrualNtfConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var msg model.AccrualNotification
		if err := json.Unmarshal(message.Value, &msg); err != nil {
			log.Warn().Err(err).Msg("unmarshaling message")
			continue
		}
		consumer.transactions <- msg.ToCanonical()
	}

	return nil
}
