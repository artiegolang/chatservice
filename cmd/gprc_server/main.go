package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"micri/pkg/note_v1"
	"net"
)

const grpcport = 50051

type server struct {
	note_v1.UnimplementedAuthAPIServer
}

func (s *server) GetUser(ctx context.Context, req *note_v1.GetUserRequest) (*note_v1.GetUserResponse, error) {
	log.Printf("Received GetUser request: ID=%d", req.Id)
	return &note_v1.GetUserResponse{
		Id:        gofakeit.Int64(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      0,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *server) CreateUser(ctx context.Context, req *note_v1.CreateUserRequest) (*note_v1.CreateUserResponse, error) {
	log.Printf("Received CreateUser request: Name=%s, Email=%s, Password=%s, PasswordConfirm=%s, Role=%v",
		req.Name, req.Email, req.Password, req.PasswordConfirm, req.Role)
	return &note_v1.CreateUserResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *note_v1.UpdateUserRequest) (*note_v1.UpdateUserResponse, error) {
	log.Printf("Received UpdateUser request: ID=%d, Name=%v, Email=%v", req.Id, req.Name, req.Email)
	return &note_v1.UpdateUserResponse{}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *note_v1.DeleteUserRequest) (*note_v1.DeleteUserResponse, error) {
	log.Printf("Received DeleteUser request: ID=%d", req.Id)
	return &note_v1.DeleteUserResponse{}, nil
}

func main() {
	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcport))
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	note_v1.RegisterAuthAPIServer(s, &server{})

	log.Printf("gRPC server listening on port %d", grpcport)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
