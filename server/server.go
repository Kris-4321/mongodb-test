package main

import (
	"context"
	"fmt"
	"log"
	pb "mongodbtest/proto"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	port           = ":50051"
	dbname         = "mongotestdb"
	collectionname = "users"
)

type UserManagmentServer struct {
	pb.UnimplementedUserManagerServer
	mongoClient *mongo.Client
}

func (s *UserManagmentServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.UserResponse, error) {
	collection := s.mongoClient.Database(dbname).Collection(collectionname)
	result, err := collection.InsertOne(ctx, bson.M{"name": in.GetName(), "email": in.GetEmail(), "password": in.GetPassword()})
	if err != nil {
		log.Printf("failed to add user into mongo %v", err)
		return nil, err
	}

	oid, err2 := result.InsertedID.(primitive.ObjectID)
	if !err2 {
		log.Printf("failed to add id %v", err2)
		return nil, fmt.Errorf("failed to convert id to ojectid %v", err2)
	}

	return &pb.UserResponse{Id: oid.Hex(), Name: in.GetName(), Email: in.GetEmail(), Password: in.GetPassword()}, nil
}

func (s *UserManagmentServer) Read(ctx context.Context, in *pb.ReadRequest) (*pb.UserResponse, error) {
	collection := s.mongoClient.Database(dbname).Collection(collectionname)
	ID, err := primitive.ObjectIDFromHex(in.GetId())
	if err != nil {
		return nil, err
	}
	var user bson.M
	if err := collection.FindOne(ctx, bson.M{"_id": ID}).Decode(&user); err != nil {
		return nil, err
	}
	return &pb.UserResponse{Id: in.GetId(), Name: user["name"].(string), Email: user["email"].(string), Password: user["password"].(string)}, nil
}

func (s *UserManagmentServer) Readall(ctx context.Context, in *emptypb.Empty) (*pb.ReadAllResponse, error) {
	collection := s.mongoClient.Database(dbname).Collection(collectionname)
	cursor, err := collection.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []*pb.UserResponse
	for cursor.Next(ctx) {
		var user bson.M
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}
		id := user["_id"].(primitive.ObjectID).Hex()
		users = append(users, &pb.UserResponse{Id: id, Name: user["name"].(string), Email: user["email"].(string), Password: user["password"].(string)})

	}
	return &pb.ReadAllResponse{Users: users}, nil
}

func main() {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("couldnt connect to mongo %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("couldnt listen %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserManagerServer(s, &UserManagmentServer{mongoClient: mongoClient})
	log.Printf("server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
