package email

import (
	"encoding/json"
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

	ids, err := repository.SaveRecipients(recipients)
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

	if mailingList == nil {
		log.Printf("no mailing list by this id")
		helpers.WriteErrorJSON(w, http.StatusNotFound, "could not find mailing list")
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

	if mailingList == nil {
		log.Printf("no mailing list by this id")
		helpers.WriteErrorJSON(w, http.StatusNotFound, "could not find mailing list")
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

	deletedRows, err := repository.DeleteRecipientsFromMailingList(intID, recipientIDS)
	if err != nil {
		log.Printf("could not delete recipient from mailing list: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not delete recipient from mailing list:")
		return
	}

	recipientStr := "recipients"

	if deletedRows <= 1 {
		recipientStr = "recipient"
	}

	log.Printf("deleted %d %s", deletedRows, recipientStr)
	helpers.WriteJSON(w, http.StatusOK, nil)
}

func SendCampaignMessage(w http.ResponseWriter, r *http.Request) {
	rbmqChanCreation := RBMQQueuecreation{
		RabbitMQChan:  RabbitMQChan,
		RabbitMQQueue: RabbitMQQueue,
	}
	muxVars := mux.Vars(r)
	campaignID := muxVars["id"]

	body := "Message to send campaign!" + campaignID
	err := rbmqChanCreation.RabbitMQChan.Publish(
		"",                                  // exchange
		rbmqChanCreation.RabbitMQQueue.Name, // routing key
		false,                               // mandatory
		false,                               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	log.Printf(" [x] Sent %s", body)
	helpers.FailOnError(err, "Failed to publish a message")
}
