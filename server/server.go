package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "mongodbtest/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	dbname        = "mongotestdb"
	collctionName = "users"
	port          = ":50051"
)

type UserManagmentServer struct {
	pb.UnimplementedUserManagerServer
}

type Userserver struct {
	idserver       int    `bson:"id"`
	nameserver     string `bson:"name"`
	emailserver    string `bson:"email"`
	passwordserver string `bson:"password"`
}

func MongoConnection() (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("conntected to mongodb")
	return client, nil
}

func createUser(client *mongo.Client, databaseName, collectionName string, user Userserver) (*mongo.InsertOneResult, error) {
	collection := client.Database(databaseName).Collection(collectionName)
	result, err := collection.InsertOne(context.Background(), user)
	fmt.Println("user created")
	return result, err
}
func (s *UserManagmentServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	newuser := Userserver{idserver: int(req.User.Id), nameserver: req.User.Email, emailserver: req.User.Email, passwordserver: req.User.Password}
	client, err := MongoConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	_, err = createUser(client, dbname, collctionName, newuser)
	response := "user created!"
	return &pb.CreateResponse{Response: response}, nil
}

func readUser(client *mongo.Client, databaseName, collectionName string, filter interface{}) ([]Userserver, error) {
	collection := client.Database(databaseName).Collection(collectionName)
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	var users []Userserver
	for cur.Next(context.Background()) {
		var user Userserver
		if err := cur.Decode(&user); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	return users, nil
}
func (s *UserManagmentServer) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	filter := bson.M{"id": req.Id}
	client, err := MongoConnection()
	if err != nil {
		log.Fatalln(err)
	}
	users, err := readUser(client, dbname, collctionName, filter)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func updateUser(client *mongo.Client, databaseName, collectionName string, filter, update interface{}) (*mongo.UpdateResult, error) {
	collection := client.Database(databaseName).Collection(collectionName)
	result, err := collection.UpdateMany(context.Background(), filter, update)
	return result, err
}
func (s *UserManagmentServer) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	filter := bson.M{"id": req.Id}
	client, err := MongoConnection()
	if err != nil {
		log.Fatalln(err)
	}
	update := bson.M{"set:": bson.M{"email": req.User.Email}}
	_, err = updateUser(client, dbname, collctionName, filter, update)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func deleterUser(client *mongo.Client, databaseName, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	collection := client.Database(databaseName).Collection(collectionName)
	result, err := collection.DeleteMany(context.Background(), filter)
	return result, err
}
func (s *UserManagmentServer) dele(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	filter := bson.M{"id": req.Id}
	client, err := MongoConnection()
	if err != nil {
		log.Fatalln(err)
	}
	_, err = deleterUser(client, dbname, collctionName, filter)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserManagerServer(s, &UserManagmentServer{})
	log.Printf("server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server %v", err)
	}

	//mongodb

}
