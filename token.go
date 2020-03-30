package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func getEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error loading env file")
		//log.Fatal(err)
	}

}

func genToken() {
	jwtSecret := os.Getenv("JWT_SECRET")
}
