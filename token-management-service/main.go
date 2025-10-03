package main

import (
	"fmt"
	"movie-management-service/utils"
	"token-management-service/grpchandler"
	"token-management-service/grpcserver"
)

func main() {
	port, err := utils.GetEnv("MOVIEPIN_TOKEN_SERVICE_PORT")
	if err != nil {
		utils.ErrorLogger.Printf("error getting token service port: %v, using default port :4550", err)
		port = "4550"
	}

	handler := grpchandler.Handler{}
	server := grpcserver.NewServer(fmt.Sprintf("localhost:%s", port), handler)

	utils.Logger.Printf("Started gRPC server on port :%s\n", port)
	if err := server.ListenAndServe(); err != nil {
		utils.ErrorLogger.Fatal("Failed to start gRPC server: ", err)
	}
}
