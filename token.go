package main

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

/*
	TODO: huge overarching todo:
	- Refactor to work with the DB
	- Make sure returning exists so the API can return the correct info
	- Build these into the API and make the calls during user creation
	- Eventually implement variable security
	- Remove collection comments after testing is complete
*/

func getEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading env file")
		//log.Fatal(err)
	}

}

func genToken(email string /*collection *gocb.Collection*/) {
	getEnv()
	jwtSecret := os.Getenv("JWT_SECRET")
	expirationTime := time.Now().Add(10 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		//TODO: handle
	}
	//TODO: put token and expire date into database
	fmt.Println(tokenString)

}

func validateToken(tokenString string, userTokenInfo UserToken /*collection *gocb.Collection*/) bool {
	//TODO: start with database query to check to see if the token is the same if not invalid
	//TODO: Remove variable userTokenInfo and implement DB search
	// Initialize a new instance of `Claims`
	if tokenString == userTokenInfo.Token {
		// Pass 1 Tokens are equal!
		if time.Now().Unix() > userTokenInfo.ExpireDate {
			// TOKEN EXPIRED !
			return false
		} else {
			claims := &Claims{}
			getEnv()
			jwtSecret := os.Getenv("JWT_SECRET")
			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					return false
				}
			}
			if !tkn.Valid {
				return false
			}
			return true
		}

	} else {
		return false
	}

}

func renewToken(email string, tokenString string, tokenExpire int64 /*collection *gocb.Collection*/) {
	// If expires in 30 seconds
	if tokenExpire > (30 * time.Second).Unix() {
		// Not expired
	} else {
		// Generate a new token, pass token back
		genToken(email, collection)

	}

}
