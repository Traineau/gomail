package users

import (
	"database/sql"
	"log"
)

// Repository struct for db connection
type Repository struct {
	Conn *sql.DB
}

//User is a user model
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//GetUser is for getting a user by username
func (repository *Repository) GetUser(username string) (*User, error) {
	row := repository.Conn.QueryRow("SELECT u.id, u.username, u.email, u.password FROM api_user u "+
		"WHERE u.username=(?)", username)
	var id int64
	var email, password string
	switch err := row.Scan(&id, &username, &email, &password); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		user := User{
			ID:       id,
			Username: username,
			Email:    email,
			Password: password,
		}
		return &user, nil
	default:
		return nil, err
	}
}

//SaveUser is for saving a new user
func (repository *Repository) SaveUser(user *User) error {
	stmt, err := repository.Conn.Prepare("INSERT INTO api_user(username, email, password) VALUES(?,?,?)")
	if err != nil {
		return err
	}

	log.Printf("\nuser : %+v", user)

	res, errExec := stmt.Exec(user.Username, user.Email, user.Password)
	if errExec != nil {
		return errExec
	}

	lastInsertedID, errInsert := res.LastInsertId()
	if errInsert != nil {
		return errInsert
	}

	user.ID = lastInsertedID

	return nil
}
