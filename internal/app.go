package internal

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialize(user, password, dbname string) {
	dbCreds := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
	fmt.Printf("Connection %s \r\n", dbCreds)

	var err error
	app.DB, err = sql.Open("postgres", dbCreds)
	if err != nil {
		log.Fatal(err)
	}

	app.Router = mux.NewRouter()
	app.initializeRoutes()
}

func (app *App) Run(addr string) {
	fmt.Printf("Listening on %s", addr)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(headers, methods, origins)(app.Router)))
}

func (app *App) initializeRoutes() {
	// adminHandler := handler.AdminHandler{DB: app.DB}
	// app.Router.HandleFunc("/tipeloket", adminHandler.CreateTipeLoket).Methods("POST")
	// app.Router.HandleFunc("/tipelokets", adminHandler.GetAllTipeLoket).Methods("GET")
	// app.Router.HandleFunc("/tipeloket/{id:[0-9]+}", adminHandler.UpdateTipeLoket).Methods("PUT")
	// app.Router.HandleFunc("/loket", adminHandler.CreateLoket).Methods("POST")
	// app.Router.HandleFunc("/lokets", adminHandler.GetAllLoket).Methods("GET")
	// app.Router.HandleFunc("/loket/{id:[0-9]+}", adminHandler.UpdateLoket).Methods("PUT")
	// app.Router.HandleFunc("/loket/{id:[0-9]+}", adminHandler.DeleteLoket).Methods("DELETE")
	// app.Router.HandleFunc("/tipeloket/{id:[0-9]+}", adminHandler.DeleteTipeLoket).Methods("DELETE")

	// kioskHandler := handler.KioskHandler{DB: app.DB}
	// app.Router.HandleFunc("/tiket/issue/{tipeloketid:[0-9]+}", kioskHandler.IssueTiket).Methods("GET")

	// loketHandler := handler.LoketHandler{DB: app.DB}
	// app.Router.HandleFunc("/loket/{id:[0-9]+}/call", loketHandler.CallTiket).Methods("GET")
	// app.Router.HandleFunc("/loket/{id:[0-9]+}/recall", loketHandler.RecallTiket).Methods("GET")
}
