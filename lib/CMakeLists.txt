add_library(vector_value)

target_sources(vector_value
        PUBLIC
        FILE_SET CXX_MODULES FILES vector_value.cppm
)

target_link_libraries(vector_value util)

add_library(helper_types)

target_sources(helper_types
        PUBLIC
        FILE_SET CXX_MODULES FILES helper_types.cppm
)

target_link_libraries(helper_types util vector_value)

add_library(function_traits)

target_sources(function_traits
        PUBLIC
        FILE_SET CXX_MODULES FILES function_traits.cppm
)

target_link_libraries(function_traits vector_value)

target_link_libraries(function_traits util)

add_library(compute)

target_sources(compute
    PUBLIC
        FILE_SET CXX_MODULES FILES Compute.cppm
)

target_link_libraries(compute util function_traits vector_value helper_types protolib)

