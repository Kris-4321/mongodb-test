package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "mongodbtest/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	connection, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect %v", err)
	}
	defer connection.Close()

	client := pb.NewUserManagerClient(connection)

	var choice int
	for {
		fmt.Print(`Choose an operation
1 Create user
2 get user by id
3 update user
4 delete user
5 list all users
Enter choice:`)
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			CreateUser(client)
		case 2:
			ReadUser(client)
		case 3:
			UpdateUser(client)
		case 4:
			DeleteUser(client)
		case 5:
			ReadallUsers(client)
		default:
			fmt.Print("invalid choice, choose again")
		}

	}

}

func CreateUser(client pb.UserManagerClient) {
	var name, email, password string
	fmt.Print("enter name:")
	fmt.Scanln(&name)
	fmt.Print("enter email")
	fmt.Scanln(&email)
	fmt.Print("enter password")
	fmt.Scanln(&password)

	timeout := 20 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	create, err := client.Create(ctx, &pb.CreateRequest{Name: name, Email: email, Password: password})
	if err != nil {
		fmt.Printf("couldnt create user %v", err)
	}
	fmt.Printf("user created %v", create.Id)
}

func ReadUser(client pb.UserManagerClient) {
	var id string
	fmt.Print("enter id: ")
	fmt.Scanln(&id)

	timeout := 20 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	read, err := client.Read(ctx, &pb.ReadRequest{Id: id})
	if err != nil {
		fmt.Printf("error reading id %v", err)
	}
	fmt.Printf("user: id=%s, name=%s, email=%s, password=%s", read.Id, read.Name, read.Email, read.Password)
}

func UpdateUser(client pb.UserManagerClient) {
	var id string
	fmt.Print("enter id: ")
	fmt.Scanln(&id)

	timeout := 20 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	read, err := client.Read(ctx, &pb.ReadRequest{Id: id})
	if err != nil {
		fmt.Printf("error reading id %v", err)
	}
	fmt.Printf("user: id=%s, name=%s, email=%s, password=%s", read.Id, read.Name, read.Email, read.Password)

	var newname, newemail, newpassword string
	fmt.Print("enter new name: ")
	fmt.Scanln(&newname)
	fmt.Print("enter new email: ")
	fmt.Scanln(&newemail)
	fmt.Print("enter new password: ")
	fmt.Scanln(&newpassword)

	update, err := client.Update(ctx, &pb.UpdateRequest{Id: id, Name: newname, Email: newemail, Password: newpassword})
	if err != nil {
		fmt.Printf("couldnt update user  %v", err)
	}
	fmt.Printf("User updated new info: id=%s, new name=%s, new email=%s, new password=%s", update.Id, update.Name, update.Email, update.Password)
}

func DeleteUser(client pb.UserManagerClient) {
	var id string
	fmt.Print("enter id: ")
	fmt.Scanln(&id)

	timeout := 20 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	delete, err := client.Delete(ctx, &pb.ReadRequest{Id: id})
	if err != nil {
		fmt.Printf("couldnt delete user %v", err)
	}
	fmt.Printf("user with id=%s deleted", delete.Id)
}

func ReadallUsers(client pb.UserManagerClient) {
	timeout := 20 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	readall, err := client.ReadAll(ctx, &pb.ListRequest{})
	if err != nil {
		log.Fatalf("failed to lsit users %v", err)
	}
	fmt.Println("users: ")
	for _, user := range readall.Users {
		fmt.Printf("id=%s, name=%s, email=%s, password=%s", user.Id, user.Name, user.Email, user.Password)
	}
}
