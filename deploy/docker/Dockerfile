# ==============================================================================
# СТАДИЯ 1: Сборка (Builder)
# ==============================================================================
FROM ubuntu:22.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive

# 1. Установка базовых утилит
RUN apt-get update && apt-get install -y \
    wget curl git zip unzip tar pkg-config \
    software-properties-common lsb-release gnupg \
    && rm -rf /var/lib/apt/lists/*

# 2. Установка CMake 3.28+
RUN wget -O - https://apt.kitware.com/keys/kitware-archive-latest.asc 2>/dev/null | gpg --dearmor - | tee /usr/share/keyrings/kitware-archive-keyring.gpg >/dev/null && \
    echo 'deb [signed-by=/usr/share/keyrings/kitware-archive-keyring.gpg] https://apt.kitware.com/ubuntu/ jammy main' | tee /etc/apt/sources.list.d/kitware.list >/dev/null && \
    apt-get update && apt-get install -y cmake && \
    rm -rf /var/lib/apt/lists/*

# 3. Установка Ninja 1.11.1+
RUN wget https://github.com/ninja-build/ninja/releases/latest/download/ninja-linux.zip && \
    unzip ninja-linux.zip -d /usr/local/bin/ && \
    rm ninja-linux.zip

# 4. Установка Clang 19 и привязка системных утилит
RUN wget https://apt.llvm.org/llvm.sh && \
    chmod +x llvm.sh && \
    ./llvm.sh 19 all && \
    rm llvm.sh

RUN update-alternatives --install /usr/bin/clang clang /usr/bin/clang-19 100 && \
    update-alternatives --install /usr/bin/clang++ clang++ /usr/bin/clang++-19 100 && \
    update-alternatives --install /usr/bin/c++ c++ /usr/bin/clang++-19 100 && \
    update-alternatives --install /usr/bin/cc cc /usr/bin/clang-19 100

ENV CC=clang-19
ENV CXX=clang++-19

WORKDIR /app

# 5. Установка vcpkg
ENV VCPKG_ROOT=/opt/vcpkg
RUN git clone https://github.com/Microsoft/vcpkg.git ${VCPKG_ROOT} && \
    ${VCPKG_ROOT}/bootstrap-vcpkg.sh -disableMetrics

# ==============================================================================
# МАГИЯ КЭШИРОВАНИЯ ЗАВИСИМОСТЕЙ
# ==============================================================================
# Создаем директорию для бинарного кэша vcpkg и указываем ее через ENV
ENV VCPKG_DEFAULT_BINARY_CACHE=/root/.cache/vcpkg
RUN mkdir -p /root/.cache/vcpkg

# Копируем ТОЛЬКО файл зависимостей.
# Если вы не меняли vcpkg.json, Docker возьмет следующий шаг из своего кэша!
COPY vcpkg.json ./
# Если у вас есть файл фиксации версий (vcpkg-configuration.json), раскомментируйте:
# COPY vcpkg-configuration.json ./

# Выполняем предварительную сборку библиотек.
# Vcpkg скачает исходники, скомпилирует их Clang'ом и положит в бинарный кэш.
RUN ${VCPKG_ROOT}/vcpkg install

# ==============================================================================
# СБОРКА ПРОЕКТА
# ==============================================================================
# Только теперь копируем исходный код нашего проекта (C++ файлы, модули, CMakeLists)
COPY . .

# Конфигурация и сборка проекта
# На этапе конфигурации CMake автоматически запросит нужные либы у vcpkg,
# но vcpkg просто достанет их из /root/.cache/vcpkg за 1-2 секунды!
RUN cmake -B build -S . \
    -GNinja \
    -DCMAKE_TOOLCHAIN_FILE=${VCPKG_ROOT}/scripts/buildsystems/vcpkg.cmake \
    -DCMAKE_C_COMPILER=clang-19 \
    -DCMAKE_CXX_COMPILER=clang++-19 \
    -DCMAKE_BUILD_TYPE=Release && \
    cmake --build build

# ==============================================================================
# СТАДИЯ 2: Релиз (Runtime) - легкий образ
# ==============================================================================
FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Устанавливаем корневые сертификаты на случай, если gRPC будет делать внешние вызовы
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем только готовый бинарник
COPY --from=builder /app/build/server/main_api ./main_api

EXPOSE 9999

CMD ["./main_api"]
