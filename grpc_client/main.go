package main

import (
	"context"
	"fmt"
	"log"
	pb "microservice/proto"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMicroserviceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// var user = pb.User{
	// 	Name: "Zues",
	// }
	// response, err := c.CreateUser(ctx, &pb.CreateUserRequest{Name: user.GetName()})
	// if err != nil {
	// 	log.Fatalf("could not create user: %v", err)
	// }
	// log.Printf("User details: Name: %v ID:%v", response.GetName(), response.GetId())
	params := &pb.GetUsersParams{}
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Print("\nUSER LIST: \n")
	fmt.Printf("r.GetUsers(): %v\n", r.GetUsers())
}
