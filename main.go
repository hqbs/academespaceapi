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
	"github.com/joho/godotenv"
)

/* API Structs */

type APIError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type DiscordAPIError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type ValidatedUser struct {
	ValidUser User     `json:"validuser"`
	UserValid bool     `json:"uservalid"`
	Errors    []string `json:"errors"`
}

type MutationPayload struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
	Token   string   `json:"token,omitempty"`
}

/* API Structs End */

/* Token Generation Structs */

type UserToken struct {
	Token      string `json:"token"`
	ExpireDate int64  `json:"expiredate"`
}

type UserPassReset struct {
	URLToken   string `json:"urltoken"`
	ExpireDate int64  `json:"expiredate"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type DiscordClaims struct {
	Email       string `json:"email"`
	ClassroomID string `json:"classroomid"`
	jwt.StandardClaims
}

type DiscordConnectToken struct {
	Token      string `json:"token"`
	ExpireDate int64  `json:"expiredate"`
}

/* Token Generation Structs End */

/* User Structs */

type User struct {
	FName           string              `json:"fname"`
	LName           string              `json:"lname"`
	Email           string              `json:"email"`
	PhoneNumber     string              `json:"phonenumber"`
	Type            string              `json:"type"`
	ID              string              `json:"id"`
	Password        string              `json:"password"`
	DiscordID       string              `json:"discordid,omitempty"`
	Token           UserToken           `json:"token,omitempty"`
	PassReset       UserPassReset       `json:"passreset,omitempty"`
	Classrooms      []UserClassroom     `json:"classrooms,omitempty"`
	ConnectionToken DiscordConnectToken `json:"connectiontoken,omitempty"`
}

type UserClassroom struct {
	CRID          string `json:"crid"`
	Professor     string `json:"professor"`
	ClassName     string `json:"classname"`
	ClassNumber   string `json:"classnumber"`
	SectionNumber string `json:"sectionnumber"`
	JoinCode      string `json:"joincode"`
}

type ModifyUser struct {
	ModifyField string `json:"modifyfield"`
	Value       string `json:"value"`
	UserToken   string `json:"token"`
	Email       string `json:"email"`
}

/* User Structs End */

/* Classroom Structs */

type Classroom struct {
	CRID              string    `json:"crid"`
	University        string    `json:"university,omitempty"`
	Professor         string    `json:"professor,omitempty"`
	ProfessorEmail    string    `json:"professoremail"`
	ClassName         string    `json:"classname,omitempty"`
	ClassNumber       string    `json:"classnumber,omitempty"`
	SectionNumber     string    `json:"sectionnumber,omitempty"`
	ProfessorDCordID  string    `json:"professordcordid"`
	AllEmails         bool      `json:"allemails,omitempty"`
	ApprovedEmails    []string  `json"approvedemails,omitempty"`
	JoinCodeServer    string    `json:"joincodeserver,omitempty"`
	StudentList       []Student `json:"studentlist,omitempty"`
	TAList            []TA      `json:"talist,omitempty"`
	DCordServerID     string    `json:"dcordserverid"`
	DCordConnected    bool      `json:"dcordconnected"`
	FrontEndConnected bool      `json:"frontendconnected,omitempty"`
}

type Student struct {
	StudentEmail string `json:"studentemail"`
	StudentName  string `json:"studentname"`
	JoinCode     string `json:"joincode"`
	DCordID      string `json:"dcordid"`
	DCordNick    string `json:"dcordnick"`
}

type TA struct {
	TAEmail   string `json:"taemail"`
	TAName    string `json:"taname"`
	JoinCode  string `json:"joincode"`
	DCordID   string `json:"dcordid"`
	DCordNick string `json:"dcordnick"`
}

/* Classroom Structs End */

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading env file")
		//log.Fatal(err)
	}
	var (
		dbUser        = os.Getenv("COUCH_USER")
		dbPass        = os.Getenv("COUCH_PASS")
		dbAddr        = os.Getenv("COUCH_ADDR")
		dbBucket      = os.Getenv("COUCH_U_BUCKET")
		dbClassBucket = os.Getenv("COUCH_CLASS_BUCKET")
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
	collectionUser := bucket.DefaultCollection()
	bucketClass := cluster.Bucket(dbClassBucket)
	collectionClass := bucketClass.DefaultCollection()
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

	// Error   bool   `json:"error"`
	// Message string `json:"message"`

	rootMutation := graphql.ObjectConfig(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type:        MutPayloadType,
				Description: "Create a user!",
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
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					validUser := ValidateInfo(params)
					returnToken, createUserErrors := NewUser(validUser, collectionUser)
					createUserErrors.Token = returnToken.Token

					return createUserErrors, nil
				},
			},
			"gendiscordtoken": &graphql.Field{
				Type:        MutPayloadType,
				Description: "Generate Discord Connection Token",
				Args: graphql.FieldConfigArgument{
					"email": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					returnPayload := MutationPayload{}
					dcordToken, apiError := DiscordTokenGen(params.Args["email"].(string), collectionUser)
					returnPayload.Token = dcordToken.Token
					if apiError.Error {
						returnPayload.Success = false
						returnPayload.Errors = append(returnPayload.Errors, apiError.Message)
					}
					return returnPayload, nil
				},
			},
			"createserverd": &graphql.Field{
				Type:        MutPayloadType,
				Description: "Create discord server from bot",
				Args: graphql.FieldConfigArgument{
					"token": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"professordcordid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"dcordserverid": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					returnPayload := MutationPayload{}
					returnPayload.Success = true
					returnPayload.Token = params.Args["token"].(string)
					apiError := CreateClassroomDiscord(params, collectionClass)
					if apiError.Error {
						returnPayload.Success = false
						returnPayload.Errors = append(returnPayload.Errors, apiError.Message)
					}

					return returnPayload, nil
				},
			},
			"createclassroomfront": &graphql.Field{
				Type:        MutPayloadType,
				Description: "Create classroom front end",
				Args: graphql.FieldConfigArgument{
					"token": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"classname": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"classnumber": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"sectionnumber": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					returnPayload := MutationPayload{}
					apiError := CreateClassroomFrontEnd(params, collectionClass)
					returnPayload.Token = params.Args["token"].(string)
					if apiError.Error {
						returnPayload.Success = false
						returnPayload.Errors = append(returnPayload.Errors, apiError.Message)
					}
					return returnPayload, nil
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
					returnVal, _ := UserExist(params.Args["email"].(string), collectionUser)

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
					getResult, err := collectionUser.LookupIn(params.Args["email"].(string), ops, &gocb.LookupInOptions{})
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
					returnVal, errorReturn := Login(params.Args["email"].(string), params.Args["password"].(string), params.Args["token"].(string), collectionUser)
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
