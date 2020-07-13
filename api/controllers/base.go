package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/MrWitold/fetcher-api/api/models"
)

// Server main engine of api
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
	W      *Work
}

// Initialize prepare api to work
func (server *Server) Initialize(DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	server.DB, err = gorm.Open("postgres", DBURL)
	if err != nil {
		log.Fatal("Cannot connect to database - This is the error:", err)
	}

	fmt.Printf("We are connected to the database \n")

	server.DB.Debug().AutoMigrate(&models.Link{}, &models.History{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()

	server.setupWorkers()
	server.initializeWorker()
	server.initializeJobPlanner()
}

// Run starts api server
func (server *Server) Run() {
	s := http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        server.Router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Starting server on port 8080")

	err := s.ListenAndServe()
	if err != nil {
		log.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
