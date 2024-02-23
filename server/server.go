package main

import (
	"context"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"

	pb "mongodbtest/proto"
)

const (
	port           = ":50051"
	dbname         = "mongotestdb"
	collectionname = "users"
	connstring     = "mongodb://localhost:27017"
)

type server struct {
	pb.UnimplementedUserManagerServer
	collection *mongo.Collection
}

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s := grpc.NewServer()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(connstring))
	if err != nil {
		log.Fatalf("failed to connect to mongo %v", err)
	}

	collection := client.Database(dbname).Collection(collectionname)

	pb.RegisterUserManagerServer(s, &server{collection: collection})
	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.UserResponse, error) {
	user := bson.D{
		{Key: "name", Value: req.Name},
		{Key: "email", Value: req.Email},
		{Key: "password", Value: req.Password},
	}
	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	return &pb.UserResponse{
		Id:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func (s *server) Read(ctx context.Context, req *pb.ReadRequest) (*pb.UserResponse, error) {
	ObjectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	var user bson.M
	if err := s.collection.FindOne(ctx, bson.M{"_id": ObjectID}).Decode(&user); err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Name:     user["name"].(string),
		Email:    user["email"].(string),
		Password: user["password"].(string),
	}, nil
}

func (s *server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UserResponse, error) {
	ObjectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"$set": bson.M{
			"name":     req.Name,
			"email":    req.Email,
			"password": req.Password,
		},
	}
	_, err2 := s.collection.UpdateOne(ctx, bson.M{"_id": ObjectID}, update)
	if err2 != nil {
		return nil, err2
	}
	return &pb.UserResponse{Id: req.Id, Name: req.Name, Email: req.Email, Password: req.Password}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.ReadRequest) (*pb.UserResponse, error) {
	ObjectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, err
	}

	_, err2 := s.collection.DeleteOne(ctx, bson.M{"_id": ObjectID})
	if err2 != nil {
		return nil, err2
	}
	return &pb.UserResponse{Id: req.Id}, nil
}

func (s *server) ReadAll(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	var users []*pb.UserResponse
	marker, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer marker.Close(ctx)
	for marker.Next(ctx) {
		var user bson.M
		if err := marker.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &pb.UserResponse{
			Id:       user["_id"].(primitive.ObjectID).Hex(),
			Name:     user["name"].(string),
			Email:    user["email"].(string),
			Password: user["password"].(string),
		})
	}
	return &pb.ListResponse{Users: users}, nil
}
