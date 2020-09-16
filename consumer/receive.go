package consumer

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Traineau/gomail/email"
	"github.com/Traineau/gomail/helpers"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var (
	DbConn *sql.DB
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	err = dbConnect()
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

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
			var campaign email.Campaign
			err := json.Unmarshal(d.Body, &campaign)
			if err != nil {
				log.Printf("ERROR %+v", err)
			}

			log.Printf("Received a message: %+v", campaign.ID)

			repository := email.Repository{Conn: DbConn}

			campaignFromRepo, err := repository.GetCampaign(campaign.ID)
			if err != nil {
				log.Printf("ERROR %+v", err)
			}

			var mailingList *email.MailingList
			if campaignFromRepo != nil {
				mailingList, err = repository.GetMailingList(campaignFromRepo.IdMailingList)
				if err != nil {
					log.Printf("ERROR %+v", err)
				}

				if mailingList != nil {
					recipients, err := repository.GetRecipientsFromMailingList(mailingList.ID)
					if err != nil {
						log.Printf("ERROR %+v", err)
					}
					if recipients != nil {
						mailingList.Recipients = recipients
					}
				}
			}

			log.Printf("campaign: %+v", campaignFromRepo)
			log.Printf("mailing list: %+v", mailingList)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}

func dbConnect() error {

	dsn := "gomail:gomail@tcp(localhost:3306)/image_gomail?parseTime=true&charset=utf8"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return err
	}

	var dbErr error
	for i := 1; i <= 8; i++ {
		dbErr = db.Ping()
		if dbErr != nil {
			if i < 8 {
				log.Printf("db connection failed, %d retry : %v", i, dbErr)
				time.Sleep(10 * time.Second)
			}
			continue
		}

		break
	}

	if dbErr != nil {
		return errors.New("can't connect to database after 3 attempts")
	}

	DbConn = db

	return nil
}
