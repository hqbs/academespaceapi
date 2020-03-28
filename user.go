package main

import (
	"fmt"
	"sort"

	"log"

	"github.com/couchbase/gocb"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/validator.v2"
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

func ValidateInfo(params graphql.ResolveParams) ValidatedUser {
	idHash, err := bcrypt.GenerateFromPassword([]byte(params.Args["email"].(string)), 10)
	validateUser := UserValidator{
		FName:       params.Args["fname"].(string),
		LName:       params.Args["lname"].(string),
		Email:       params.Args["email"].(string),
		PhoneNumber: params.Args["phonenumber"].(string),
		Type:        params.Args["type"].(string),
		ID:          string(idHash),
		Username:    params.Args["username"].(string),
		Password:    params.Args["password"].(string),
	}

	if err != nil {
		//TODO: Handle
	}
	validateUser.ID = string(idHash)

	err = validator.Validate(validateUser)
	var errOuts []string
	var userValid bool
	if err == nil {
		//println("Values are valid!")
		userValid = true
	} else {
		errs := err.(validator.ErrorMap)

		for f, e := range errs {
			errOuts = append(errOuts, fmt.Sprintf("\t - %s (%v)\n", f, e))
		}

		sort.Strings(errOuts)
		userValid = false
	}
	passHash, err := bcrypt.GenerateFromPassword([]byte(params.Args["password"].(string)), 12)
	if err != nil {
		//TODO: handle
	}
	returnUser := ValidatedUser{
		ValidUser: User{
			FName:       params.Args["fname"].(string),
			LName:       params.Args["lname"].(string),
			Email:       params.Args["email"].(string),
			PhoneNumber: params.Args["phonenumber"].(string),
			Type:        params.Args["type"].(string),
			ID:          string(idHash),
			Username:    params.Args["username"].(string),
			Password:    string(passHash),
		},

		UserValid: userValid,
		Errors:    errOuts,
	}

	return returnUser
}

func validatePassword(password string) bool {

	return false
}

func NewUser(userInfo ValidatedUser, collection *gocb.Collection) bool {
	if userInfo.UserValid {

		checkUser, err := collection.Get(userInfo.ValidUser.Email, nil)
		if err != nil {
			// TODO: API json panic
			log.Fatal(err)
		}
		if checkUser != nil {
			// TODO: API JSON return of why didnt work
		} else {
			_, err = collection.Upsert(userInfo.ValidUser.Email, userInfo.ValidUser, &gocb.UpsertOptions{})
			if err != nil {
				//TODO: Handle
				log.Fatal(err)
			}

			if err != nil {
				//TODO: handle
				log.Fatal(err)
			}

		}

	} else {
		return false
	}

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
