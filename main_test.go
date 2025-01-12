package main_test

import (
	"os"
	"testing"

	"github.com/kadzany/frosty/internal"
)

var app internal.App

func TestMain(m *testing.M) {
	app = internal.App{}
	app.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_TEST_DB_NAME"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_PORT"))

	code := m.Run()
	os.Exit(code)
}
