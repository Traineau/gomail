package email

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"gomail/database"
	"gomail/helpers"
	"log"
	"net/http"
)

func CreateMailingList(w http.ResponseWriter, r *http.Request) {
	var mailingList MailingList
	err := json.NewDecoder(r.Body).Decode(&mailingList)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not decode request body")
		return
	}

	db := database.DbConn
	repository := Repository{Conn: db}

	err = repository.SaveMailingList(&mailingList)
	if err != nil {
		log.Printf("could not save mailing list: %v", err)
		return
	}

	log.Printf("mailing list : %+v", mailingList)

	helpers.WriteJSON(w, http.StatusOK, mailingList)
}

func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var campaign Campaign
	err := json.NewDecoder(r.Body).Decode(&campaign)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not decode request body")
		return
	}

	db := database.DbConn
	repository := Repository{Conn: db}

	err = repository.SaveCampaign(&campaign)
	if err != nil {
		log.Printf("could not save campaign: %v", err)
		return
	}

	log.Printf("campaign : %+v", campaign)
	helpers.WriteJSON(w, http.StatusOK, campaign)
}

func AddRecipientToMailinglist(w http.ResponseWriter, r *http.Request) {
	var recipients []*Recipient
	err := json.NewDecoder(r.Body).Decode(&recipients)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not decode request body")
		return
	}

	db := database.DbConn
	repository := Repository{Conn: db}

	ids, err := repository.AddRecipients(recipients)
	if err != nil {
		log.Printf("could not save recipients list: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not save recipients")
		return
	}

	muxVars := mux.Vars(r)
	mailingListID, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		log.Printf("could not save parse id into int: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not parse id")
		return
	}

	err = repository.AddRecipientToMailingList(ids, mailingListID)
	if err != nil {
		log.Printf("could not save recipients list: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not add recipients to mailing list")
		return
	}

	mailingList, err := repository.GetMailingList(mailingListID)
	if err != nil {
		log.Printf("could not get mailing list: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not get mailing list")
		return
	}

	mailingList.Recipients = recipients

	helpers.WriteJSON(w, http.StatusOK, mailingList)
}

func GetMailingList(w http.ResponseWriter, r *http.Request) {

	db := database.DbConn
	repository := Repository{Conn: db}
	muxVar := mux.Vars(r)
	strID := muxVar["id"]
	intID, err := helpers.ParseInt64(strID)
	if err != nil {
		log.Printf("could not get recipientst: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not get recipients")
		return
	}

	recipients, err := repository.GetRecipientsFromMailingList(intID)
	if err != nil {
		log.Printf("could not get recipientst: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not get recipients")
		return
	}

	mailingList, err := repository.GetMailingList(intID)
	if err != nil {
		log.Printf("could not get mailing list: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not get mailing list")
		return
	}

	mailingList.Recipients = recipients

	helpers.WriteJSON(w, http.StatusOK, mailingList)
}

func DeleteRecipientsFromMailinglist(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn
	repository := Repository{Conn: db}
	muxVar := mux.Vars(r)
	strID := muxVar["id"]
	intID, err := helpers.ParseInt64(strID)
	if err != nil {
		log.Printf("could not parse id: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not parse id")
		return
	}

	recipientIDS := make([]int64, 0)
	err = json.NewDecoder(r.Body).Decode(&recipientIDS)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not decode request body")
		return
	}

	fmt.Printf("recipients : %v", recipientIDS)

	deletedRows, err := repository.DeleteRecipientsFromMailingList(intID, recipientIDS)
	if err != nil {
		log.Printf("could not delete recipient from mailing list: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not delete recipient from mailing list:")
		return
	}
	log.Printf("deleted %d recipients", deletedRows)
	helpers.WriteJSON(w, http.StatusOK, nil)
}

type RabbitMqEnv struct {
	RabbitMqHost string `env:"RABBITMQ_HOST"`
	RabbitMqPort string `env:"RABBITMQ_PORT"`
	RabbitMqUser string `env:"RABBITMQ_DEFAULT_USER"`
	RabbitMqPass string `env:"RABBITMQ_DEFAULT_PASS"`
}

func SendCampaignMessage(w http.ResponseWriter, r *http.Request) {
	cfg := RabbitMqEnv{}
	if err := env.Parse(&cfg); err != nil {
		helpers.FailOnError(err, "Failed to parse env")
	}

	fmt.Printf("%+v", cfg)

	muxVars := mux.Vars(r)
	campaignID := muxVars["id"]

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMqPass,
		cfg.RabbitMqUser,
		cfg.RabbitMqHost,
		cfg.RabbitMqPort,
	))
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

	body := "Message to send campaign!" + campaignID
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	helpers.FailOnError(err, "Failed to publish a message")
}

//TODO: add env vars
