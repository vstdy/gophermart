package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"

	canonical "github.com/vstdy/gophermart/model"
	"github.com/vstdy/gophermart/provider/notification/kafka/model"
)

// ProduceAccrualNotifications sends notifications of accruals.
func (kfk *Kafka) ProduceAccrualNotifications(objs []canonical.Transaction) error {
	notifications := model.NewAccrualNotificationsFromCanonical(objs)
	var msgs []*sarama.ProducerMessage

	for _, notification := range notifications {
		marshal, err := json.Marshal(notification)
		if err != nil {
			return fmt.Errorf("marshaling accrual notification: %w", err)
		}

		msgs = append(msgs, &sarama.ProducerMessage{
			Topic: kfk.config.AccrualsTopicName,
			Value: sarama.ByteEncoder(marshal),
		})
	}

	return kfk.syncProducer.SendMessages(msgs)
}
