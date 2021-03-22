package services

import (
	"time"

	"github.com/alanwade2001/spa-common/rabbitmq"
	"github.com/alanwade2001/spa-messaging/spa-msg-initiation-instruction/types"
)

func NewMessaging(url string, queueName string, timeout time.Duration) types.MessagingAPI {
	return &rabbitmq.Messaging{
		Url:       url,
		QueueName: queueName,
		Timeout:   timeout,
	}
}
