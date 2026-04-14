# CMakeLists.txt для модуля server

# Создание библиотеки server
add_library(server)

# Добавление исходных файлов к библиотеке
target_sources(server
    PUBLIC
        FILE_SET CXX_MODULES FILES server.cppm
)

# Линковка необходимых библиотек
# Здесь нужно будет добавить библиотеки gRPC после их установки через vcpkg
# target_link_libraries(server <grpc_libraries>)

# Создание исполняемого файла main_api
add_executable(main_api main_api.cpp)

# Линковка библиотек к исполняемому файлу
target_link_libraries(main_api PRIVATE server)

# После установки gRPC через vcpkg потребуется добавить соответствующие зависимости
# target_link_libraries(main_api PRIVATE <grpc_libraries>)