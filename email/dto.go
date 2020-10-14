package email

import (
	"database/sql"
)

// Repository struct for db connection
type Repository struct {
	Conn *sql.DB
}

// Recipient of an email
type Recipient struct {
	ID        int64  `json:"id,omitempty" db:"id"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	UserName  string `json:"username" db:"username"`
}

// Campaign is a marketing campaign
type Campaign struct {
	ID            int64  `json:"id"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	TemplateName  string `json:"template_name,omitempty"`
	TemplatePath  string `json:"template_path,omitempty"`
	IDMailingList int64  `json:"id_mailing_list,omitempty"`
}

// MailingList with recipients emails
type MailingList struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Recipients  []*Recipient `json:"recipients" db:"-"`
}

//CampaignID represent a campaign ID
type CampaignID struct {
	ID int64 `json:"id"`
}
