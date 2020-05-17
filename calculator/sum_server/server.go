package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/grpc-go-course/calculator/sumpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	a := req.GetArguments().GetA()
	b := req.GetArguments().GetB()
	result := a + b
	res := sumpb.SumResponse{
		Result: result,
	}

	return &res, nil
}

func isPrime(n int32) bool {
	if n <= 1 {
		return false
	}

	if n == 2 || n == 3 {
		return true
	}

	if n%2 == 0 || n%3 == 0 {
		return false
	}

	for i := int32(5); i*i <= n; i = i + 2 {
		if n%i == 0 || n%i+2 == 0 {
			return false
		}
	}
	return true
}

func (*server) NumberPrimeDecompose(req *sumpb.NumberPrimeDecomposeRequest, stream sumpb.SumService_NumberPrimeDecomposeServer) error {
	a := req.GetArgument().GetNum()
	i := int32(2)
	for a > 1 {
		if isPrime(i) && a%i == 0 {
			a = a / i
			res := sumpb.NumberPrimeDecomposeResponse{
				Result: i,
			}
			stream.Send(&res)
			continue
		}
		i++
		time.Sleep(time.Second)
	}
	return nil
}

func (*server) Average(stream sumpb.SumService_AverageServer) error {
	fmt.Println("Average function was invoked")
	sum := int32(0)
	count := int32(0)

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&sumpb.AverageResponse{
				Result: sum / count,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		n := msg.GetArgument().GetNum()
		sum += n
		count++
	}
}

func (*server) Max(stream sumpb.SumService_MaxServer) error {
	fmt.Println("Max function was invoked")
	max := int32(-1)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		num := req.GetArgument().GetNum()
		if num > max {
			max = num
			err = stream.Send(&sumpb.MaxResponse{
				Result: num,
			})
			if err != nil {
				log.Fatalf("Error while sending server stream: %v", err)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *sumpb.SquareRootRequest) (*sumpb.SquareRootResponse, error) {
	fmt.Println("SquareRoot function was invoked")
	number := req.GetArgument().GetNum()

	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %d", number),
		)
	}
	return &sumpb.SquareRootResponse{
		SquareRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Hello")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	sumpb.RegisterSumServiceServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
