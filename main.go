package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fabwi987/irecommend/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		decoded := []byte(os.Getenv("AUTH0_CLIENT_SECRET_API"))
		if len(decoded) == 0 {
			return nil, errors.New("Missing Client Secret")
		}

		//Extracting the userID from the jwt token
		Claim := token.Claims.(jwt.MapClaims)
		//models.UseridClaim = "admin"
		models.UseridClaim = Claim["sub"].(string)
		return decoded, nil
	},
})

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Handle("/confirm/{id}", http.FileServer(http.Dir("./views/start")))

	http.ListenAndServe(":7070", handlers.LoggingHandler(os.Stdout, r))
}
