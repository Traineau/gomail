package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"jwt-todo/database"
	"jwt-todo/helpers"
	"log"
	"net/http"
	"time"
)

var jwtKey = []byte("my_secret_key")

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

	log.Printf("creds : %v", creds)

	db := database.DbConn
	repository := Repository{Conn: db}

	user, err := repository.getUser(creds.Username, creds.Password)
	if err != nil {
		log.Printf("could not get user: %v", err)
	}
	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
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

	tokenString, err := token.SignedString(jwtKey)
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

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not decode request body")
		return
	}

	log.Printf("userrr : %+v", user)

	db := database.DbConn
	repository := Repository{Conn: db}

	userFromDB, err := repository.getUser(user.Username, user.Password)
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

	err = repository.saveUser(&user)
	if err != nil {
		log.Print(err)
		helpers.WriteErrorJSON(w, http.StatusInternalServerError, "could not save user in db")
		return
	}

}

func Welcome(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			helpers.WriteErrorJSON(w, http.StatusUnauthorized, "user is not logged in")
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			helpers.WriteErrorJSON(w, http.StatusUnauthorized, "invalid signature")
			return
		}
		helpers.WriteErrorJSON(w, http.StatusBadRequest, "Bad request")
		return
	}
	if !tkn.Valid {
		helpers.WriteErrorJSON(w, http.StatusUnauthorized, "invalid token")
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
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
		return jwtKey, nil
	})
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
