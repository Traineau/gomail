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
	Name          string `json:"name"`
	Description   string `json:"description"`
	TemplateName  string `json:"template_name"`
	TemplatePath  string `json:"template_path"`
	IdMailingList int64  `json:"id_mailing_list"`
}

// MailingList with recipients emails
type MailingList struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Recipients  []*Recipient `json:"recipients" db:"-"`
}

type CampaignID struct {
	ID          int64        `json:"id"`
}
