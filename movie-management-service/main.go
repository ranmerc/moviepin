package main

import (
	"database/sql"
	"fmt"
	"movie-management-service/domain"
	"movie-management-service/grpcclient"
	"movie-management-service/handlers"
	"movie-management-service/middleware"
	"movie-management-service/routes"
	"movie-management-service/utils"
	"token-management-service/tokenpb"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	db, err := openDB()
	if err != nil {
		utils.ErrorLogger.Printf("error connecting to DB: %v", err)
		return
	}

	utils.Logger.Println("DB connection is successful")

	defer db.Close()

	tokenServiceConnection, err := grpc.Dial("localhost:4550", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.ErrorLogger.Printf("error connecting to token service grpc: %v", err)
		return
	}
	defer tokenServiceConnection.Close()

	client := tokenpb.NewTokenServiceClient(tokenServiceConnection)

	tokenServiceClient := grpcclient.NewTokenServiceClient(client)
	movieService := domain.NewMovieService(db)

	movieHandler := handlers.NewMovieHandler(movieService, tokenServiceClient)

	server := gin.Default()
	apiRoutes := routes.NewRoutes(movieHandler)
	authMiddleware := middleware.NewAuthMiddleware(tokenServiceClient)

	routes.AttachRoutes(server, authMiddleware, apiRoutes)

	server.Run(":4545")
}
