package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello i'm client")

	tls := true
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatalf("Failed loading certificates: %v", err)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)

	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created a client: %f\n", c)
	doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBiDiStreaming(c)
	// doUnaryWithDeadline(c, 5*time.Second) //should complete
	// doUnaryWithDeadline(c, 2*time.Second) //should timeout
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("Starting to do a Unary RPC...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bekarys",
			LastName:  "Kuspan",
		},
	}
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit!!!" + statusErr.Message())
			}
		} else {
			fmt.Println("Unexpected err: " + statusErr.Message())

		}
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi streaming RPC...")
	requests := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bekarys",
				LastName:  "Kuspan",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Helmet",
				LastName:  "NSD",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Gavnyuk",
				LastName:  "vonyu4iy",
			},
		},
	}
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Println("sending request: " + req.Greeting.FirstName)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}

			fmt.Printf("Received: %v\n", res.GetResult())
		}

		close(waitc)

	}()

	<-waitc
}
func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bekarys",
				LastName:  "Kuspan",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Helmet",
				LastName:  "NSD",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Gavnyuk",
				LastName:  "vonyu4iy",
			},
		},
	}
	ctx := context.Background()
	stream, err := c.LongGreet(ctx)
	if err != nil {
		log.Fatalf("error while calling client streaming GreetManyTimes RPC: %v", err)
	}

	for _, req := range requests {
		fmt.Println("sending request: " + req.Greeting.FirstName)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while calling client streaming GreetManyTimes RPC: %v", err)
	}

	fmt.Printf("LongGreet response  from server: %v", res)
}
func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")
	ctx := context.Background()
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bekarys",
			LastName:  "Kuspan",
		},
	}

	stream, err := c.GreetManyTimes(ctx, req)
	if err != nil {
		log.Fatal("error while calling server streaming GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("error while calling server streaming GreetManyTimes RPC: %v", err)
		}
		log.Printf("Response form GreetManyTimes: %v", msg.GetResult())
	}

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	ctx := context.Background()
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Bekarys",
			LastName:  "Kuspan",
		},
	}
	res, err := c.Greet(ctx, req)
	if err != nil {
		log.Fatal("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}
