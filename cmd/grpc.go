/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"microservice/db"
	pb "microservice/proto"
	"net"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "used to start grpc",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("grpc called")
		StartGrpc()

	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)
}

// needed by grpc server
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
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}
	created_user := &pb.User{}
	err = tx.QueryRow(context.Background(), "INSERT INTO users (name) VALUES ($1) RETURNING *", in.GetName()).Scan(&created_user.Id, &created_user.Name)
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

func (server *Server) GetWallet(ctx context.Context, in *pb.WalletRequest) (*pb.Wallet, error) {

	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}
	wallet := &pb.Wallet{}
	var createdDate time.Time
	err = tx.QueryRow(context.Background(),
		"SELECT id_wallet, created_date, balance, currency, users_id FROM wallet WHERE users_id=$1",
		in.GetId()).Scan(&wallet.Id, &createdDate, &wallet.Balance, &wallet.Currency, &wallet.UserId)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())
	fmt.Println(createdDate)
	wallet.CreatedDate = timestamppb.New(createdDate)
	return wallet, nil
}
func (server *Server) GetBalance(ctx context.Context, in *pb.BalanceRequest) (*pb.Balance, error) {
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}
	balance := &pb.Balance{}
	err = tx.QueryRow(context.Background(),
		"SELECT balance FROM wallet WHERE users_id=$1",
		in.GetUserId()).Scan(&balance.Balance)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}
	tx.Commit(context.Background())
	return balance, nil
}

func StartGrpc() {
	var s *Server = NewServer()
	conn, _ := db.DbConnect()
	defer conn.Close(context.Background())
	s.conn = conn
	if err := s.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
