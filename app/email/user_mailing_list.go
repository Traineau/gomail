package email

import (
	"database/sql"
	"fmt"
)

// Repository struct for db connection
func (repository *Repository) GetUserEmailsFromMailingList(id int64) ([]string, error) {
	var usersEmails []string
	var email string
	rows, err := repository.Conn.Query("SELECT user.email AS 'email' FROM user, "+
		"mailing_list, user_mailing_list WHERE user.id = user_mailing_list.id_user "+
		"AND mailing_list.id = (?);", id)
	if err != nil {
		return nil, fmt.Errorf("could not prepare query: %v", err)
	}

	for rows.Next() {
		err := rows.Scan(&email)
		if err == sql.ErrNoRows {
			return nil, nil
		}

		fmt.Printf("\n email %s", email)
		if err != nil {
			return nil, fmt.Errorf("could not get emails : %v", err)
		}

		usersEmails = append(usersEmails, email)

	}

	return usersEmails, nil
}
