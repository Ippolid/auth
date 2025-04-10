package main

import (
	"context"
	"fmt"
	auth_v2 "github.com/Ippolid/auth/pkg/auth_v1"
	"log"
	"net"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	auth_v2.UnimplementedAuthV1Server
}

// Get ...
func (s *server) Get(_ context.Context, req *auth_v2.GetRequest) (*auth_v2.GetResponse, error) {
	log.Printf("Note id: %d", req.GetId())

	return &auth_v2.GetResponse{
		Note: &auth_v2.UserGet{
			Id: req.GetId(),
			Info: &auth_v2.UserInfo{
				Name:  gofakeit.BeerName(),
				Email: gofakeit.Email(),
			},
			CreatedAt: timestamppb.New(gofakeit.Date()),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Create(_ context.Context, req *auth_v2.CreateRequest) (*auth_v2.CreateResponse, error) {
	//чето делается
	fmt.Printf("name +%v\n", req.Info)

	return &auth_v2.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Update(_ context.Context, req *auth_v2.UpdateRequest) (*emptypb.Empty, error) {
	//чето делается
	fmt.Printf("User id: %d", req.GetId())
	fmt.Printf("name +%v\n", req.Info)

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, req *auth_v2.DeleteRequest) (*emptypb.Empty, error) {
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
	auth_v2.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
