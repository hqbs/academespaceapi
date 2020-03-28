package main

import (
	"fmt"
	"log"
	"os"

	"github.com/couchbase/gocb"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
)

type User struct {
	FName       string        `json:"fname"`
	LName       string        `json:"lname"`
	Email       string        `json:"email"`
	PhoneNumber string        `json:"phonenumber"`
	Type        string        `json:"type"`
	ID          string        `json:"id"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	DiscordID   string        `json:"discordid,omitempty"`
	Servers     []DCordServer `json:"server,omitempty"`
}

type ValidatedUser struct {
	ValidUser User     `json:"validuser"`
	UserValid bool     `json:"uservalid"`
	Errors    []string `json:"errors"`
}

type DCordServer struct {
	Name     string `json:"name"`
	UserType string `json:"usertype"`
	/*
		Under this will have lots of omit if empty
		Profs/Teachers/Owners will need admin info
		Students will need basic info
		TAs will need semi admin/privledged access
	*/
	Role        string `json:"role"`
	ID          string `json:"id"`
	AccessHash  string `json:"accesshash"`
	DisplayName string `json:"displayname"`
}

type MutationPayload struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
	Token   string   `json:"token,omitempty"`

}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading env file")
		log.Fatal(err)
	}
	var (
		dbUser   = os.Getenv("COUCH_USER")
		dbPass   = os.Getenv("COUCH_PASS")
		dbAddr   = os.Getenv("COUCH_ADDR")
		dbBucket = os.Getenv("COUCH_U_BUCKET")
	)

	// Connect to DB
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			string(dbUser),
			string(dbPass),
		},
	}

	cluster, err := gocb.Connect(string(dbAddr), opts)
	if err != nil {
		//TODO: Handle error - will need to be sent through API
		//TODO: Most likely a struct with a payload setup of some sort
		log.Fatal(err)
	}
	bucket := cluster.Bucket(dbBucket)
	collection := bucket.DefaultCollection()
	//TODO: implement


	// GraphQL
	mutPayloadType := graphql.NewObject(graphql.ObjectConfig{
		Name: "mutpayload",
		Fields: graphql.Fields{
			"Success": &graphql.Field{
				Type: graphql.Boolean,
			},
			"Errors": &graphql.Field{
				Type: graphql.String[]
			},
		}
	})
}
