package email

import (
	"encoding/json"
	"github.com/gorilla/mux"
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
		return
	}

	muxVars := mux.Vars(r)
	mailingListID, err := helpers.ParseInt64(muxVars["id"])
	if err != nil {
		log.Printf("could not save parse id into int: %v", err)
		return
	}

	err = repository.AddRecipientToMailingList(ids, mailingListID)
	if err != nil {
		log.Printf("could not save recipients list: %v", err)
		return
	}

	log.Printf("ids des recipients : %+v", ids)

	helpers.WriteJSON(w, http.StatusOK, "c bon")
}

// func DeleteRecipientFromMailinglist(w http.ResponseWriter, r *http.Request){
// 	muxVars := mux.Vars(r)
// 	id := muxVars["id"]

// 	var recipient Recipient
// 	err := json.NewDecoder(r.Body).Decode(&recipient)
// 	if err != nil {
// 		log.Print(err)
// 		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not decode request body")
// 		return
// 	}

// 	log.Printf("mailing list : %+v", recipient)

// }
