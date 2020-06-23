package email

import (
	"database/sql"
	"log"
)

// Repository struct for db connection
type Repository struct {
	Conn *sql.DB
}

type Campaign struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	IdDiffusionList int64  `json:"id_diffusion_list"`
}

func (repository *Repository) GetCampaign(id int64) (*Campaign, error) {
	row := repository.Conn.QueryRow("SELECT c.id, c.name, c.description, c.id_diffusion_list FROM campaign c "+
		"WHERE c.id=(?)", id)
	var idDiffusionList int64
	var name, description string
	switch err := row.Scan(&id, &name, &description, &idDiffusionList); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		campaign := Campaign{
			ID:              id,
			Name:            name,
			Description:     description,
			IdDiffusionList: idDiffusionList,
		}
		return &campaign, nil
	default:
		return nil, err
	}
}

func (repository *Repository) SaveCampaign(campaign *Campaign) error {
	stmt, err := repository.Conn.Prepare("INSERT INTO campaign(name, description, id_diffusion_list) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	log.Printf("\nuser : %+v", campaign)

	res, errExec := stmt.Exec(campaign.Name, campaign.Description, campaign.IdDiffusionList)
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
