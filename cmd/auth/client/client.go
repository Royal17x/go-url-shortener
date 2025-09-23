package main

import (
	"context"
	"log"
	"time"

	"github.com/Royal17x/go-url-shortener/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	urlClient := pb.NewURLServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// registerResp, err := client.Register(ctx, &pb.RegisterRequest{
	// 	Username: "petya",
	// 	Email:    "petyanagibator@example.com",
	// 	Password: "qwe123",
	// })

	// if err != nil {
	// 	log.Fatalf("Не удалось получить ответ на Register Request: %v", err)
	// }

	// log.Printf("Register Response: userID=%v", registerResp.UserID)

	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Username: "petya",
		Email:    "petyanagibator@example.com",
		Password: "qwe123",
	})

	if err != nil {
		log.Fatalf("Не удалось получить ответ на Login Request: %v", err)
	}
	log.Printf("Login Response: token=%v", loginResp.Token)

	shortResp, err := urlClient.ShortenURL(ctx, &pb.ShortenURLRequest{
		UserId:      2,
		OriginalUrl: "https://github.com/Royal17x",
	})
	if err != nil {
		log.Fatalf("не удалось восстановить ссылку:%v", err)
	}
	log.Printf("Short URL:%s", shortResp.ShortCode)

	resolveResp, err := urlClient.ResolveURL(ctx, &pb.ResolveURLRequest{ShortCode: shortResp.ShortCode})
	if err != nil {
		log.Fatalf("не удалось восстановить ссылку:%v", err)
	}
	log.Printf("Original URL:%s", resolveResp.OriginalUrl)
}
