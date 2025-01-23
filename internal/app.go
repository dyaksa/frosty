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

func (app *App) Initialize(user, password, dbname, dbhost, dbport string) {
	dbCreds := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, dbhost, dbport)
	fmt.Printf("Connection established...")

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
	wfHandler := WorkflowHandler{DB: app.DB}
	app.Router.HandleFunc("/workflow/node", wfHandler.CreateNode).Methods("POST")
	app.Router.HandleFunc("/workflow/node/{id:[0-9a-fA-F-]+}", wfHandler.GetNode).Methods("GET")
	app.Router.HandleFunc("/workflow/node/{id:[0-9a-fA-F-]+}/relationship", wfHandler.AddRelationship).Methods("POST")
	app.Router.HandleFunc("/workflow/{id:[0-9a-fA-F-]+}/execute", wfHandler.ExecuteWorkflow).Methods("POST")
	app.Router.HandleFunc("/workflow", wfHandler.CreateWorkflow).Methods("POST")
}
