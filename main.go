package main

import (
	"os"

	"github.com/kadzany/frosty/internal"
)

func main() {
	app := internal.App{}

	app.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	app.Run(":8080")
}
