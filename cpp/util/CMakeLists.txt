add_library(util)

target_sources(util
        PUBLIC
        FILE_SET CXX_MODULES FILES util.cppm
)

target_link_libraries(util PUBLIC fmt::fmt)
