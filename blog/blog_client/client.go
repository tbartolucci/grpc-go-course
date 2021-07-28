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
		AuthorId: "TomB",
		Title: "Reading Blog",
		Content: "Body of the Blog",
	}

	createRes, createErr := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if createErr != nil {
		log.Fatalf("Unexpected error: %v", createErr)
	}
	blogId := createRes.GetBlog().GetId()
	fmt.Printf("Finished creating the Blog: %v\n", createRes)

	// read Blog
	fmt.Println("Testing a read failure")
	_, readErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "asdfasdfasd"})
	if readErr != nil {
		fmt.Printf("Error happened while reading: %v", readErr)
	}

	fmt.Println("Reading the blog back out")
	readBlog, readErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogId})
	if readErr != nil {
		fmt.Printf("Error happened while reading: %v", readErr)
	}

	fmt.Printf("Blog was read: %v\n", readBlog)
}
