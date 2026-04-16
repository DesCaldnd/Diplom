#include <iostream>
#include <grpcpp/grpcpp.h>
import server;

int main() {
    std::cout << "Starting gRPC server..." << std::endl;

    GridServer service;
    
    // TODO: Добавьте здесь код для запуска gRPC сервера
    grpc::ServerBuilder builder;
    builder.AddListeningPort("0.0.0.0:9999", grpc::InsecureServerCredentials());
    builder.RegisterService(&service);
    std::unique_ptr<grpc::Server> server(builder.BuildAndStart());
    server->Wait();
    
    return 0;
}