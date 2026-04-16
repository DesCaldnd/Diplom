cmake_minimum_required(VERSION 3.30)

project(Diplom LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_SCAN_FOR_MODULES ON)

# Настройка путей для LLVM/Clang на macOS
if(APPLE AND CMAKE_CXX_COMPILER_ID MATCHES "Clang")
    # Получаем путь к macOS SDK
    execute_process(
        COMMAND xcrun --show-sdk-path
        OUTPUT_VARIABLE MACOS_SDK_PATH
        OUTPUT_STRIP_TRAILING_WHITESPACE
    )
    
    # Устанавливаем флаги компилятора для поиска системных заголовков
    set(CMAKE_OSX_SYSROOT "${MACOS_SDK_PATH}")
endif()

find_package(Protobuf CONFIG REQUIRED)
find_package(gRPC CONFIG REQUIRED)
find_package(OpenSSL REQUIRED)
find_path(EXPRTK_INCLUDE_DIRS "exprtk.hpp")

add_library(protolib proto/grid.proto)
target_link_libraries(protolib PUBLIC gRPC::grpc++)
target_include_directories(protolib PUBLIC ${CMAKE_CURRENT_BINARY_DIR})

get_target_property(grpc_cpp_plugin_location gRPC::grpc_cpp_plugin LOCATION)

protobuf_generate(TARGET protolib LANGUAGE cpp)
protobuf_generate(TARGET protolib LANGUAGE grpc
        GENERATE_EXTENSIONS .grpc.pb.h .grpc.pb.cc
        PLUGIN "protoc-gen-grpc=${grpc_cpp_plugin_location}")

add_subdirectory(util)
add_subdirectory(lib)
add_subdirectory(server)
add_subdirectory(examples)
add_subdirectory(test)
