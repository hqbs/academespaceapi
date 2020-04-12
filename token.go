package main

import (
	"fmt"
	"os"
	"time"

	"github.com/couchbase/gocb"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func getEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading env file")
		//log.Fatal(err)
	}

}

func GenToken(email string, id string) (UserToken, APIError) {
	newToken := UserToken{}
	newError := APIError{
		Error:   false,
		Message: "",
	}
	// getEnv() For prod
	jwtSecret := fmt.Sprintf("%s%s", os.Getenv("JWT_SECRET"), id)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		newError.Error = true
		newError.Message = "Validation Error, please try again later. Dev Code: JWTERRGEN"
	}
	newToken.Token = tokenString
	newToken.ExpireDate = expirationTime.Unix()

	return newToken, newError

}

func ValidateToken(tokenString string, userTokenInfo UserToken, id string) bool {

	// Initialize a new instance of `Claims`
	if tokenString == userTokenInfo.Token {

		// Pass 1 Tokens are equal!
		if time.Now().Unix() > userTokenInfo.ExpireDate {
			// TOKEN EXPIRED !
			return false
		} else {
			claims := &Claims{}
			// getEnv() for prod
			jwtSecret := fmt.Sprintf("%s%s", os.Getenv("JWT_SECRET"), id)
			// Parse the JWT string and store the result in `claims`.
			// Note that we are passing the key in this method as well. This method will return an error
			// if the token is invalid (if it has expired according to the expiry time we set on sign in),
			// or if the signature does not match
			tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
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

func RenewToken(email string, id string, tokenString string, tokenExpire int64) (UserToken, APIError) {
	renewToken := UserToken{}
	renewError := APIError{}
	// If expires in 30 seconds
	if tokenExpire > (time.Now().Add(time.Hour * 24)).Unix() {
		// Not expired
	} else {
		// Generate a new token, pass token back
		renewToken, renewError = GenToken(email, id)

	}

	return renewToken, renewError

}

func GetCurrentToken(email string, collection *gocb.Collection) (UserToken, APIError) {
	returnToken := UserToken{}
	returnError := APIError{}
	ops := []gocb.LookupInSpec{
		gocb.GetSpec("token.token", &gocb.GetSpecOptions{}),
		gocb.GetSpec("token.expiredate", &gocb.GetSpecOptions{}),
	}
	getResult, err := collection.LookupIn(email, ops, &gocb.LookupInOptions{})
	if err != nil {
		panic(err)
		//TODO: Create API Err
	}

	var currentToken string
	var currentExpireDate int64
	err = getResult.ContentAt(0, &currentToken)
	if err != nil {
		panic(err)
		// Create API Err
	}
	err = getResult.ContentAt(1, &currentExpireDate)
	if err != nil {
		panic(err)
		// Create API Err
	}
	returnToken.Token = currentToken
	returnToken.ExpireDate = currentExpireDate
	return returnToken, returnError
}
