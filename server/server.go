package main

import (
	pb "mongodbtest/proto"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserManagmentServer struct {
	pb.UnimplementedUserManagerServer
	mongoClient *mongo.Client
}

const (
	dbname = ":"
	port   = ":50051"
)
