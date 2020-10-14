package email

import (
	"database/sql"
	"fmt"
)

//GetRecipientsFromMailingList is for getting all recipients from a mailing list
func (repository *Repository) GetRecipientsFromMailingList(id int64) ([]*Recipient, error) {
	rows, err := repository.Conn.Query("SELECT r.id, r.email, r.first_name, r.last_name, r.username"+
		"\nFROM recipient_mailing_list rml "+
		"\nINNER JOIN recipient r "+
		"\nON rml.id_recipient = r.id "+
		"\nINNER JOIN mailing_list ml"+
		"\nON rml.id_mailing_list = ml.id"+
		"\n WHERE ml.id = (?);", id)
	if err != nil {
		return nil, fmt.Errorf("could not prepare query: %v", err)
	}
	var recipients []*Recipient
	var email, firstName, lastName, username string
	defer rows.Close()
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

	err = rows.Err()
	if err != nil {
		fmt.Print(err)
	}

	return recipients, nil
}

//SaveRecipients is for saving new recipients
func (repository *Repository) SaveRecipients(recipients []*Recipient) ([]int64, error) {
	sqlStr := "INSERT INTO recipient(email, first_name, last_name, username) VALUES "
	var values []interface{}

	for _, row := range recipients {
		sqlStr += "(?, ?, ?, ?),"
		values = append(values, row.Email, row.FirstName, row.LastName, row.UserName)
	}
	//trim the last
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := repository.Conn.Prepare(sqlStr)
	if err != nil {
		return nil, fmt.Errorf("could not prepare stmt: %v", err)
	}

	//format all vals at once
	res, errExec := stmt.Exec(values...)
	if errExec != nil {
		return nil, fmt.Errorf("could not exec stmt: %v", err)
	}

	var insertedIDs []int64
	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("could not exec stmt: %v", err)
	}

	for i := range recipients {
		insertedIDs = append(insertedIDs, id+int64(i))
	}

	return insertedIDs, nil
}

//AddRecipientToMailingList is for adding recipients to a mailing list
func (repository *Repository) AddRecipientToMailingList(recipientIDs []int64, mailingListID int64) error {
	sqlStr := "INSERT INTO recipient_mailing_list(id_recipient, id_mailing_list) VALUES "
	var values []interface{}

	for _, recipientID := range recipientIDs {
		sqlStr += "(?, ?),"
		values = append(values, recipientID, mailingListID)
	}
	//trim the last
	sqlStr = sqlStr[0 : len(sqlStr)-1]
	//prepare the statement
	stmt, err := repository.Conn.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("could not prepare stmt: %v", err)
	}

	//format all vals at once
	_, errExec := stmt.Exec(values...)
	if errExec != nil {
		return fmt.Errorf("could not exec stmt: %v", err)
	}

	return nil
}

//DeleteRecipientsFromMailingList is for deleting recipients from a mailing list
func (repository *Repository) DeleteRecipientsFromMailingList(mailingListID int64, recipientIDs []int64) (int64, error) {
	if len(recipientIDs) == 0 {
		return 0, nil
	}
	queryString := fmt.Sprintf("DELETE FROM recipient_mailing_list WHERE id_mailing_list=%d", mailingListID)
	queryString += fmt.Sprintf(" AND id_recipient=%d", recipientIDs[0])

	for i, id := range recipientIDs {
		if i == 0 {
			continue
		}
		queryString += fmt.Sprintf("\nOR id_recipient=%d", id)
	}

	res, err := repository.Conn.Exec(queryString)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
