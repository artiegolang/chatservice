package main

import (
	"chat/pkg/note_v1"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	address = "localhost:50052"
	noteID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := note_v1.NewChatAPIClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.CreateChat(ctx, &note_v1.CreateChatRequest{Usermanes: []string{"user1", "user2"}})
	if err != nil {
		log.Fatalf("could not create chat: %v", err)
	}

	log.Printf("Chat: %v", r)
}
