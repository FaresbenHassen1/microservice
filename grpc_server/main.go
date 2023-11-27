package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	pb "microservice/proto"
	"net"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8081, "The server port")
)

func NewServer() *Server {
	return &Server{}
}

type Server struct {
	conn *pgx.Conn
	pb.UnimplementedMicroserviceServer
}

func (server *Server) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMicroserviceServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (server *Server) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
	created_user := &pb.User{Name: in.GetName()}
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}
	_, err = tx.Exec(context.Background(), "INSERT INTO users (name) VALUES ($1)", created_user.Name)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())
	return created_user, nil
}

func (server *Server) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UsersList, error) {
	var users *pb.UsersList = &pb.UsersList{}
	rows, err := server.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return nil, err
		}
		users.Users = append(users.Users, &user)
	}
	return users, nil
}

func main() {
	db_url := "postgres://postgres:postgres@localhost:5432/goproject"
	var s *Server = NewServer()
	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())
	s.conn = conn
	if err := s.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
