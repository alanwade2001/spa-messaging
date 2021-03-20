package main

import (
	"log"

	"github.com/alanwade2001/spa-messaging/spa-msg-initiation-instruction/services"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	// ch, err := conn.Channel()
	// failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	// q, err := ch.QueueDeclare(
	// 	"SPA-INITIATION-INSTRUCTION", // name
	// 	false,                        // durable
	// 	false,                        // delete when unused
	// 	false,                        // exclusive
	// 	false,                        // no-wait
	// 	nil,                          // arguments
	// )
	// failOnError(err, "Failed to declare a queue")

	// msgs, err := ch.Consume(
	// 	q.Name, // queue
	// 	"",     // consumer
	// 	true,   // auto-ack
	// 	false,  // exclusive
	// 	false,  // no-local
	// 	false,  // no-wait
	// 	nil,    // args
	// )
	// failOnError(err, "Failed to register a consumer")
	var err error
	var msgs <-chan amqp.Delivery
	messagingService := services.NewMessaging("amqp://guest:guest@localhost:5672/", "SPA-INITIATION-INSTRUCTION")

	if err := messagingService.Connect(); err != nil {
		failOnError(err, "failed to connect")
	}
	defer messagingService.Disconnect()
	msgs, err = messagingService.Consume()
	failOnError(err, "failed to consume messages")

	forever := make(chan bool)

	go func() {
		messageService := services.NewMessage()
		for d := range msgs {
			if err := messageService.Process(d.Body); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
