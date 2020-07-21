package controllers

import (
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
	"starter/app/models"
	"starter/app/middlewares"
)

type App struct {
	Router *mux.Router
	DB *gorm.DB
}

// Route List
func (app *App) Routes() {
	app.Router = mux.NewRouter().StrictSlash(true)

	app.Router.Use(middlewares.SetContentTypeJSON)

	app.Router.HandleFunc("/api/register", app.Register).Methods("POST")
	app.Router.HandleFunc("/api/login", app.Login).Methods("POST")

	// Require auth
	ProtectedRoute := app.Router.PathPrefix("/protected").Subrouter()
	ProtectedRoute.Use(middlewares.AuthJwtVerify)

	ProtectedRoute.HandleFunc("/api/users", app.GetAllUsers).Methods("GET")
	ProtectedRoute.HandleFunc("/api/articles", app.CreateArticle).Methods("POST")
	ProtectedRoute.HandleFunc("/api/articles/{id}", app.DeleteArticle).Methods("DELETE")
}

// Init the App
func (app *App) Init(DBHost, DBPort, DBUser, DBName, DBPassword string) {
	var err error

	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DBHost, DBPort, DBUser, DBName, DBPassword)

	app.DB, err = gorm.Open("postgres", DBURI)

	if err != nil {
		fmt.Println("Failed connected to database")
	}

	fmt.Println("Connected to database")

	app.DB.Debug().AutoMigrate(&models.User{}, &models.Article{})
	app.Routes()
}

// Run The Server
func (app *App) RunServer() {
	fmt.Printf("\n Server started at port 8000")
	log.Fatal(http.ListenAndServe(":8000", app.Router))
}