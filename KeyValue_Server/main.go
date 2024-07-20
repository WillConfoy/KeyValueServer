package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "cs498.com/KeyValueServer/KeyValue"
)

type server struct {
	pb.UnimplementedMapServiceServer
}

var theMap = make(map[string]string)

var (
	port = flag.Int("port", 50052, "The server port")
)

// Note that keys and values are strings
func (s *server) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {
	log.Printf("Received %v as key, %v as value for PutRequest", in.GetKey(), in.GetValue())
	theMap[in.GetKey()] = in.GetValue()
	return &pb.PutResponse{Success: true}, nil
}


func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	log.Printf("Received %v as key for GetRequest", in.GetKey())
	return &pb.GetResponse{Key: in.GetKey(), Value: theMap[in.GetKey()]}, nil
}


func (s *server) Append(ctx context.Context, in *pb.AppendRequest) (*pb.AppendResponse, error) {
	log.Printf("Received %v as key, %v as value for AppendRequest", in.GetKey(), in.GetValue())
	oldVal := theMap[in.GetKey()]
	theMap[in.GetKey()] += in.GetValue()
	return &pb.AppendResponse{Oldval: oldVal}, nil
}


func main() {
	// outCh := make(chan [][]string)
	// go doPut(theMap, outCh)
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMapServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}