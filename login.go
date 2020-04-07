package main

import (
	"time"

	"github.com/couchbase/gocb"
	"golang.org/x/crypto/bcrypt"
)

func Login(email string, password string, jwt string, collection *gocb.Collection) (UserToken, APIError) {
	returnError := APIError{}
	emptyUserReturn := UserToken{}
	var dbPassword string
	var dbID string
	var dbToken string
	var dbExpireDate int64
	var dbUserToken UserToken
	userExists, errors := UserExist(email, collection)
	if errors.Error {
		return emptyUserReturn, errors
	} else if !userExists {
		returnError.Error = true
		returnError.Message = "User does not exist!"
		return emptyUserReturn, returnError
	} else {
		// User Exists and no errors!
		ops := []gocb.LookupInSpec{
			gocb.GetSpec("password", &gocb.GetSpecOptions{}),
			gocb.GetSpec("id", &gocb.GetSpecOptions{}),
			gocb.GetSpec("token.token", &gocb.GetSpecOptions{}),
			gocb.GetSpec("token.expiredate", &gocb.GetSpecOptions{}),
		}
		getResult, err := collection.LookupIn(email, ops, &gocb.LookupInOptions{})
		if err != nil {
			returnError.Error = true
			returnError.Message = "Account Retrieval Error, please try again later. Dev Code: LOGGETERR"
			return emptyUserReturn, returnError
		}

		err = getResult.ContentAt(0, &dbPassword)
		if err != nil {
			returnError.Error = true
			returnError.Message = "Account Retrieval Error, please try again later. Dev Code: LOGGETERR"
			return emptyUserReturn, returnError
		}
		err = getResult.ContentAt(1, &dbID)
		if err != nil {
			returnError.Error = true
			returnError.Message = "Account Retrieval Error, please try again later. Dev Code: LOGGETERR"
			return emptyUserReturn, returnError
		}
		err = getResult.ContentAt(2, &dbToken)
		if err != nil {
			returnError.Error = true
			returnError.Message = "Account Retrieval Error, please try again later. Dev Code: LOGGETERR"
			return emptyUserReturn, returnError
		}
		err = getResult.ContentAt(3, &dbExpireDate)
		if err != nil {
			returnError.Error = true
			returnError.Message = "Account Retrieval Error, please try again later. Dev Code: LOGGETERR"
			return emptyUserReturn, returnError
		}
		dbUserToken = UserToken{
			Token:      dbToken,
			ExpireDate: dbExpireDate,
		}
		if !ValidateToken(jwt, dbUserToken, dbID) {

			err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
			if err != nil {
				returnError.Error = true
				returnError.Message = "Incorrect password!"
				return emptyUserReturn, returnError
			} else {
				//TODO: Generate and return new token
				newToken, tokenGenErrors := GenToken(email, dbID)
				if tokenGenErrors.Error {
					returnError.Error = true
					returnError.Message = "Token Generation Error, please try again later. Dev Code: LOGTOKGENERR"
					return emptyUserReturn, returnError
				} else {
					// No error return token and add to DB
					mops := []gocb.MutateInSpec{
						gocb.UpsertSpec("token.token", newToken.Token, &gocb.UpsertSpecOptions{}),
						gocb.UpsertSpec("token.expiredate", newToken.ExpireDate, &gocb.UpsertSpecOptions{}),
					}
					_, err := collection.MutateIn(email, mops, &gocb.MutateInOptions{
						Timeout: 50 * time.Millisecond,
					})
					if err != nil {
						returnError.Error = true
						returnError.Message = "Account Update Error, please try again later. Dev Code: LOGTOKENUPERR"
						return emptyUserReturn, returnError
					}
					return newToken, returnError
				}
			}
		}

	}
	return dbUserToken, returnError
}
