package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/razasayed/protoapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RandomService struct {
	protoapi.UnimplementedRandomServer
}

// Implement GetDate method

func (RandomService) GetDate(ctx context.Context, r *protoapi.GetDateRequest) (*protoapi.GetDateResponse, error) {
	currentTime := time.Now()
	response := &protoapi.GetDateResponse{Value: currentTime.String()}
	return response, nil
}

// Implement GetRandomInt method

var min = 0
var max = 100

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func (RandomService) GetRandomInt(ctx context.Context, r *protoapi.GetRandomIntRequest) (*protoapi.GetRandomIntResponse, error) {
	rand.NewSource(r.GetSeed())
	place := r.GetPlace()
	temp := random(min, max)
	for {
		place--
		if place <= 0 {
			break
		}
		temp = random(min, max)
	}

	response := &protoapi.GetRandomIntResponse{
		Value: int64(temp),
	}
	return response, nil
}

// Implement GetRandomPass method

func getString(len int64) string {
	temp := ""
	startChar := "!"
	var i int64 = 1
	for {
		// For getting valid ASCII characters
		myRand := random(0, 94)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == len {
			break
		}
		i++
	}
	return temp
}

func (RandomService) GetRandomPass(ctx context.Context, r *protoapi.GetRandomPassRequest) (*protoapi.GetRandomPassResponse, error) {
	rand.NewSource(r.GetSeed())
	temp := getString(r.GetLength())
	response := &protoapi.GetRandomPassResponse{
		Password: temp,
	}
	return response, nil
}

var port = ":8080"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Using default port:", port)
	} else {
		port = os.Args[1]
	}

	// Create a new gRPC server and register the service with it
	server := grpc.NewServer()
	var randomService RandomService
	protoapi.RegisterRandomServer(server, randomService)
	/*
		Below is optional and is only used for testing and debugging as it allows clients to discover servers API  without needing the
		protobuf definition files.
	*/
	reflection.Register(server)

	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Serving requests...")
	server.Serve(listen)
}
