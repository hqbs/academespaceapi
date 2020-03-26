package main

import (
	"fmt"
	"sort"

	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/validator.v2"
	"github.com/joho/godotenv"
	"github.com/couchbase/gocb "
)
var (
	dbUser := os.Getenv("COUCH_USER")
	dbPass := os.Getenv("COUCH_PASS")
	dbAddr := os.Getenv("COUCH_ADDR")
	ddBucket := os.Getenv("COUCH_U_BUCKET")
)
type UserValidator struct {
	FName       string `validate:"nonzero,min=2,max=100"`
	LName       string `validate:"nonzero,min=2,max=100"`
	Email       string `validate:"nonzero"`                                                // https://www.golangprograms.com/regular-expression-to-validate-email-address.html
	PhoneNumber string `validate:"min=4,max=40,regexp=^(\(?\+?[0-9]*\)?)?[0-9_\- \(\)]*$"` // http://regexlib.com/REDetails.aspx?regexp_id=73
	Type        string `validate:"nonzero"`
	ID          string `validate:"nonzero"`
	Username    string `validate:"min=4,max=40"`
	Password    string `validate:"min=14,max=350"`
}

func ValidateInfo() ValidatedUser {

	validateUser := UserValidator{
		FName:       params.Args["fname"].(string),
		LName:       params.Args["lname"].(string),
		Email:       params.Args["email"].(string),
		PhoneNumber: params.Args["phonenumber"].(string),
		Type:        params.Args["type"].(string),
		ID:          hashID,
		Username:    params.Args["username"].(string),
		Password:    params.Args["password"].(string),
	}

	idHash, err := bcrypt.GenerateFromPassword([]byte(params.Args["email"].(string)), 10)

	if err != nil {
		//TODO: Handle 
	}
	validateUser.ID = string(idHash)
	println(string(idHash))
	err = validator.Validate(validateUser)
	var errOuts []string
	var userValid bool
	if err == nil {
		println("Values are valid!")
		userValid = true
	} else {
		errs := err.(validator.ErrorMap)

		for f, e := range errs {
			errOuts = append(errOuts, fmt.Sprintf("\t - %s (%v)\n", f, e))
		}

		sort.Strings(errOuts)
		userValid = false
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(params.Args["password"].(string),12))
	if err != nil {
		//TODO: handle
	}
	returnUser := ValidatedUser{
		FName:       params.Args["fname"].(string),
		LName:       params.Args["lname"].(string),
		Email:       params.Args["email"].(string),
		PhoneNumber: params.Args["phonenumber"].(string),
		Type:        params.Args["type"].(string),
		ID:          hashID,
		Username:    params.Args["username"].(string),
		Password:    passHash,
		UserValid:   userValid,
		Errors:      errOuts,
	}

	return returnUser
}

func validatePassword(password string) bool {

	return false
}

func NewUser(userInfo User) bool {
	// Connect to DB
	connOpts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			dbUser,
			dbPass,
		},
	}
	cluster, err := gocb.Connect(dbAddr, connOpts)
	if err != nil {
		//TODO: Handle error - will need to be sent through API
		//TODO: Most likely a struct with a payload setup of some sort
		log.Fatal(err)
	}
	bucket := cluster.Bucket(dbBucket)
	collection := bucket.DefaultCollection()
	return false
}

func UpdateUser(userInfo User, userInfoUpdated User) bool {
	//TODO: Implement
	return false
}

func RemoveUser(userInfo User) bool {
	//TODO: Implement
	return false
}
