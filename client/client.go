package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "mongodbtest/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("couldnt connect %v", err)
	}
	if err == nil {
		fmt.Println("connected successfully")
	}

	defer conn.Close()
	client := pb.NewUserManagerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Print(`
	select an option:
	1 create new user
	2 Search for user
	3 Update user
	4 delete user
	5 list all users
	choice: `)
	var choice int
	_, err3 := fmt.Scan(&choice)
	if err3 != nil {
		log.Fatalf("error : %v", err3)
	}

	switch choice {
	case 1:
		createUser(client, ctx)

	case 2:
		readUser(client, ctx)

	case 3:
		updateUser(client, ctx)

	case 4:
		deleteUser(client, ctx)

	case 5:
		readallUsers(client, ctx)
	}
}

func createUser(client pb.UserManagerClient, ctx context.Context) {
	var name, email, password string

	fmt.Print("enter name:")
	_, err := fmt.Scanln(&name)
	if err != nil {
		log.Fatalf("failed reading name %v", err)
	}

	fmt.Print("enter email:")
	_, err2 := fmt.Scanln(&email)
	if err2 != nil {
		log.Fatalf("failed reading email %v", err2)
	}

	fmt.Print("enter password:")
	_, err3 := fmt.Scanln(&password)
	if err3 != nil {
		log.Fatalf("failed reading password %v", err3)
	}

	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	user, err4 := client.Create(ctx, &pb.CreateRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})

	if err4 != nil {
		log.Printf("failed adding user %v", err4)
	}
	log.Printf("added user %v", user)
}

func readUser(client pb.UserManagerClient, ctx context.Context) {
	var id string
	fmt.Print("enter id:")
	_, err := fmt.Scanln(&id)
	if err != nil {
		log.Fatalf("failed reading id %v", err)
	}

	user, err := client.Read(ctx, &pb.ReadRequest{
		Id: id,
	})
	if err != nil {
		log.Fatal("couldnt find user", err)
	}
	fmt.Printf("Read user: %v", user)
}

func updateUser(client pb.UserManagerClient, ctx context.Context) {
	var id, name, email, password string

	fmt.Print("enter id:")
	_, err0 := fmt.Scanln(&id)
	if err0 != nil {
		log.Fatalf("failed reading id %v", err0)
	}

	fmt.Print("enter name:")
	_, err := fmt.Scanln(&name)
	if err != nil {
		log.Fatalf("failed reading name %v", err)
	}

	fmt.Print("enter email:")
	_, err2 := fmt.Scanln(&email)
	if err2 != nil {
		log.Fatalf("failed reading email %v", err2)
	}

	fmt.Print("enter password:")
	_, err3 := fmt.Scanln(&password)
	if err3 != nil {
		log.Fatalf("failed reading password %v", err3)
	}

	user, err4 := client.Update(ctx, &pb.UpdateRequest{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err4 != nil {
		log.Fatalf("couldnt update user %v", err4)
	}
	log.Printf("user updated %v", user)
}

func deleteUser(client pb.UserManagerClient, ctx context.Context) {
	var id string
	fmt.Print("enter id:")
	_, err0 := fmt.Scanln(&id)
	if err0 != nil {
		log.Fatalf("failed reading id %v", err0)
	}

	_, err := client.Delete(ctx, &pb.ReadRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("couldnt delete user with id %v", err)
	}
	log.Println("deleted user:")
}

func readallUsers(client pb.UserManagerClient, ctx context.Context) {
	users, err := client.ReadAll(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("couldnt list users %v", err)
	}
	log.Printf("Users: %v", users)
}
