package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/Royal17x/go-url-shortener/internal/auth"
	"github.com/Royal17x/go-url-shortener/internal/pb"
	"github.com/Royal17x/go-url-shortener/internal/service"
	"github.com/Royal17x/go-url-shortener/internal/storage"
	"github.com/Royal17x/go-url-shortener/internal/validation"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	auth *service.AuthService
}

type URLServer struct {
	pb.UnimplementedURLServiceServer
	urlService *service.URLService
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if !validation.ValidateUsername(req.Username) {
		return nil, status.Errorf(codes.InvalidArgument, "некорректное имя пользователя")
	}
	if !validation.ValidateEmail(req.Email) {
		return nil, status.Errorf(codes.InvalidArgument, "некорректная почта пользователя")
	}
	if !validation.ValidatePassword(req.Password) {
		return nil, status.Errorf(codes.InvalidArgument, "некорректный пароль пользователя")
	}

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
		return nil, status.Errorf(codes.NotFound, "пользователя с почтой:%s не существует", req.Email)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.Password)); err != nil {
		fmt.Printf("DEBUG логин: сохраненный хэш = %s\n", existingUser.Password)
		fmt.Printf("DEBUG Login: пароль при вводе = %s\n", req.Password)
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	token, err := auth.GenerateToken(existingUser.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка генерации токена:%v", err)
	}
	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (s *URLServer) ShortenURL(ctx context.Context, req *pb.ShortenURLRequest) (*pb.ShortenURLResponse, error) {
	shortCode, err := s.urlService.ShortenURL(ctx, int(req.UserId), req.OriginalUrl)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при сохранении URL: %v", err)
	}
	return &pb.ShortenURLResponse{ShortCode: shortCode}, nil
}

func (s *URLServer) ResolveURL(ctx context.Context, req *pb.ResolveURLRequest) (*pb.ResolveURLResponse, error) {
	original, err := s.urlService.ResolveURL(ctx, req.ShortCode)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "короткий код не найден:%v", err)
	}
	return &pb.ResolveURLResponse{OriginalUrl: original}, nil
}

func main() {
	dsn := os.Getenv("DB_URL")

	store, err := storage.NewDB(dsn)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД:%v", err)
	}
	defer store.DB.Close()
	authService := service.NewAuthService(store)
	urlService := service.NewURLService(store)
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &AuthServer{auth: authService})
	pb.RegisterURLServiceServer(grpcServer, &URLServer{urlService: urlService})
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка в создании listener:%v", err)
	}
	log.Println("gRPC запущен на :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
	}
}
