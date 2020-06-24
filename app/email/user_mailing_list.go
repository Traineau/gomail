package email

import (
	"database/sql"
	"fmt"
	"gomail/users"
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
		if err != nil {
			return nil, fmt.Errorf("could not get emails : %v", err)
		}

		usersEmails = append(usersEmails, email)

	}

	return usersEmails, nil
}

func (repository *Repository) AddUserToMailingList(user *users.User, mailingListID int64) error {
	stmt, err := repository.Conn.Prepare("INSERT INTO user_mailing_list(id_user, id_mailing_list) VALUES(?,?)")
	if err != nil {
		return err
	}

	res, errExec := stmt.Exec(user.ID, mailingListID)
	if errExec != nil {
		return fmt.Errorf("could not exec stmt: %v", errExec)
	}

	_, errInsert := res.LastInsertId()
	if errInsert != nil {
		return fmt.Errorf("could not retrieve last inserted ID: %v", errInsert)
	}

	return nil
}

func (repository *Repository) DeleteUserFromMailingList(user *users.User, mailingListID int64) (int64, error) {
	res, err := repository.Conn.Exec("DELETE FROM user_mailing_list WHERE id_user=(?) "+
		"AND id_mailing_list=(?)", user.ID, mailingListID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
