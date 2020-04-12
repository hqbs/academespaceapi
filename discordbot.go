package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
)

func DiscordTokenGen(email string) (DiscordConnectToken, APIError) {
	newToken := DiscordConnectToken{}
	newError := APIError{
		Error:   false,
		Message: "",
	}
	getEnv()
	jwtSecret := fmt.Sprintf("%s%s", os.Getenv("JWT_SECRET"), os.Getenv("DISCORD_JWT_SECRET"))
	expirationTime := time.Now().Add(30 * time.Minute)
	classroomID, err := bcrypt.GenerateFromPassword([]byte(email), 4)
	if err != nil {
		newError.Error = true
		newError.Message = "Discord Connection Setup Error, please try again later. Dev Code: ERRHASHGEN"
	}
	claims := DiscordClaims{
		Email:       email,
		ClassroomID: classroomID,
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

func DiscordValidateToken(tokenString string) (bool, string, string) {

	// Called from discord classroom creation
	// Returns bool if correct and the prof email
	email := ""
	claims := &DiscordClaims{}
	getEnv()
	jwtSecret := fmt.Sprintf("%s%s", os.Getenv("JWT_SECRET"), os.Getenv("DISCORD_JWT_SECRET"))
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {

			return false, email, ""
		}
	}
	if !tkn.Valid {

		return false, email, ""
	}
	if time.Now().Unix() > claims.StandardClaims.ExpiresAt {
		return false, email, ""
	}
	return true, claims.Email, claims.ClassroomID
}
