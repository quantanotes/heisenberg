package utils

import "github.com/joho/godotenv"

var Env map[string]string // Environment variables

func init() {
	// Load environment variables from .env
	var err error
	Env, err = godotenv.Read()
	if err != nil {
		panic(err)
	}
}
