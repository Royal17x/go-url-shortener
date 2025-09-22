package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/Royal17x/go-url-shortener/internal/pb"
	"github.com/Royal17x/go-url-shortener/internal/service"
	"github.com/Royal17x/go-url-shortener/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	auth *service.AuthService
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	uid, err := s.auth.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка регистрации:%v", err)
	}
	return &pb.RegisterResponse{UserID: strconv.Itoa(uid)}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	fmt.Printf("Вызван login: username=%s, email=%s", req.Username, req.Email)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "имя, почта и пароль пользователя необходимы")
	}
	existingUser, err := s.auth.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка БД: %v", err)
	}
	if existingUser == nil {
		return nil, status.Errorf(codes.AlreadyExists, "пользователя с почтой:%s не существует", req.Email)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}
	return &pb.LoginResponse{
		Token: "someJWTtoken",
	}, nil
}

func main() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = "postgres://postgres:33tangoqwe@localhost:5432/urlshortener?sslmode=disable"
	}

	store, err := storage.NewDB(dsn)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД:%v", err)
	}
	defer store.DB.Close()
	authService := service.NewAuthService(store)
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &AuthServer{auth: authService})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка в создании listener:%v", err)
	}
	log.Println("gRPC запущен на :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
	}
}
