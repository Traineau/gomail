package email

import (
	"encoding/json"
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

// func CreateCampaign(w http.ResponseWriter, r *http.Request){
// 	var campaign Campaign
// 	err := json.NewDecoder(r.Body).Decode(&campaign)
// 	if err != nil {
// 		log.Print(err)
// 		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not decode request body")
// 		return
// 	}

// 	log.Printf("campaign : %+v", campaign)
// }

// func AddRecipientToMailinglist(w http.ResponseWriter, r *http.Request){
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
