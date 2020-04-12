package main

import (
	"fmt"
	"sort"

	"github.com/couchbase/gocb"
	"github.com/graphql-go/graphql"
	"golang.org/x/crypto/bcrypt"
	validator "gopkg.in/validator.v2"
)

/*
	TODO: Update user validation to deal with new fields in the User Struct in main
	TODO: Determine what is needed for new user creation with the new fields
	TODO: Hold it all together and keep moving forward!!
*/
type UserValidator struct {
	FName       string `validate:"nonzero,min=2,max=100"`
	LName       string `validate:"nonzero,min=2,max=100"`
	Email       string `validate:"nonzero"`
	PhoneNumber string `validate:"min=4,max=40,regexp=^(\(?\+?[0-9]*\)?)?[0-9_\- \(\)]*$"` // http://regexlib.com/REDetails.aspx?regexp_id=73
	Type        string `validate:"nonzero"`
	ID          string `validate:"nonzero"`
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
			Password:    string(passHash),
		},

		UserValid: userValid,
		Errors:    errOuts,
	}
	newToken, tokenErr := GenToken(params.Args["email"].(string), returnUser.ValidUser.ID)
	returnUser.ValidUser.Token = newToken
	if tokenErr.Error {
		returnUser.Errors = append(returnUser.Errors, tokenErr.Message)
	}
	return returnUser
}

func NewUser(userInfo ValidatedUser, collection *gocb.Collection) (UserToken, MutationPayload) {
	returnErr := MutationPayload{}
	returnErr.Success = true
	returnToken := userInfo.ValidUser.Token
	if userInfo.UserValid {
		exists, _ := UserExist(userInfo.ValidUser.Email, collection)
		if exists {
			returnErr.Success = false
			returnErr.Errors = append(returnErr.Errors, "Email already in use!")
			returnErr.Token = ""
			returnToken.Token = ""
			returnToken.ExpireDate = 0000

		} else {
			_, err := collection.Upsert(userInfo.ValidUser.Email, userInfo.ValidUser, &gocb.UpsertOptions{})
			if err != nil {
				returnErr.Success = false
				returnErr.Errors = append(returnErr.Errors, "Account Creation Error, please try again later. Dev Code: ERRNEWUSRDBUP")
				returnErr.Token = ""
				returnToken.Token = ""
				returnToken.ExpireDate = 0000
			}

		}

	} else {
		returnErr.Errors = append(returnErr.Errors, userInfo.Errors...)
		returnErr.Success = false
		returnErr.Errors = append(returnErr.Errors, "Account Creation Error, please try again later. Dev Code: ERRNEWUSRNVU")
		returnErr.Token = ""
		returnToken.Token = ""
		returnToken.ExpireDate = 0000
		return returnToken, returnErr
	}

	return returnToken, returnErr
}

func UpdateUser(modifyDetails ModifyUser, collection *gocb.Collection) bool {
	//TODO: Implement
	return false
}

func RemoveUser(userInfo User) bool {
	//TODO: Implement
	return false
}

func UserExist(email string, collection *gocb.Collection) (bool, APIError) {
	apiErr := APIError{}
	checkUser, _ := collection.Get(email, nil)
	//TODO: Swap to check exist function couchbase
	// if err != nil {
	// 	if err.error_name == "KEY_ENOENT" {
	// 		return false, apiErr
	// 	}
	// 	apiErr.Error = true
	// 	apiErr.Message = "Account Validation Error, please try again later."
	// 	panic(err)
	// 	return false, apiErr
	// }
	if checkUser != nil {
		return true, apiErr
	}
	return false, apiErr

}
