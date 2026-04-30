package main

//go:generate protoc --proto_path=../../proto --go_out=../../gridProto --go_opt=paths=source_relative --go-grpc_out=../../gridProto --go-grpc_opt=paths=source_relative grid.proto

import (
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	pb "Diplom/gridProto"
	"Diplom/internal/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting server...")

	// Запуск HTTP сервера для pprof и Prometheus метрик
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		fmt.Println("Starting HTTP server for metrics and pprof on :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("failed to serve http: %v", err)
		}
	}()

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
