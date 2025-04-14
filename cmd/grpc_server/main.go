package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Ippolid/auth/internal/config"
	"github.com/Ippolid/auth/pkg/auth_v1"
	"github.com/Ippolid/auth/postgres/query"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	auth_v1.UnimplementedAuthV1Server
	db *pgx.Conn
}

// Get ...
func (s *server) Get(_ context.Context, req *auth_v1.GetRequest) (*auth_v1.GetResponse, error) {
	log.Printf("User id: %d", req.GetId())

	ctx := context.Background()
	User, err := query.GetUser(ctx, s.db, int(req.GetId()))
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}
	defer ctx.Done()
	fmt.Printf("User: %+v", User)
	return &auth_v1.GetResponse{
		Note: &auth_v1.UserGet{
			Id: req.GetId(),
			Info: &auth_v1.UserInfo{
				Name:  User.Name,
				Email: User.Email,
			},
			CreatedAt: timestamppb.New(User.CreatedAt),
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Create(_ context.Context, req *auth_v1.CreateRequest) (*auth_v1.CreateResponse, error) {
	//чето делается
	fmt.Printf("name +%v\n", req.Info)
	name := req.GetInfo().GetUser().Name
	email := req.GetInfo().GetUser().Email
	password := req.GetInfo().GetPassword()
	// Преобразуем Role в bool (предполагая, что Role - это enum или int32)
	role := req.GetInfo().GetRole() > 0

	ctx := context.Background()
	id, err := query.InsertUser(ctx, s.db, name, email, password, role)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}
	defer ctx.Done()

	return &auth_v1.CreateResponse{
		Id: int64(id),
	}, nil
}

func (s *server) Update(_ context.Context, req *auth_v1.UpdateRequest) (*emptypb.Empty, error) {
	//чето делается
	fmt.Printf("User id: %d", req.GetId())
	fmt.Printf("name +%v\n", req.Info)
	name := req.Info.Name
	email := req.Info.Email
	ctx := context.Background()
	err := query.UpdateUser(ctx, s.db, int(req.GetId()), name, email)
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}
	defer ctx.Done()

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, req *auth_v1.DeleteRequest) (*emptypb.Empty, error) {
	//чето делается
	fmt.Printf("User id: %d", req.GetId())

	ctx := context.Background()
	err := query.DeleteUser(ctx, s.db, int(req.GetId()))
	if err != nil {
		log.Fatalf("failed to insert user: %v", err)
	}
	defer ctx.Done()

	return &emptypb.Empty{}, nil
}

func main() {
	ctx := context.Background()

	err := config.Load("./.env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dsn, err := config.NewPGConfig()
	fmt.Println(dsn.DSN())
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	con, err := pgx.Connect(ctx, dsn.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func(con *pgx.Conn, ctx context.Context) {
		err := con.Close(ctx)
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(con, ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	auth_v1.RegisterAuthV1Server(s, &server{db: con})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
