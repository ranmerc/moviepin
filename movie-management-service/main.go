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
	user, err := utils.GetEnv("MOVIEPIN_DB_USER")
	if err != nil {
		return nil, err
	}

	dbName, err := utils.GetEnv("MOVIEPIN_DB_NAME")
	if err != nil {
		return nil, err
	}

	host, err := utils.GetEnv("MOVIEPIN_DB_HOST")
	if err != nil {
		return nil, err
	}

	port, err := utils.GetEnv("MOVIEPIN_DB_PORT")
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	utils.Logger.Println("DB connection is successful")

	return db, nil
}

func main() {
	db, err := openDB()
	if err != nil {
		utils.ErrorLogger.Printf("error connecting to DB: %v", err)
		return
	}

	defer db.Close()

	tokenServicePort, err := utils.GetEnv("MOVIEPIN_TOKEN_SERVICE_PORT")
	if err != nil {
		utils.ErrorLogger.Printf("error getting token service port: %v, using default port :4550", err)
		tokenServicePort = "4550"
	}

	tokenServiceConnection, err := grpc.NewClient(fmt.Sprintf("localhost:%s", tokenServicePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		utils.ErrorLogger.Printf("error connecting to token service grpc: %v", err)
		return
	}
	defer tokenServiceConnection.Close()

	client := tokenpb.NewTokenServiceClient(tokenServiceConnection)

	tokenServiceClient := grpcclient.NewTokenServiceClient(client)
	movieService := domain.NewMovieService(db)

	movieHandler := handlers.NewMovieHandler(movieService, tokenServiceClient)

	port, err := utils.GetEnv("MOVIEPIN_SERVER_PORT")
	if err != nil {
		utils.ErrorLogger.Printf("error getting server port: %v, using default port :4545", err)
		port = "4545"
	}

	server := gin.Default()
	apiRoutes := routes.NewRoutes(movieHandler)
	authMiddleware := middleware.NewAuthMiddleware(tokenServiceClient)

	routes.AttachRoutes(server, authMiddleware, apiRoutes)

	server.Run(fmt.Sprintf("localhost:%s", port))
}
