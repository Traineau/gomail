package email

import (
	"database/sql"
)

// Repository struct for db connection
type Repository struct {
	Conn *sql.DB
}

type Campaign struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	IdMailingList int64  `json:"id_mailing_list"`
}

func (repository *Repository) GetCampaign(id int64) (*Campaign, error) {
	row := repository.Conn.QueryRow("SELECT c.id, c.name, c.description, c.id_mailing_list FROM campaign c "+
		"WHERE c.id=(?)", id)
	var idMailingList int64
	var name, description string
	switch err := row.Scan(&id, &name, &description, &idMailingList); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		campaign := Campaign{
			ID:              id,
			Name:            name,
			Description:     description,
			IdMailingList: idMailingList,
		}
		return &campaign, nil
	default:
		return nil, err
	}
}

func (repository *Repository) SaveCampaign(campaign *Campaign) error {
	stmt, err := repository.Conn.Prepare("INSERT INTO campaign(name, description, id_mailing_list) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	res, errExec := stmt.Exec(campaign.Name, campaign.Description, campaign.IdMailingList)
	if errExec != nil {
		return errExec
	}

	lastInsertedID, errInsert := res.LastInsertId()
	if errInsert != nil {
		return errInsert
	}

	campaign.ID = lastInsertedID

	return nil
}
