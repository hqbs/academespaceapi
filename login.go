package main

import (
	"time"

	"github.com/couchbase/gocb"
	"golang.org/x/crypto/bcrypt"
)

func LogIn(email string, password string, jwt string, collection *gocb.Collection) (UserToken, APIError) {
	returnError := APIError{}
	ops := []gocb.LookupInSpec{
		gocb.GetSpec("passoword", &gocb.GetSpecOptions{}),
		gocb.GetSpec("id", &gocb.GetSpecOptions{}),
		gocb.GetSpec("token.token", &gocb.GetSpecOptions{}),
		gocb.GetSpec("token.expiredate", &gocb.GetSpecOptions{}),
	}
	getResult, err := collection.LookupIn(email, ops, &gocb.LookupInOptions{})
	if err != nil {
		panic(err)
		//TODO: Create API Err
	}

	var dbPassword string
	var dbID string
	var dbToken string
	var dbExpireDate int64
	err = getResult.ContentAt(0, &dbPassword)
	if err != nil {
		panic(err)
		//TODO: Create API Err
	}
	err = getResult.ContentAt(1, &dbID)
	if err != nil {
		panic(err)
		//TODO: Create API Err
	}
	err = getResult.ContentAt(2, &dbToken)
	if err != nil {
		panic(err)
		//TODO: Create API Err
	}
	err = getResult.ContentAt(3, &dbExpireDate)
	if err != nil {
		panic(err)
		//TODO: Create API Err
	}
	dbUserToken := UserToken{
		Token:      dbToken,
		ExpireDate: dbExpireDate,
	}
	if !ValidateToken(jwt, dbUserToken) {
		// Need to login
		userExists, errors := UserExist(email, collection)
		if errors.Error {
			//TODO: Handle error
		} else if !userExists {
			//TODO: Handle User doesn't exist
		} else {
			// User Exists and no errors!

			err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
			if err != nil {
				//TODO: Pass is incorrect
			} else {
				//TODO: Generate and return new token
				newToken, tokenGenErrors := GenToken(email, dbID)
				if tokenGenErrors.Error {
					//TODO: Handle error
				} else {
					// No error return token and add to DB
					mops := []gocb.MutateInSpec{
						gocb.UpsertSpec("token", newToken, &gocb.UpsertSpecOptions{}),
					}
					_, err := collection.MutateIn("customer123", mops, &gocb.MutateInOptions{
						Timeout: 50 * time.Millisecond,
					})
					if err != nil {
						// TODO: token update failed
						panic(err)
					}
					return newToken, returnError
				}
			}
		}
	} else {
		return dbUserToken, returnError
	}
	return dbUserToken, returnError
}
