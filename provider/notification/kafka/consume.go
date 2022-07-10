package kafka

import (
	"context"

	"github.com/rs/zerolog/log"

	canonical "github.com/vstdy/gophermart/model"
)

// ConsumeAccrualNotifications receives accrual notifications.
func (kfk *Kafka) ConsumeAccrualNotifications(ctx context.Context) {
	for {
		if err := kfk.accNtfConsumerGroup.Consume(ctx, kfk.accNtfTopics(), kfk.accNtfConsumer); err != nil {
			log.Warn().Err(err).Msg("consumer:")
		}

		if ctx.Err() != nil {
			return
		}

		kfk.accNtfConsumer.ready = make(chan struct{})
	}
}

// GetAccrualNotificationsChan returns the accrual notifications consumer channel.
func (kfk *Kafka) GetAccrualNotificationsChan() chan canonical.Transaction {
	return kfk.accNtfConsumer.transactions
}
