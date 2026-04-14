# Dockerfile для сборки контейнера с gRPC сервером

# Используем официальный образ Ubuntu как базовый
FROM ubuntu:22.04

# Установка необходимых зависимостей
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    && rm -rf /var/lib/apt/lists/*

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Установка vcpkg и сборка только необходимых библиотек
RUN git clone https://github.com/Microsoft/vcpkg.git /opt/vcpkg && \
    cd /opt/vcpkg && \
    ./bootstrap-vcpkg.sh && \
    # Установка только необходимых библиотек (grpc в данном случае)
    ./vcpkg install grpc

# Установка переменной окружения для vcpkg
ENV VCPKG_ROOT=/opt/vcpkg

# Сборка проекта
RUN mkdir build && cd build && \
    cmake .. -DCMAKE_TOOLCHAIN_FILE=${VCPKG_ROOT}/scripts/buildsystems/vcpkg.cmake && \
    cmake --build . --config Release

# Открываем порт для gRPC сервера (предполагаем, что используется порт 50051)
EXPOSE 50051

# Команда для запуска сервера
CMD ["./build/server/main_api"]