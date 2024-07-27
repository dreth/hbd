# hbd

HBD is a simple application (which you can self-host) that serves birthday reminders through telegram.

## Usage

The application usage is very straightforward:

1. Sign up for an account
    + Requires an email, password, reminder time, timezone, [telegram bot API key](#bot-api-key), and [telegram chat ID](#chat-id)
2. Add birthdays
3. Receive reminders

### Getting a bot API key and chat ID

#### Bot API key

1. Open Telegram
2. Search for `BotFather`
3. Start a chat with `BotFather`
4. Use the `/newbot` command to create a new bot
5. Follow the instructions to create a new bot
6. Copy the API key

Through this API key the application can send messages to you through the bot.

#### Chat ID

1. Open Telegram on your mobile device
2. Send `/start` to your newly created bot
3. Send a message to the bot
4. Open the following URL in your browser: `https://api.telegram.org/bot<API_KEY>/getUpdates`
5. Look for the `chat` object in the JSON response
6. Copy the `id` field

Using this ID the application can send messages to your chat specifically.

## Self-hosting

### Docker

The application is containerized using docker, we use sqlite as the database by default.

#### Docker compose

The repo includes two docker-compose files. Feel free to customize them to your needs.

1. `docker-compose.yml` - This file is used to run the application locally with no reverse proxy. This works fine for use within a local network or for development. We map the ports `8417` and `8418` to the host machine. Port `8417` is the backend port and port `8418` is the frontend port. 

```yaml
---
---
services:
  hbd:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: hbd
    volumes:
      - ./data:/app/data
    ports:
      - "8417:8417"
      - "8418:8418"
    environment:
      - DB_TYPE=sqlite
      - DATABASE_URL=/app/data/hbd.db
      - MASTER_KEY=35e150e7ca83247f18cb1a37d61d8e161dddec06027f5db009b34da48c25f1b5
      - PORT=8418
      - ENVIRONMENT=development
      - CUSTOM_DOMAIN_FRONTEND=https://hbd.lotiguere.com
      - CUSTOM_DOMAIN_BACKEND=https://hbd-api.lotiguere.com
      - GIN_MODE=debug
```

2. `docker-compose.prod.yml` - This file is used to run the application in a production environment. We use traefik as a reverse proxy. We can optionally map ports to a local machine but it is not necessary, we opted not to.

```yaml
---
---
services:
  hbd:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: hbd
    volumes:
      - ./data:/app/data
    restart: always
    environment:
      - DB_TYPE=sqlite
      - DATABASE_URL=/app/data/hbd.db
      - MASTER_KEY=35e150e7ca83247f18cb1a37d61d8e161dddec06027f5db009b34da48c25f1b5
      - PORT=8418
      - ENVIRONMENT=production
      - CUSTOM_DOMAIN_FRONTEND=https://hbd.lotiguere.com
      - CUSTOM_DOMAIN_BACKEND=https://hbd-api.lotiguere.com
      - GIN_MODE=release
    labels:
      # Frontend
      - "traefik.enable=true"
      - "traefik.http.routers.hbd.rule=Host(`hbd.lotiguere.com`)"
      - "traefik.http.routers.hbd.entrypoints=websecure"
      - "traefik.http.routers.hbd.tls.certresolver=myresolver"
      - "traefik.http.routers.hbd.service=hbd-service"
      - "traefik.http.services.hbd-service.loadbalancer.server.port=8418"

      # Backend
      - "traefik.http.routers.hbd-api.rule=Host(`hbd-api.lotiguere.com`)"
      - "traefik.http.routers.hbd-api.entrypoints=websecure"
      - "traefik.http.routers.hbd-api.tls.certresolver=myresolver"
      - "traefik.http.routers.hbd-api.service=hbd-api-service"
      - "traefik.http.services.hbd-api-service.loadbalancer.server.port=8417"
    networks:
      - proxy

networks:
  proxy:
    name: proxy
    external: true
```

## Migrations

### SQLite

#### Up

```bash
migrate -database 'sqlite3://hbd.db' -path ./migrations up
```

#### Down

```bash
migrate -database 'sqlite3://hbd.db' -path ./migrations down
```

#### Down to last migration

```bash
migrate -database 'sqlite3://hbd.db' -path ./migrations down 1
```



## SQLBoiler

Generate models

```bash
sqlboiler psql
```

Point to the `.sqlboiler.toml` file

```bash
sqlboiler psql --config .sqlboiler.toml
```

## All together

```bash
migrate -database 'postgres://postgres:postgres@localhost:6684/postgres?sslmode=disable' -path ./migrations down && migrate -database 'postgres://postgres:postgres@localhost:6684/postgres?sslmode=disable' -path ./migrations up && sqlboiler psql --config .sqlboiler.toml && air
```

## Update swagger docs

```bash
swag init
```
