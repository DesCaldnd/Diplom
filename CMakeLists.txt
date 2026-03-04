cmake_minimum_required(VERSION 3.30)

project(Diplom LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 23)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_SCAN_FOR_MODULES ON)

add_subdirectory(util)
add_subdirectory(lib)
add_subdirectory(examples)
add_subdirectory(test)
