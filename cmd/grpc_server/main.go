package main

import (
	"chat/pkg/note_v1"
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct {
	note_v1.UnimplementedChatAPIServer
}

func (s *server) CreateChat(ctx context.Context, req *note_v1.CreateChatRequest) (*note_v1.CreateChatResponse, error) {
	log.Printf("Received CreateChat request: Name=%s", req.Usermanes)
	return &note_v1.CreateChatResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) DeleteChat(ctx context.Context, req *note_v1.DeleteChatRequest) (*note_v1.DeleteChatResponse, error) {
	log.Printf("Received DeleteChat request: ID=%d", req.Id)
	return &note_v1.DeleteChatResponse{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *note_v1.SendMessageRequest) (*note_v1.SendMessageResponse, error) {
	log.Printf("Received SendMessage request: From=%s, Text=%s, At=%s", req.From, req.Text, req.Timestamp)
	return &note_v1.SendMessageResponse{}, nil
}

const (
	grpcport = 50052
)

func main() {
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
