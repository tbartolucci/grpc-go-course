package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var mongoClient *mongo.Client
var collection *mongo.Collection

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

type blogItem struct {
	ID primitive.ObjectID 	`bson:"_id,omitempty"`
	AuthorID string			`bson:"author_id"`
	Content string			`bson:"content"`
	Title string			`bson:"title"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title: blog.GetTitle(),
		Content: blog.GetContent(),
	}

	one, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	oid, ok := one.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal,fmt.Sprintf("Cannot convert to OID"))
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

func main() {
	// if we crash the go code we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile )

	fmt.Println("Bootstrapping Blog Server...")

	// Setting up mongo client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, mongoErr := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if mongoErr != nil {
		log.Fatalf("Failed to connect to Mongodb")
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
		fmt.Println("Safely disconnected from MongoDB")
	}()

	collection = client.Database("grpcdb").Collection("blogs")
	fmt.Println("Connected to MongoDb")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	// Starting listening in a separate thread
	go func() {
		fmt.Println("Blog Service Listening...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to server:  %v", err)
		}
		fmt.Println("Blog Service Started")
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Server End")
}
