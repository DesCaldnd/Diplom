#include <iostream>
#include <grpcpp/grpcpp.h>
import server;

int main() {
    std::cout << "Starting gRPC server..." << std::endl;

    GridServer service;

    grpc::ServerBuilder builder;
    builder.AddListeningPort("0.0.0.0:9999", grpc::InsecureServerCredentials());
    builder.RegisterService(&service);
    std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
    if (!server) {
        std::cerr << "Failed to start gRPC server" << std::endl;
        return 1;
    }
    server->Wait();
    
    return 0;
}