package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/couchbase/gocb"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
)

type APIError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type UserToken struct {
	Token      string `json:"token"`
	ExpireDate int64  `json:"expiredate"`
}

type UserPassReset struct {
	URLToken   string `json:"urltoken"`
	ExpireDate int64  `json:"expiredate"`
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

type ModifyUser struct {
	ModifyField string `json:"modifyfield"`
	Value       string `json:"value"`
	UserToken   string `json:"token"`
	Email       string `json:"email"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("error loading env file")
	// 	//log.Fatal(err)
	// }
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
		//log.Fatal(err)
	}
	bucket := cluster.Bucket(dbBucket)
	collection := bucket.DefaultCollection()
	//TODO: implement

	// GraphQL
	MutPayloadType := graphql.NewObject(graphql.ObjectConfig{
		Name: "mutpayload",
		Fields: graphql.Fields{
			"success": &graphql.Field{
				Type: graphql.Boolean,
			},
			"errors": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
			"token": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// StudentClassType := graphql.NewObject(graphql.ObjectConfig{
	// 	Name: "studentclass",
	// 	Fields: graphql.Fields{
	// 		"ClassName": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"CourseNumber": &graphql.Field{
	// 			Type: graphql.Int,
	// 		},
	// 		"CourseSection": &graphql.Field{
	// 			Type: graphql.Int,
	// 		},
	// 		"Professor": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"ProfEmail": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"University": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"UniversityID": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 	},
	// })

	// UserTokenType := graphql.NewObject(graphql.ObjectConfig{
	// 	Name: "usertoken",
	// 	Fields: graphql.Fields{
	// 		"Token": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"ExpireDate": &graphql.Field{
	// 			Type: graphql.Int,
	// 		},
	// 	},
	// })

	// UserPassResetType := graphql.NewObject(graphql.ObjectConfig{
	// 	Name: "userpassreset",
	// 	Fields: graphql.Fields{
	// 		"URLToken": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"ExpireDate": &graphql.Field{
	// 			Type: graphql.Int,
	// 		},
	// 	},
	// })

	// UserType := graphql.NewObject(graphql.ObjectConfig{
	// 	Name: "user",
	// 	Fields: graphql.Fields{
	// 		"FName": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"LName": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"Email": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"PhoneNumber": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"Type": &graphql.Field{
	// 			Type: graphql.String,

	// 		},
	// 		"ID": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"Password": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 		"DiscordID": &graphql.Field{
	// 			Type: graphql.String,
	// 		},
	// 	},
	// })

	rootMutation := graphql.ObjectConfig(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type:        MutPayloadType,
				Description: "Update existing todo, mark it done or not done",
				Args: graphql.FieldConfigArgument{
					"fname": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"lname": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"phonenumber": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"type": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"discordid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					validUser := ValidateInfo(params)
					returnToken, createUserErrors := NewUser(validUser, collection)
					createUserErrors.Token = returnToken.Token

					return createUserErrors, nil
				},
			},
		},
	})

	rootQuery := graphql.ObjectConfig(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"userExists": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Check to see if a user exists",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					returnVal, _ := UserExist(params.Args["email"].(string), collection)

					return returnVal, nil
				},
			},
			"validateUserToken": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Checks user JWT token for validity",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"currenttoken": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var id string
					var userToken UserToken
					ops := []gocb.LookupInSpec{
						gocb.GetSpec("id", &gocb.GetSpecOptions{}),
						gocb.GetSpec("token", &gocb.GetSpecOptions{}),
					}
					getResult, err := collection.LookupIn(params.Args["email"].(string), ops, &gocb.LookupInOptions{})
					if err != nil {

						return false, nil
					}

					err = getResult.ContentAt(0, &id)
					if err != nil {

						return false, nil
					}
					err = getResult.ContentAt(1, &userToken)
					if err != nil {

						return false, nil
					}
					return ValidateToken(params.Args["currenttoken"].(string), userToken, id), nil
				},
			},
			"login": &graphql.Field{
				Type:        MutPayloadType,
				Description: "Login!",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"password": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"token": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					returnVal, errorReturn := Login(params.Args["email"].(string), params.Args["password"].(string), params.Args["token"].(string), collection)
					var errors []string
					errors = append(errors, errorReturn.Message)
					returnPayload := MutationPayload{

						Errors: errors,
						Token:  returnVal.Token,
					}
					if errorReturn.Error {
						returnPayload.Success = false
					} else {
						returnPayload.Success = true
					}
					return returnPayload, nil
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
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)

	})
	http.ListenAndServe(":4000", nil)
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
