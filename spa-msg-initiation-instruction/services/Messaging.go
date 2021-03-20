package services

import (
	"github.com/alanwade2001/spa-messaging/spa-msg-initiation-instruction/types"
	"github.com/streadway/amqp"
)

type Messaging struct {
	url       string
	queueName string

	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewMessaging(url string, queueName string) types.MessagingAPI {
	return &Messaging{url: url, queueName: queueName}
}

func (m *Messaging) Connect() (err error) {
	if m.conn, err = amqp.Dial(m.url); err != nil {
		return err
	}

	if m.ch, err = m.conn.Channel(); err != nil {
		return err
	}

	if m.q, err = m.ch.QueueDeclare(
		"SPA-INITIATION-INSTRUCTION", // name
		false,                        // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	); err != nil {
		m.Disconnect()
		return err
	}

	return nil
}

func (m *Messaging) Consume() (<-chan amqp.Delivery, error) {
	return m.ch.Consume(
		m.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
}

func (m *Messaging) Disconnect() error {
	m.ch.Close()
	m.conn.Close()

	return nil
}
