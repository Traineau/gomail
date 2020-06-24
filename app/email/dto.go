package email

import "database/sql"

// Repository struct for db connection
type Repository struct {
	Conn *sql.DB
}

// Recipient of an email
type Recipient struct {
	ID        int64
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	UserName  string `json:"user_name,omitempty"`
}

// Campaign is a marketing campaign
type Campaign struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	IdMailingList int64  `json:"id_mailing_list"`
}

// MailingList with recipients emails
type MailingList struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
