package main

//go:generate protoc --proto_path=../../proto --go_out=../../gridProto --go_opt=paths=source_relative --go-grpc_out=../../gridProto --go-grpc_opt=paths=source_relative grid.proto

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"Diplom/internal/server"
	pb "Diplom/gridProto"
)

func main() {
	fmt.Println("Starting server...")
	lis, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGridServiceServer(s, &server.GridServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
