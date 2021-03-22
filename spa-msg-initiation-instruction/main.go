package main

import (
	"log"

	"github.com/alanwade2001/spa-messaging/spa-msg-initiation-instruction/services"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {

	services.NewConfigService().Load(".")

	var err error
	var msgs <-chan amqp.Delivery

	messagingUrl := viper.GetString("MESSAGE_SERVICE_URI")
	queueName := viper.GetString("INITIATION_QUEUE")
	timeout := viper.GetDuration("MESSAGING_TIMEOUT")

	messagingService := services.NewMessaging(messagingUrl, queueName, timeout)

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
