package main

import (
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
		returnError.Error = true
		returnError.Message = "Token invalid"
	}
	_, err := collectionClass.Upsert(classID, newClassroom, &gocb.UpsertOptions{})
	//log.Println("this is an issue" + result)
	if err != nil {

		returnError.Error = true
		returnError.Message = "Database upsert problem"
	}
	return returnError

}

func CreateClassroomFrontEnd(params graphql.ResolveParams, collectionClass *gocb.Collection, collectionUser *gocb.Collection) APIError {
	newClassroom := Classroom{}
	returnError := APIError{}

	valid, email, classID := DiscordValidateToken(params.Args["token"].(string))

	if valid {

		var newlistTA []TA
		var newlistStudent []Student
		newClassroom = Classroom{
			ClassName:     params.Args["classname"].(string),
			ClassNumber:   params.Args["classnumber"].(string),
			SectionNumber: params.Args["sectionnumber"].(string),
			TAList:        newlistTA,
			StudentList:   newlistStudent,
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
		_, err := collectionClass.MutateIn(classID, mops, &gocb.MutateInOptions{})

		if err != nil {

			//TODO: handle
		}
		mops = []gocb.MutateInSpec{
			gocb.ArrayAppendSpec("classrooms", classID, &gocb.ArrayAppendSpecOptions{
				HasMultiple: true,
				CreatePath:  true,
			}),
		}
		_, err = collectionUser.MutateIn(email, mops, &gocb.MutateInOptions{})

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
