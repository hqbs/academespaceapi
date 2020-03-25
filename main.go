package main

type User struct {
	FName       string        `json:"fname"`
	LName       string        `json:"lname"`
	Email       string        `json:"email"`
	PhoneNumber string        `json:"phonenumber"`
	Type        string        `json:"type"`
	ID          int64         `json:"id"`
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

func main() {
	//TODO: implement
}
