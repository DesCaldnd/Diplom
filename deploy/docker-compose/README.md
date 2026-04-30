# Docker Compose deployment

Конфигурация в [`deploy/docker-compose/docker-compose.yml`](deploy/docker-compose/docker-compose.yml) поднимает:
- [`nginx`](deploy/nginx/nginx.conf) как внешний reverse proxy с TLS termination
- несколько реплик `server` внутри docker-сети
- [`envoy`](deploy/envoy.yaml) для grpc-web
- `web`, `prometheus`, `grafana`

## Маршрутизация

Внешние запросы идут через [`nginx`](deploy/nginx/nginx.conf):
- `/` -> редирект на `/web/`
- `/web/` -> `web:80`
- `/api/...` -> gRPC на `server:9999`
- `/envoy/` -> `envoy:8082`
- `/envoy-adm/` -> `envoy:9901`
- `/grafana/` -> `grafana:3000`

Публично открыты только порты `80` и `443` у контейнера `nginx`. Остальные сервисы доступны только внутри сети [`app_network`](deploy/docker-compose/docker-compose.yml).

## HTTP локально / HTTPS на сервере

Используется один [`nginx`](deploy/nginx/nginx.conf), который всегда публикует `80` и `443`.

Логика простая:
- основной конфиг [`deploy/nginx/nginx.conf`](deploy/nginx/nginx.conf) всегда поднимает HTTP-сервер на `80`
- при старте контейнера команда из [`deploy/docker-compose/docker-compose.yml`](deploy/docker-compose/docker-compose.yml:10) проверяет наличие сертификатов
- если сертификаты есть, подключается дополнительный include [`deploy/nginx/ssl-server.conf.template`](deploy/nginx/ssl-server.conf.template), который включает `443 ssl` и редирект с HTTP на HTTPS
- если сертификатов нет, HTTPS-блок не подключается, и система работает локально только по HTTP

То есть:
- локально без сертификатов автоматически работает `http://localhost/web/`
- на сервере при наличии сертификатов автоматически работает HTTPS

## Локальный запуск

```bash
cd deploy/docker-compose
docker compose up --build
```

Приложение будет доступно по `http://localhost/web/`.

## Запуск на сервере с HTTPS

Перед запуском положи сертификаты в каталог [`deploy/docker-compose/certs`](deploy/docker-compose/certs):
- `fullchain.pem`
- `privkey.pem`

Затем:

```bash
cd deploy/docker-compose
docker compose up --build -d
```

## Получение сертификата для diplom.funandchecks.ru

На сервере нужно:

1. Настроить DNS:
   - `A` запись `diplom.funandchecks.ru` на IP сервера.

2. Открыть порты `80` и `443` во внешнем firewall.

3. Установить `certbot`.

Для Ubuntu/Debian:

```bash
sudo apt update
sudo apt install -y certbot
```

4. Остановить всё, что слушает `80` порт, либо временно остановить [`docker compose`](deploy/docker-compose/docker-compose.yml).

5. Выпустить сертификат standalone-командой:

```bash
sudo certbot certonly --standalone -d diplom.funandchecks.ru
```

6. Скопировать сертификаты в каталог проекта:

```bash
sudo mkdir -p /home/descaldnd/Diplom/deploy/docker-compose/certs
sudo cp /etc/letsencrypt/live/diplom.funandchecks.ru/fullchain.pem /home/descaldnd/Diplom/deploy/docker-compose/certs/fullchain.pem
sudo cp /etc/letsencrypt/live/diplom.funandchecks.ru/privkey.pem /home/descaldnd/Diplom/deploy/docker-compose/certs/privkey.pem
sudo chown -R $USER:$USER /home/descaldnd/Diplom/deploy/docker-compose/certs
```

7. Запустить compose:

```bash
cd /home/descaldnd/Diplom/deploy/docker-compose
docker compose up --build -d
```

## Автообновление сертификатов

Пример cron:

```bash
0 4 * * * certbot renew --quiet && cp /etc/letsencrypt/live/diplom.funandchecks.ru/fullchain.pem /home/descaldnd/Diplom/deploy/docker-compose/certs/fullchain.pem && cp /etc/letsencrypt/live/diplom.funandchecks.ru/privkey.pem /home/descaldnd/Diplom/deploy/docker-compose/certs/privkey.pem && docker compose -f /home/descaldnd/Diplom/deploy/docker-compose/docker-compose.yml exec nginx nginx -s reload
```

## Несколько реплик `server`

Если используется обычный [`docker compose`](deploy/docker-compose/docker-compose.yml), масштабирование делается так:

```bash
docker compose up --build --scale server=3
```

Внутри docker-сети имя `server` будет резолвиться в несколько IP, а [`nginx`](deploy/nginx/nginx.conf) распределяет запросы через upstream [`api_upstream`](deploy/nginx/nginx.conf:35).
