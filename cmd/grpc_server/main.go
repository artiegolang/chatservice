package main

import (
	"chat/pkg/note_v1"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

type server struct {
	note_v1.UnimplementedChatAPIServer
}

func (s *server) CreateChat(ctx context.Context, req *note_v1.CreateChatRequest) (*note_v1.CreateChatResponse, error) {
	log.Printf("Получен запрос CreateChat: Usernames=%v", req.Usernames)

	var id int64
	err := dbPool.QueryRow(ctx, "INSERT INTO chats (usernames) VALUES ($1) RETURNING id", req.Usernames).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании чата: %v", err)
	}

	return &note_v1.CreateChatResponse{
		Id: id,
	}, nil
}

func (s *server) DeleteChat(ctx context.Context, req *note_v1.DeleteChatRequest) (*note_v1.DeleteChatResponse, error) {
	log.Printf("Received DeleteChat request: ID=%d", req.Id)

	_, err := dbPool.Exec(ctx, "DELETE FROM chats WHERE id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("error deleting chat: %v", err)
	}

	return &note_v1.DeleteChatResponse{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *note_v1.SendMessageRequest) (*note_v1.SendMessageResponse, error) {
	log.Printf("Получен запрос SendMessage: From=%s, Text=%s, At=%s", req.From, req.Text, req.Timestamp)

	_, err := dbPool.Exec(ctx, "INSERT INTO messages (chat_id, sender, text, timestamp) VALUES ($1, $2, $3, $4)",
		req.ChatId, req.From, req.Text, req.Timestamp.AsTime())
	if err != nil {
		return nil, fmt.Errorf("Ошибка при отправке сообщения: %v", err)
	}

	return &note_v1.SendMessageResponse{}, nil
}

const (
	grpcport = 50052
)

var dbPool *pgxpool.Pool

func main() {
	err := godotenv.Load("/Users/anastasiapasenko/GolandProjects/chat/.env")
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	dbPool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer dbPool.Close()

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcport))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	note_v1.RegisterChatAPIServer(s, &server{})

	log.Printf("Starting gRPC server on port %d", grpcport)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
