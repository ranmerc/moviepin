package grpcserver

import (
	"context"
	"testing"
	"token-management-service/grpchandler"
	"token-management-service/tokenpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestServer_ListenAndServe(t *testing.T) {
	port := ":50051"

	// Start the gRPC server in a separate goroutine
	go func() {
		handler := grpchandler.Handler{}
		server := NewServer(port, handler)

		if err := server.ListenAndServe(); err != nil {
			t.Errorf("failed to start gRPC server: %v", err)
		}
	}()

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	service := tokenpb.NewTokenServiceClient(conn)

	// Test the gRPC server methods
	res, err := service.GenerateToken(context.Background(), &tokenpb.GenerateTokenRequest{
		Username: "username",
	})
	if err != nil {
		t.Errorf("failed to call GenerateToken: %v", err)
	}

	_, err = service.VerifyToken(context.Background(), &tokenpb.VerifyTokenRequest{
		Token: res.Token,
	})
	if err != nil {
		t.Errorf("failed to call VerifyToken: %v", err)
	}
}
