package main

import "github.com/graphql-go/graphql"

type UserValidator struct {
	FName     string        `validate:"nonzero"`
	LName     string        `validate:"nonzero"`
	Email     string        `validate:"regexp=^[a-zA-Z0-9.!#$%&'*+/=?^_{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"`
	Type      string        `json:"type"`
	ID        int64         `json:"id"`
	Username  string        `json:"username"`
	Password  string        `json:"password"`
	DiscordID string        `json:"discordid,omitempty"`
	Servers   []DCordServer `json:"server,omitempty"`
}

func ValidateInfo(params graphql.ResolveParams) ValidatedUser {

	returnUser := ValidatedUser{
		// User {
		// FName: params.Args["fname"].(string),
		// LName: params.Args["lname"].(string),
		// Email: params.Args["email"].(string),
		// Type:  params.Args["type"].(string),
		// },
	}
	if len(params.Args["username"].(string)) <= 100 {
		// do something
	}
	return returnUser
}

func validatePassword(password string) bool {

	return false
}

func NewUser(userInfo User) bool {
	//TODO: Implement

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
