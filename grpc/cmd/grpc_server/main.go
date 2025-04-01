package main

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Ippolid/auth/tree/main/grpc/pkg/auth_v1"
)

const grpcPort = 50051

type server struct {
	auth_v1.UnimplementedAuthV1Server
}

// Get ...
func (s *server) Get(ctx context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	log.Printf("Note id: %d", req.GetId())

	return &auth_v1.GetResponse{
		Note: &auth_v1.UserGet{
			Id: req.GetId(),
			Info: &auth_v1.UserInfo{
				Name:  gofakeit.BeerName(),
				Email: gofakeit.Email(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	//чето делается
	fmt.Printf("name +%v\n", req.Info)

	return &auth_v1.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Update(ctx context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	//чето делается
	fmt.Printf("User id: %d", req.GetId())
	fmt.Printf("name +%v\n", req.Info)

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	//чето делается
	fmt.Printf("User id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	auth_v1.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
