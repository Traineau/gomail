package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	router = gin.Default()
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//A sample user
var user = User{
	ID:       1,
	Username: "username",
	Password: "password",
}

func main() {
	router.POST("/login", Login)
	router.POST("/signin", Signin)
	router.GET("/welcome", Welcome)
	log.Fatal(router.Run(":8080"))
}

//TODO : Check the user from the database and not from the hardcoded user
func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	
	c.SetCookie("jwt", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, token)
}

//TODO : Create the new user and store it in database
func Signin(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	token, err := CreateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	
	c.SetCookie("jwt", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, token)
}


//TODO : Check the JWT stored in the user browser cookie
func Welcome(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, cookie)
}


func CreateToken(userid uint64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}