package router

import (
	"github.com/Traineau/gomail/email"
	"github.com/Traineau/gomail/helpers"
	"github.com/Traineau/gomail/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Route struct defining all of this project routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Public      bool
}

// Routes slice of Route
type Routes []Route

// newRouter registers public routes
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	authenticatedRouter := router.PathPrefix("/").Subrouter()

	for _, route := range routes {
		appRouter := authenticatedRouter
		if route.Public {
			appRouter = router
		}
		appRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	authenticatedRouter.Use(loggingMiddleware)

	return router
}

var routes = Routes{
	Route{
		Name:        "Sign In",
		Method:      "POST",
		Pattern:     "/signin",
		HandlerFunc: user.Signin,
		Public:      true,
	},
	Route{
		Name:        "Sign Up",
		Method:      "POST",
		Pattern:     "/signup",
		HandlerFunc: user.SignUp,
		Public:      true,
	},
	Route{
		Name:        "Refresh",
		Method:      "GET",
		HandlerFunc: user.Refresh,
		Public:      false,
	},
	Route{
		Name:        "Create mailing list",
		Method:      "POST",
		Pattern:     "/mailinglist",
		HandlerFunc: email.CreateMailingList,
		Public:      false,
	},
	Route{
		Name:        "Add recipient to mailing list",
		Method:      "POST",
		Pattern:     "/mailinglist/{id}",
		HandlerFunc: email.AddRecipientToMailinglist,
		Public:      false,
	},
	Route{
		Name:        "Create campaign",
		Method:      "POST",
		Pattern:     "/campaign",
		HandlerFunc: email.CreateCampaign,
		Public:      false,
	},
	Route{
		Name:        "Send message",
		Method:      "POST",
		Pattern:     "/campaign/{id}/send",
		HandlerFunc: email.SendCampaignMessage,
		Public:      false,
	},
	Route{
		Name:        "Get mailing list",
		Method:      "GET",
		Pattern:     "/mailinglist/{id}",
		HandlerFunc: email.GetMailingList,
		Public:      false,
	},
	Route{
		Name:        "Delete recipients from mailing list",
		Method:      "DELETE",
		Pattern:     "/mailinglist/{id}/recipients",
		HandlerFunc: email.DeleteRecipientsFromMailinglist,
		Public:      false,
	},
	// TODO : Get campaign template, path and recipients
	// Route{
	// 	Name:        "Get campaign template, path and recipients",
	// 	Method:      "GET",
	// 	Pattern:     "",
	// 	HandlerFunc: ,
	// 	Public:      false,
	// },
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
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
		claims := &user.Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return user.JwtKey, nil
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

		user.Refresh(w, r)

		next.ServeHTTP(w, r)
	})
}
