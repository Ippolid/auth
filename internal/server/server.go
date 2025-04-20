package server

import (
	"context"
	"fmt"
	"log"

	"github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/Ippolid/auth/postgres/query"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Get возвращает информацию о пользователе по ID
func (s *Server) Get(_ context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	ctx := context.Background()
	user, err := s.db.GetUser(ctx, int(req.GetId()))
	if err != nil {
		log.Printf("failed to get user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}
	fmt.Printf("User: %+v", user)
	return &auth_v1.GetResponse{
		Note: &auth_v1.UserGet{
			Id: req.GetId(),
			Info: &auth_v1.UserInfo{
				Name:  user.User.Name,
				Email: user.User.Email,
			},
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

// Create создает нового пользователя
func (s *Server) Create(_ context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	//чето делается
	name := req.GetInfo().GetUser().Name
	email := req.GetInfo().GetUser().Email
	password := req.GetInfo().GetPassword()
	// Преобразуем Role в bool (предполагая, что Role - это enum или int32)
	role := req.GetInfo().GetRole() > 0

	ctx := context.Background()
	user := query.User{
		Name:  name,
		Email: email,
	}
	id, err := s.db.InsertUser(ctx, user, password, role)
	if err != nil {
		log.Printf("failed to insert user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to insert user: %v", err)
	}
	fmt.Printf("User id: %d", id)

	return &auth_v1.CreateResponse{
		Id: int64(id),
	}, nil
}

// Update обновляет информацию о пользователе
func (s *Server) Update(_ context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	//чето делается
	fmt.Printf("User id: %d", req.GetId())
	fmt.Printf("name +%v\n", req.Info)
	name := req.Info.Name
	email := req.Info.Email
	user := query.User{
		Name:  name,
		Email: email,
	}
	ctx := context.Background()
	err := s.db.UpdateUser(ctx, int(req.GetId()), user)
	if err != nil {
		log.Printf("failed to update user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// Delete удаляет пользователя по ID
func (s *Server) Delete(_ context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	//чето делается
	fmt.Printf("User id: %d", req.GetId())

	ctx := context.Background()
	err := s.db.DeleteUser(ctx, int(req.GetId()))
	if err != nil {
		log.Printf("failed to delete user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &emptypb.Empty{}, nil
}
