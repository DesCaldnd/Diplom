add_library(server)

target_sources(server
    PUBLIC
        FILE_SET CXX_MODULES FILES server.cppm
)

target_link_libraries(server PUBLIC protolib compute)

add_executable(main_api main_api.cpp)

target_link_libraries(main_api PRIVATE server)