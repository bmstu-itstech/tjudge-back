package jwtauth

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	_ = godotenv.Load()
	secret = os.Getenv("SECRET_KEY")
)
