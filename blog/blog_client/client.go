package main

import (
	"context"
	"fmt"
	"github.com/tbartolucci/udemy-grpc/blog/blogpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	fmt.Println("Hello I'm a Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial grpc server: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)
	fmt.Printf("Connected to client: %v\n", c)

	fmt.Printf("Creating the Blog\n")
	blog := &blogpb.Blog{
		AuthorId: "Tom",
		Title: "First Blog",
		Content: "Body of the Blog",
	}

	createRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Finished creating the Blog: %v\n", createRes)
}