package user

import (
	"encoding/json"
	"github.com/Traineau/gomail/database"
	"github.com/Traineau/gomail/helpers"
	"github.com/Traineau/gomail/helpers/password"
	"github.com/Traineau/gomail/users"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)

var JwtKey = []byte("my_secret_key")

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

//
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Signin(w http.ResponseWriter, r *http.Request) {

	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := database.DbConn
	repository := users.Repository{Conn: db}

	user, err := repository.GetUser(creds.Username)
	if err != nil {
		log.Printf("could not get user: %v", err)
		return
	}
	if user == nil {
		log.Print("no user found")
		helpers.WriteErrorJSON(w, http.StatusBadRequest, "no user to connect")
		return
	}
	isMatching, err := password.ComparePasswordAndHash(creds.Password, user.Password)
	if err != nil {
		log.Printf("could not compare password: %v", err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not compare password")
		return
	}
	if !isMatching {
		log.Print("password not matching")
		helpers.WriteErrorJSON(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	helpers.WriteJSON(w, http.StatusOK, "user logged in")
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not decode request body")
		return
	}

	db := database.DbConn
	repository := users.Repository{Conn: db}

	userFromDB, err := repository.GetUser(user.Username)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not get user from db")
		return
	}

	if userFromDB != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusBadRequest, "user already exists")
		return
	}

	hash, err := password.GenerateFromPassword(user.Password)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not safely save user")
		return
	}
	user.Password = hash

	err = repository.SaveUser(&user)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not save user in db")
		return
	}

	helpers.WriteJSON(w, http.StatusOK, "user registered")

}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(2 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
