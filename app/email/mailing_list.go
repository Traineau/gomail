package email

import (
	"database/sql"
)

func (repository *Repository) GetMailingList(id int64) (*MailingList, error) {
	row := repository.Conn.QueryRow("SELECT m.id, m.name, m.description FROM mailing_list m "+
		"WHERE m.id=(?)", id)
	var name, description string
	switch err := row.Scan(&id, &name, &description); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		mailingList := MailingList{
			ID:          id,
			Name:        name,
			Description: description,
		}
		return &mailingList, nil
	default:
		return nil, err
	}
}

func (repository *Repository) SaveMailingList(mailingList *MailingList) error {
	stmt, err := repository.Conn.Prepare("INSERT INTO mailing_list(name, description) VALUES(?,?)")
	if err != nil {
		return err
	}

	res, errExec := stmt.Exec(mailingList.Name, mailingList.Description)
	if errExec != nil {
		return errExec
	}

	lastInsertedID, errInsert := res.LastInsertId()
	if errInsert != nil {
		return errInsert
	}

	mailingList.ID = lastInsertedID

	return nil
}
