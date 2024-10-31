package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"micri/pkg/note_v1"
	"time"
)

const (
	address = "localhost:50051"
	noteID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := note_v1.NewAuthAPIClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetUser(ctx, &note_v1.GetUserRequest{Id: noteID})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("User: %v", r)
}
