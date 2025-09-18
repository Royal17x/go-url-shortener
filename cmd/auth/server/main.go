package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Royal17x/go-url-shortener/internal/pb"
	"google.golang.org/grpc"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	fmt.Printf("Вызван register: username=%s, email=%s.\n", req.Username, req.Email)
	return &pb.RegisterResponse{
		UserID: "123",
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	fmt.Printf("Вызван login: username=%s, email=%s", req.Username, req.Email)
	return &pb.LoginResponse{
		Token: "someJWTtoken",
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка в создании listener:%v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &AuthServer{})

	log.Println("gRPC запущен на :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
	}
}
