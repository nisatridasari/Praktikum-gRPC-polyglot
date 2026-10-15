package main

import (
	"context"
	"log"
	"net"

	pb "go-server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUserProfile(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Printf("Request diterima untuk user_id=%d", req.GetUserId())

	if req.GetUserId() == 1 {
		return &pb.UserResponse{
			Name:  "Nisa Tridasari",
			Email: "nisa@universitas.ac.id",
		}, nil
	}

	return nil, status.Error(codes.NotFound, "user tidak ditemukan")
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Gagal membuka port: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{})

	log.Println("gRPC server berjalan di port 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Server gagal: %v", err)
	}
}
