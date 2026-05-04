include(GoogleTest)

add_executable(tests main.cpp)
target_link_libraries(tests PRIVATE compute util GTest::gtest GTest::gtest_main)
gtest_discover_tests(tests)

add_executable(bench bench.cpp)
target_link_libraries(bench PRIVATE compute util GTest::gtest GTest::gtest_main)
gtest_discover_tests(bench)