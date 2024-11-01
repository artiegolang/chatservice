package main

import (
	"chat/pkg/note_v1"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	address = "localhost:50052"
)

func main() {
	// Устанавливаем соединение с gRPC-сервером
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	// Создаем клиента
	c := note_v1.NewChatAPIClient(conn)

	// Устанавливаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// 1. Создаем новый чат
	createChatResp, err := c.CreateChat(ctx, &note_v1.CreateChatRequest{
		Usernames: []string{"user1", "user2"}, // Исправлено с "Usermanes" на "Usernames"
	})
	if err != nil {
		log.Fatalf("Не удалось создать чат: %v", err)
	}
	chatID := createChatResp.Id
	log.Printf("Создан чат с ID: %d", chatID)

	// 2. Отправляем сообщение в чат
	sendMessageReq := &note_v1.SendMessageRequest{
		ChatId:    chatID,
		From:      "user1",
		Text:      "Привет, это тестовое сообщение!",
		Timestamp: timestamppb.Now(),
	}
	_, err = c.SendMessage(ctx, sendMessageReq)
	if err != nil {
		log.Fatalf("Не удалось отправить сообщение: %v", err)
	}
	log.Println("Сообщение отправлено")

	// 3. Удаляем чат
	_, err = c.DeleteChat(ctx, &note_v1.DeleteChatRequest{
		Id: chatID,
	})
	if err != nil {
		log.Fatalf("Не удалось удалить чат: %v", err)
	}
	log.Println("Чат удален")
}
