package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "mongodbtest/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("didnt connect %v", err)
	}
	if err == nil {
		fmt.Println("connection successful")
	}
	defer conn.Close()

	client := pb.NewUserManagerClient(conn)

	fmt.Print(`Choose an option:
1 create user
2 read user
3 update user
4 delete user`)

	var choice int
	_, err1 := fmt.Scan(&choice)
	if err1 != nil {
		log.Fatal(err)
	}

	switch choice {
	case 1:
		var id int32
		_, err2 := fmt.Scan(&id)
		if err2 != nil {
			fmt.Print(err2)
		}

		var name string
		_, err3 := fmt.Scan(&name)
		if err3 != nil {
			fmt.Print(err3)
		}

		var email string
		_, err4 := fmt.Scan(&email)
		if err4 != nil {
			fmt.Print(err4)
		}

		var password string
		_, err5 := fmt.Scan(&password)
		if err5 != nil {
			fmt.Print(err5)
		}

		newuser := pb.User{Id: id, Name: name, Email: email, Password: password}

		createusr, err := client.Create(context.Background(), &pb.CreateRequest{User: &newuser})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(createusr.CreateResponse)

	case 2:

	}

}
