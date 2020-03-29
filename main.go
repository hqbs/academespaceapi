package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/couchbase/gocb"
	"github.com/friendsofgo/graphiql"
	"github.com/graphql-go/graphql"
	"github.com/joho/godotenv"
)

type UserToken struct {
	Token      string `json:"token"`
	ExpireDate int    `json:"expiredate"`
}

type UserPassReset struct {
	URLToken   string `json:"urltoken"`
	ExpireDate int    `json:"expiredate"`
	UserEmail  string `json:"useremail"`
	UserToken  string `json:"usertoken"`
}

type StudentClass struct {
	ClassName     string `json:"classname"`
	CourseNumber  int    `json:"coursenumber,omitempty"`
	CourseSection int    `json:"coursesection,omitempty"`
	Professor     string `json:"professor"`
	ProfEmail     string `json:"profemail"`
	University    string `json:"university"`
	UniversityID  string `json:"universityid"`
}

type User struct {
	FName       string         `json:"fname"`
	LName       string         `json:"lname"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phonenumber"`
	Type        string         `json:"type"`
	ID          string         `json:"id"`
	Password    string         `json:"password"`
	DiscordID   string         `json:"discordid,omitempty"`
	Servers     []DCordServer  `json:"server,omitempty"`
	Token       UserToken      `json:"token,omitempty"`
	PassReset   UserPassReset  `json:"passreset,omitempty"`
	Classes     []StudentClass `json:"classes,omitempty"`
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
	// MutPayloadType := graphql.NewObject(graphql.ObjectConfig{
	// 	Name: "mutpayload",
	// 	Fields: graphql.Fields{
	// 		"Success": &graphql.Field{
	// 			Type: graphql.Boolean,
	// 		},
	// 		"Errors": &graphql.Field{
	// 			Type: graphql.NewList(graphql.String),
	// 		},
	// 		"Token": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 	},
	// })

	StudentClassType := graphql.NewObject(graphql.ObjectConfig{
		Name: "studentclass",
		Fields: graphql.Fields{
			"ClassName": &graphql.Field{
				Type: graphql.String,
			},
			"CourseNumber": &graphql.Field{
				Type: graphql.Int,
			},
			"CourseSection": &graphql.Field{
				Type: graphql.Int,
			},
			"Professor": &graphql.Field{
				Type: graphql.String,
			},
			"ProfEmail": &graphql.Field{
				Type: graphql.String,
			},
			"University": &graphql.Field{
				Type: graphql.String,
			},
			"UniversityID": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	UserTokenType := graphql.NewObject(graphql.ObjectConfig{
		Name: "usertoken",
		Fields: graphql.Fields{
			"Token": &graphql.Field{
				Type: graphql.String,
			},
			"ExpireDate": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	UserPassResetType := graphql.NewObject(graphql.ObjectConfig{
		Name: "userpassreset",
		Fields: graphql.Fields{
			"URLToken": &graphql.Field{
				Type: graphql.String,
			},
			"ExpireDate": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})

	UserType := graphql.NewObject(graphql.ObjectConfig{
		Name: "user",
		Fields: graphql.Fields{
			"FName": &graphql.Field{
				Type: graphql.String,
			},
			"LName": &graphql.Field{
				Type: graphql.String,
			},
			"Email": &graphql.Field{
				Type: graphql.String,
			},
			"PhoneNumber": &graphql.Field{
				Type: graphql.String,
			},
			"Type": &graphql.Field{
				Type: graphql.String,
			},
			"ID": &graphql.Field{
				Type: graphql.String,
			},
			"Password": &graphql.Field{
				Type: graphql.String,
			},
			"DiscordID": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	rootMutation := graphql.ObjectConfig(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: UserType,
				Args: graphql.FieldConfigArgument{
					"FName": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"LName": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"Email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"PhoneNumber": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"Type": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"ID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"Password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"DiscordID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					validUser := ValidateInfo(params)
					success := NewUser(validUser, collection)
					payload := MutationPayload{
						Success: success,
						Errors:  nil,
						Token:   "1234",
					}
					return payload, nil
				},
			},
		},
	})

	rootQuery := graphql.ObjectConfig(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"User": &graphql.Field{
				Type: graphql.NewList(UserType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
		},
	})

	schemaConfig := graphql.SchemaConfig{
		Query:    graphql.NewObject(rootQuery),
		Mutation: graphql.NewObject(rootMutation),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":12345", nil)
	graphiqlHandler, err := graphiql.NewGraphiqlHandler("http://localhost:12345/graphql")
	if err != nil {
		panic(err)
	}

	http.Handle("/graphiql", graphiqlHandler)
	http.ListenAndServe(":4040", nil)
}
