package email

import (
	"database/sql"
	"fmt"
)

// Repository struct for db connection
func (repository *Repository) GetRecipientsFromMailingList(id int64) ([]*Recipient, error) {
	rows, err := repository.Conn.Query("SELECT r.id, r.email, r.first_name, r.last_name, r.username"+
		"FROM recipient r, mailing_list ml, recipient_mailing_list rml "+
		"WHERE r.id = rml.id_recipient AND ml.id = (?);", id)
	if err != nil {
		return nil, fmt.Errorf("could not prepare query: %v", err)
	}
	var recipients []*Recipient
	var email, firstName, lastName, username string
	for rows.Next() {
		err := rows.Scan(&id, &email, &firstName, &lastName, &username)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("could not get emails : %v", err)
		}
		recipient := &Recipient{
			ID:        id,
			Email:     email,
			FirstName: firstName,
			LastName:  lastName,
			UserName:  username,
		}
		recipients = append(recipients, recipient)
	}

	return recipients, nil
}

func (repository *Repository) AddRecipientToMailingList(recipientID int64, mailingListID int64) error {
	stmt, err := repository.Conn.Prepare("INSERT INTO recipient_mailing_list(id_recipient, id_mailing_list) VALUES(?,?)")
	if err != nil {
		return err
	}

	res, errExec := stmt.Exec(recipientID, mailingListID)
	if errExec != nil {
		return fmt.Errorf("could not exec stmt: %v", errExec)
	}

	_, errInsert := res.LastInsertId()
	if errInsert != nil {
		return fmt.Errorf("could not retrieve last inserted ID: %v", errInsert)
	}

	return nil
}

func (repository *Repository) DeleteRecipientFromMailingList(recipientID, mailingListID int64) (int64, error) {
	res, err := repository.Conn.Exec("DELETE FROM recipient_mailing_list WHERE id_recipient=(?) "+
		"AND id_mailing_list=(?)", recipientID, mailingListID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
