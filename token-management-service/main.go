package main

import (
	"movie-management-service/utils"
	"token-management-service/grpchandler"
	"token-management-service/grpcserver"
)

func main() {
	handler := grpchandler.Handler{}
	server := grpcserver.NewServer(":4550", handler)

	utils.Logger.Println("Started gRPC server on port :4550")
	if err := server.ListenAndServe(); err != nil {
		utils.ErrorLogger.Fatal("Failed to start gRPC server: ", err)
	}
}
