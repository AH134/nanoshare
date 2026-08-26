#

<h1 align="center">
    <img src="assets/logo.png" width="200" height="200" alt="nanoshare">
    <br/>
    <span>Nanoshare<span>
</h1>

<p align="center">File sharing web app built with <a href="https://go.dev">Go</a> and <a href="https://react.dev">React</a></p>

Drag and drop your files, set a max download count and expiry date and get a shareable link for sharing.

## Getting Started

### Docker

Start the container with `docker run`:

```bash
docker run \
    --name nanoshare \
    -p 8080:8080 \
    -v nanoshare-data:/app/data \
    -e ADMIN_USERNAME=admin \
    -e ADMIN_PASSWORD=admin \
    -d \
    ghcr.io/xw134/nanoshare:latest
```

The app will be available at http://localhost:8080.

`ADMIN_USERNAME` and `ADMIN_PASSWORD` are required. The user will be seeded upon container creation.

### Docker Compose

```yaml
services:
  nanoshare:
    image: xw134/nanoshare:latest
    ports:
      - "8080:8080"
    env_file:
      - .env
    environment:
      ADMIN_USERNAME: admin
      ADMIN_PASSWORD: admin
    volumes:
      - data:/app/data
    healthcheck:
      test:
        ["CMD", "wget", "--spider", "-q", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 3s
      retries: 5
      start_period: 5s

    restart: unless-stopped

networks:
  nanoshare_net:
    driver: bridge

volumes:
  data:
```

Then run:

```bash
docker compose up -d
```

### Using an env file

Instead of passing variables, you can use a `.env` file

```bash
# .env

ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin
```

```bash
# docker run
docker run \
...
--env-file .env \
...

# docker compose
services:
  nanoshare:
    image: xw134/nanoshare:latest
    ...
    env_file:
      - .env
    ...
```

### Envrionment Variables

| Variable         | Description                                 | Default               |
| ---------------- | ------------------------------------------- | --------------------- |
| `ADMIN_USERNAME` | Username for admin login                    | (required)            |
| `ADMIN_PASSWORD` | Password for admin login                    | (required)            |
| `PORT`           | Port the server listens on                  | `8080`                |
| `DB_PATH`        | Path to SQLite database                     | `./data/nanoshare.db` |
| `STORAGE_PATH`   | Path were the uploaded files are stored     | `./data/storage/`     |
| `PROD`           | Makes session cookies `Secure` (HTTPS only) | `false`               |

> [!Warning]
> **About `PROD`:** this controls whether session cookies are marked `Secure`. Leave it `false` for local, LAN, or plain-HTTP access. Set it to `true` once nanoshare is served over HTTPS (eg. behind a reverse-proxy such as Nginx, Caddy, Traefik).

### Building from Source

You can build the application from source and run it directly using the provided `Makefile`

#### Quickstart

```bash
git clone https://github.com/AH134/nanoshare.git
cd nanoshare

cp .env.example .env

make build

./bin/nanoshare
```

This project is under the [MIT License](LICENSE)
