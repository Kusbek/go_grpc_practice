package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/grpc-go-course/calculator/sumpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello i'm client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)

	}
	defer cc.Close()

	c := sumpb.NewSumServiceClient(cc)
	// fmt.Printf("Created a client: %f\n", c)
	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBiDiStreaming(c)
	doErrorUnary(c, 9)
	doErrorUnary(c, -10)
}

func doErrorUnary(c sumpb.SumServiceClient, number int32) {
	// fmt.Println("Starting to do a error Unary RPC...")
	ctx := context.Background()
	req := &sumpb.SquareRootRequest{
		Argument: &sumpb.Argument{
			Num: number,
		},
	}
	res, err := c.SquareRoot(ctx, req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably have sent negative number!")
			}
		} else {
			log.Fatalf("Big error calling SquareRoot: %v", respErr)
		}
		return
	}
	log.Printf("Response from SquareRoot: %v", res.SquareRoot)
}

func doBiDiStreaming(c sumpb.SumServiceClient) {
	nums := []int32{
		2, 1, 0, 3, 2, 1, 0, 10, 5, 79, 125,
	}
	fmt.Println("Starting to do a BiDi streaming RPC...")
	ctx := context.Background()
	stream, err := c.Max(ctx)
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}
	waitc := make(chan struct{})

	go func() {
		for _, n := range nums {
			fmt.Printf("Sending number:%d\n", n)
			stream.Send(&sumpb.MaxRequest{
				Argument: &sumpb.Argument{
					Num: n,
				},
			})
			time.Sleep(1 * time.Second)
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
func doClientStreaming(c sumpb.SumServiceClient) {
	nums := []int32{
		2, 10, 5, 79, 125,
	}
	fmt.Println("Starting to do a client streaming RPC...")
	ctx := context.Background()
	stream, err := c.Average(ctx)
	if err != nil {
		log.Fatalf("error while calling client streaming Average RPC: %v", err)
	}

	for _, n := range nums {
		fmt.Printf("Sending number: %d\n", n)
		stream.Send(
			&sumpb.AverageRequest{
				Argument: &sumpb.Argument{
					Num: n,
				},
			},
		)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while calling client streaming Average RPC: %v", err)
	}

	fmt.Printf("Average response  from server: %v\n", res.GetResult())
}
func doServerStreaming(c sumpb.SumServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")
	ctx := context.Background()
	req := &sumpb.NumberPrimeDecomposeRequest{
		Argument: &sumpb.Argument{
			Num: 1024,
		},
	}

	stream, err := c.NumberPrimeDecompose(ctx, req)
	if err != nil {
		log.Fatal("error while calling server streaming GreetManyTimes RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("error while calling server streaming NumberPrimeDecompose RPC: %v", err)
		}
		log.Printf("Response form NumberPrimeDecompose: %v", msg.GetResult())
	}

}
func doUnary(c sumpb.SumServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	ctx := context.Background()
	req := &sumpb.SumRequest{
		Arguments: &sumpb.Arguments{
			A: 6,
			B: 15,
		},
	}
	res, err := c.Sum(ctx, req)
	if err != nil {
		log.Fatal("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.Result)
}
