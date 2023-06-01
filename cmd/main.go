package main

import (
	"github.com/hramov/go-excel-parser/internal"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	godotenv.Load(".env")
	internal.NewServer(os.Getenv("PORT"))
}
