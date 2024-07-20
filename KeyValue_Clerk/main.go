package main

import (
	"context"
	"flag"
	"log"
	"time"
	"net"
	"fmt"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "cs498.com/KeyValueServer/KeyValue"
	clerk "cs498.com/KeyValueServer/KeyValueClerk"
)


// const (
// 	defaultName = "world"
// 	defaultKey = "Hello"
// )

type server struct {
	clerk.UnimplementedMapServiceClerkServer
}

var (
	connaddr = flag.String("connaddr", "localhost:50052", "the address to connect to")
	port = flag.Int("port", 50051, "The server port")
	conn *grpc.ClientConn
	c pb.MapServiceClient
	myCtx context.Context
	cancel context.CancelFunc
	err error
	outCh = make(chan []string)
	order = [][]string{}
	// name = flag.String("name", defaultName, "Name to greet")
	// key = flag.String("key", defaultKey, "Key to get")
)

// func connect() {
// 	conn, err := grpc.Dial(*connaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	c := pb.NewMapServiceClient(conn)

// 	// Contact the server and print out its response.
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()
// }

// func serve() {
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}
// 	s := grpc.NewServer()
// 	clerk.RegisterMapServiceClerkServer(s, &server{})
// 	log.Printf("server listening at %v", lis.Addr())
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }

func handleOrder() {
	for {
		out := <- outCh
		order = append(order, out)
		if rand.Intn(100) > 95 {fmt.Println(order)}
	}
}

func (s *server) Get(ctx context.Context, in *clerk.GetRequest) (*clerk.GetResponse, error){
	conn, err = grpc.Dial(*connaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewMapServiceClient(conn)

	// Contact the server and print out its response.
	myCtx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rGet, errGet := c.Get(myCtx, &pb.GetRequest{Key: in.GetKey()})
	log.Printf("Passing %v onto server as Put", in.GetKey())
	outCh <- []string{"get", in.GetKey()}
	return &clerk.GetResponse{Key: rGet.GetKey(), Value: rGet.GetValue()}, errGet
}

func (s *server) Put(ctx context.Context, in *clerk.PutRequest) (*clerk.PutResponse, error){
	conn, err = grpc.Dial(*connaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewMapServiceClient(conn)

	// Contact the server and print out its response.
	myCtx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rPut, errPut := c.Put(myCtx, &pb.PutRequest{Key: in.GetKey(), Value: in.GetValue()})
	log.Printf("Passing (%v, %v) onto server as Put", in.GetKey(), in.GetValue())
	outCh <- []string{"put", in.GetKey(), in.GetValue()}
	return &clerk.PutResponse{Success: rPut.GetSuccess()}, errPut
}

func (s *server) Append(ctx context.Context, in *clerk.AppendRequest) (*clerk.AppendResponse, error){
	conn, err = grpc.Dial(*connaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c = pb.NewMapServiceClient(conn)

	// Contact the server and print out its response.
	myCtx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	rAppend, errAppend := c.Append(myCtx, &pb.AppendRequest{Key: in.GetKey(), Value: in.GetValue()})
	log.Printf("Passing (%v, %v) onto server as Append", in.GetKey(), in.GetValue())
	outCh <- []string{"append", in.GetKey(), in.GetValue()}
	return &clerk.AppendResponse{Oldval: rAppend.GetOldval()}, errAppend
}


// var conn, err = grpc.Dial(*connaddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
// var c = pb.NewMapServiceClient(conn)

// // Contact the server and print out its response.
// var ctx, cancel = context.WithTimeout(context.Background(), time.Second)

func main() {
	flag.Parse()
	go handleOrder()
	// Set up server and listen
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	clerk.RegisterMapServiceClerkServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
