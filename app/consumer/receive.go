package main

import (
	"github.com/streadway/amqp"
	"github.com/Traineau/gomail/helpers"
	"log"
	"github.com/Traineau/gomail/email"
	"encoding/json"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")

	helpers.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	helpers.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	helpers.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	helpers.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var campaignID email.CampaignID
			err := json.Unmarshal(d.Body, &campaignID)
			if err != nil {
				log.Printf("ERROR %+v", err)
			}

			log.Printf("Received a message: %+v", campaignID)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
