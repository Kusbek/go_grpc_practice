package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello i'm client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v\n", err)

	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// createBlogTest(c)
	// readBlogTest(c)
	// UpdateBlogTest(c)
	// DeleteBlogTest(c)
	ListBlogTest(c)
}

func ListBlogTest(c blogpb.BlogServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")
	ctx := context.Background()

	stream, err := c.ListBlog(ctx, &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling server streaming ListBlog RPC: %v\n", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while calling server streaming ListBlog RPC: %v\n", err.Error())
		}
		log.Printf("Response form NumberPrimeDecompose: %v\n", msg.GetBlog())
	}
}

func DeleteBlogTest(c blogpb.BlogServiceClient) {
	//Delete Blog
	fmt.Println("Deleting the blog")
	_, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: "5ec0fbafc7c51ef5cfc1a8d0",
	})
	if err != nil {
		fmt.Printf("Error happened  while deleting: %v\n", err)
	}

	res, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: "5ec0fbafc7c50ef5cfc1a8d0",
	})
	if err != nil {
		fmt.Printf("Error happened  while deleting: %v\n", err)
	}

	fmt.Printf("Blog was deleted: %v", res)

}

func UpdateBlogTest(c blogpb.BlogServiceClient) {
	//Update Blog
	fmt.Println("Update the blog")
	blog := &blogpb.Blog{
		Id:       "5ec0fbafc7c50ef5cfc1a8d0",
		AuthorId: "UpdatedID",
		Title:    "Changed author from  Bekarys to HeLMeT",
		Content:  "Content content test blog",
	}

	cbRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: blog,
	})

	if err != nil {
		log.Fatalf("Unexpected error while updating: %v\n", err)
	}
	fmt.Printf("Blog has been updated: %v\n", cbRes)
}

func readBlogTest(c blogpb.BlogServiceClient) {
	//Read Blog
	fmt.Println("Reading the blog")
	_, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "5ec0fbafc7c51ef5cfc1a8d0",
	})
	if err != nil {
		fmt.Printf("Error happened  while reading: %v\n", err)
	}

	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "5ec0fbafc7c50ef5cfc1a8d0",
	})
	if err != nil {
		fmt.Printf("Error happened  while reading: %v\n", err)
	}

	fmt.Printf("Blog was read: %v", res)

}

func createBlogTest(c blogpb.BlogServiceClient) {
	//Creating Blog
	fmt.Println("Creating blog")
	blog := &blogpb.Blog{
		AuthorId: "Bekarys",
		Title:    "Test Blog Create",
		Content:  "Content content test blog",
	}
	cbRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})

	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", cbRes)
}
