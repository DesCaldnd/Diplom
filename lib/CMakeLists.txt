add_library(function_traits)

target_sources(function_traits
        PUBLIC
        FILE_SET CXX_MODULES FILES function_traits.cppm
)

target_link_libraries(function_traits util)

add_library(compute)

target_sources(compute
    PUBLIC
        FILE_SET CXX_MODULES FILES compute.cppm
)

target_link_libraries(compute util function_traits)

