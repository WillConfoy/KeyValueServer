package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "cs498.com/KeyValueServer/KeyValueClerk"
)

// const (
// 	defaultName = "world"
// 	defaultKey = "Hello"
// )

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	// name = flag.String("name", defaultName, "Name to greet")
	// key = flag.String("key", defaultKey, "Key to get")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMapServiceClerkClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	arr := []string{"Hi", "this", "is", "a", "cool", "sentence", "there", "isn't", "that", "much", "to", "say", "about", "it", "it", "is", "going", "to", "have", "to", "repeat", "some", "stuff", "but", "hopefully", "it'll", "all", "be", "worth", "it", "in", "the", "end"}

	for _, i := range arr {
		_, err := c.Append(ctx, &pb.AppendRequest{Key: i, Value: i})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
	}

	rGet, errGet := c.Get(ctx, &pb.GetRequest{Key: arr[9]})
	if errGet != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Key: %s --> Value: %s", rGet.GetKey(), rGet.GetValue())
}
