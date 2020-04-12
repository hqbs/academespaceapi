package main

import (
	"time"

	"github.com/couchbase/gocb"
	"github.com/graphql-go/graphql"
)

func CreateClassroomDiscord(params graphql.ResolveParams, collectionClass *gocb.Collection) APIError {
	// Classroom creation from the Discord Bot
	// Information passed: Linking Token: Token contains: Classroom ID, prof email
	// Prof Discord ID passed, Discord server ID
	returnError := APIError{}
	newClassroom := Classroom{}
	valid, email, classID := DiscordValidateToken(params.Args["token"].(string))
	if valid {
		newClassroom = Classroom{
			CRID:             classID,
			ProfessorDCordID: params.Args["professordcordid"].(string),
			DCordServerID:    params.Args["dcordserverid"].(string),
			DCordConnected:   true,
			ProfessorEmail:   email,
		}
	} else {
		//TODO: Return token error
	}
	_, err := collectionClass.Upsert(classID, newClassroom, &gocb.UpsertOptions{})
	if err != nil {
		//TODO: Return db upsert error
	}
	return returnError

}

func CreateClassroomFrontEnd(params graphql.ResolveParams, collectionClass *gocb.Collection) APIError {
	newClassroom := Classroom{}
	returnError := APIError{}
	valid, _, classID := DiscordValidateToken(params.Args["token"].(string))
	if valid {
		// Search by the classID
		// Verify email
		// Add additional details and generate join code
		newClassroom = Classroom{
			ClassName:     params.Args["classname"].(string),
			ClassNumber:   params.Args["classnumber"].(string),
			SectionNumber: params.Args["sectionnumber"].(string),
		}
		mops := []gocb.MutateInSpec{
			gocb.UpsertSpec("classname", newClassroom.ClassName, &gocb.UpsertSpecOptions{}),
			gocb.UpsertSpec("classnumber", newClassroom.ClassNumber, &gocb.UpsertSpecOptions{}),
			gocb.UpsertSpec("sectionnumber", newClassroom.SectionNumber, &gocb.UpsertSpecOptions{}),
		}
		_, err := collectionClass.MutateIn(classID, mops, &gocb.MutateInOptions{
			Timeout: 50 * time.Millisecond,
		})
		if err != nil {
			//TODO: API error mutate error
		}
	}
	return returnError
}
