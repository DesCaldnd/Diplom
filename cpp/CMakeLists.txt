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

find_package(fmt CONFIG REQUIRED)
find_package(GTest CONFIG REQUIRED)

add_subdirectory(util)
add_subdirectory(lib)
add_subdirectory(examples)
add_subdirectory(test)
