package main

import (
	"github.com/couchbase/gocb"
	"github.com/graphql-go/graphql"
)

/*
type Classroom struct {
	CRID              string    `json:"crid"`
	ProfessorEmail    string    `json:"professoremail"`
	ProfessorID       string    `json:"professorid"`
	DCordServerID     string    `json:"dcordserverid"`
	DCordConnected    bool      `json:"dcordconnected"`
}
*/

func CreateClassroomDiscord(params graphql.Params, collection gocb.Collection) {
	// Classroom creation from the Discord Bot
	// Information passed: Linking Token: Token contains: Classroom ID, prof email
	// Prof Discord ID passed, Discord server ID
}

func CreateClassroomFrontEnd(params graphql.Params, collection gocb.Collection) {
	//TODO: Implement
}
