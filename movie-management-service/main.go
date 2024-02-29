package main

import (
	"database/sql"
	"fmt"
	"log"
	"movie-management-service/domain"
	"movie-management-service/handlers"
	"movie-management-service/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func openDB() (*sql.DB, error) {
	var (
		host   = "localhost"
		port   = 5432
		user   = "ranmerc"
		dbname = "moviepin"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	server := gin.Default()

	db, err := openDB()
	if err != nil {
		log.Printf("error connecting DB: %v", err)
		return
	}

	log.Println("DB connection is successful")

	defer db.Close()

	movieService := domain.NewMovieService(db)

	movieHandler := handlers.NewMovieHandler(movieService)
	apiRoutes := routes.NewRoutes(movieHandler)
	routes.AttachRoutes(server, apiRoutes)

	server.Run(":4545")
}