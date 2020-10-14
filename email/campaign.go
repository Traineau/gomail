package email

import (
	"database/sql"
)

//GetCampaign get a campaign by id in the database
func (repository *Repository) GetCampaign(id int64) (*Campaign, error) {
	row := repository.Conn.QueryRow("SELECT c.id, c.name, c.description, c.template_name, c.template_path, c.id_mailing_list FROM campaign c "+
		"WHERE c.id=(?)", id)
	var idMailingList int64
	var name, description, templateName, templatePath string
	switch err := row.Scan(&id, &name, &description, &templateName, &templatePath, &idMailingList); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		campaign := Campaign{
			ID:            id,
			Name:          name,
			Description:   description,
			TemplateName:  templateName,
			TemplatePath:  templatePath,
			IDMailingList: idMailingList,
		}
		return &campaign, nil
	default:
		return nil, err
	}
}

//SaveCampaign save a new campaign in database
func (repository *Repository) SaveCampaign(campaign *Campaign) error {
	stmt, err := repository.Conn.Prepare(`INSERT INTO campaign(name, description, 
id_mailing_list, template_name, template_path) VALUES(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	res, errExec := stmt.Exec(campaign.Name, campaign.Description, campaign.IDMailingList, 
		campaign.TemplateName, campaign.TemplatePath)

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
