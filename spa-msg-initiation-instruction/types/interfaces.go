package types

import "github.com/streadway/amqp"

type MessagingAPI interface {
	Connect() error
	Disconnect() error
	Consume() (<-chan amqp.Delivery, error)
}

type MessageAPI interface {
	Process(body []byte) error
}

// ConfigAPI si
type ConfigAPI interface {
	Load(configPath string) error
}
