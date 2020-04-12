package main

import (
	"log"

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
	result, err := collectionClass.Upsert(classID, newClassroom, &gocb.UpsertOptions{})
	log.Println(result)
	if err != nil {

		//TODO: Handle (err)
	}
	return returnError

}

func CreateClassroomFrontEnd(params graphql.ResolveParams, collectionClass *gocb.Collection, collectionUser *gocb.Collection) APIError {
	newClassroom := Classroom{}
	returnError := APIError{}

	valid, email, classID := DiscordValidateToken(params.Args["token"].(string))
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
			gocb.UpsertSpec("classname", newClassroom.ClassName, &gocb.UpsertSpecOptions{
				CreatePath: true,
			}),
			gocb.UpsertSpec("classnumber", newClassroom.ClassNumber, &gocb.UpsertSpecOptions{
				CreatePath: true,
			}),
			gocb.UpsertSpec("sectionnumber", newClassroom.SectionNumber, &gocb.UpsertSpecOptions{
				CreatePath: true,
			}),
		}
		result, err := collectionClass.MutateIn(classID, mops, &gocb.MutateInOptions{})

		if err != nil {

			//TODO: Handle (err)
		}
		mops = []gocb.MutateInSpec{
			gocb.ArrayAppendSpec("classrooms", classID, &gocb.ArrayAppendSpecOptions{
				HasMultiple: true,
				CreatePath:  true,
			}),
		}
		result, err = collectionUser.MutateIn(email, mops, &gocb.MutateInOptions{})
		log.Println(result)
		if err != nil {

			//TODO: Handle (err)
		}
	}
	return returnError
}

func GetClassrooms(email string, collectionUser *gocb.Collection, collectionClass *gocb.Collection) ([]Classroom, APIError) {
	var userClasses []string
	var userClassroomsInfo []Classroom
	returnErr := APIError{}
	ops := []gocb.LookupInSpec{
		gocb.GetSpec("classrooms", &gocb.GetSpecOptions{}),
	}
	getResult, err := collectionUser.LookupIn(email, ops, &gocb.LookupInOptions{})

	if err != nil {

		//TODO: Handle (err)
	}

	err = getResult.ContentAt(0, &userClasses)
	if err != nil {

		//TODO: Handle (err)
	}

	for i := 0; i < len(userClasses); i++ {
		ops := []gocb.LookupInSpec{
			gocb.GetSpec("classname", &gocb.GetSpecOptions{}),
			gocb.GetSpec("classnumber", &gocb.GetSpecOptions{}),
			gocb.GetSpec("sectionnumber", &gocb.GetSpecOptions{}),
			gocb.GetSpec("dcordserverid", &gocb.GetSpecOptions{}),
			gocb.GetSpec("studentlist", &gocb.GetSpecOptions{}),
			gocb.GetSpec("talist", &gocb.GetSpecOptions{}),
		}
		getResult, err := collectionClass.LookupIn(userClasses[i], ops, &gocb.LookupInOptions{})

		if err != nil {

			//TODO: Handle (err)
		}

		tempClassroom := Classroom{}
		err = getResult.ContentAt(0, &tempClassroom.ClassName)
		if err != nil {
			//TODO: Handle (err)
		}
		err = getResult.ContentAt(1, &tempClassroom.ClassNumber)
		if err != nil {
			//TODO: Handle (err)
		}
		err = getResult.ContentAt(2, &tempClassroom.SectionNumber)
		if err != nil {
			//TODO: Handle (err)
		}
		err = getResult.ContentAt(3, &tempClassroom.DCordServerID)
		if err != nil {
			//TODO: Handle (err)
		}
		err = getResult.ContentAt(4, &tempClassroom.StudentList)
		if err != nil {
			//TODO: Handle (err)
		}
		err = getResult.ContentAt(5, &tempClassroom.TAList)
		if err != nil {
			//TODO: Handle (err)
		}

		userClassroomsInfo = append(userClassroomsInfo, tempClassroom)

	}

	//TODO: Return array of class info
	return userClassroomsInfo, returnErr

}

func AddStudents() {
	//TODO: Implement
}

func AddTAs() {
	//TODO: Implement
}
