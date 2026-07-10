# sahil-api

Live "now" data aggregator for [sahil.im](https://sahil.im). Fetches the latest activity from personal media services and exposes them over a unified JSON API.

## Endpoints

| Method | Path                          | Description                          |
| ------ | ----------------------------- | ------------------------------------ |
| GET    | `/health`                     | Health check                         |
| GET    | `/api/now`                    | Aggregated song, book, anime, movie   |
| GET    | `/api/lastfm`                 | Last played track                    |
| GET    | `/api/anilist`                | Last updated anime                   |
| GET    | `/api/letterboxd`             | Last movie diary entry               |
| GET    | `/api/goodreads`              | Currently reading book                |
| GET    | `/api/cms/feed`               | Get feed entry (public)              |
| POST   | `/api/cms/feed`               | Upsert feed entry (Bearer auth)      |
| GET    | `/api/greetings`              | Greetings in 99 languages            |
| GET    | `/api/statuses`               | Random status sentences              |
| GET    | `/api/openapi.yaml`           | OpenAPI 3.0 specification            |

## Quick Start

```bash
cp .env.example .env
# edit .env with your credentials
make run
```

The server starts on port `8080` (or `$PORT`).

## Configuration

All configuration is via environment variables:

| Variable            | Required | Description                     |
| ------------------- | -------- | ------------------------------- |
| `PORT`              |          | Server port (default: 8080)     |
| `API_BASE_URL`      |          | Public base URL for self-refs   |
| `LASTFM_API_KEY`     | for lastfm    | Last.fm API key                 |
| `LASTFM_USERNAME`    | for lastfm    | Last.fm username                |
| `ANILIST_USERNAME`   | for anilist   | AniList username                |
| `LETTERBOXD_USERNAME`| for letterboxd| Letterboxd username             |
| `GOODREADS_USER_ID`  | for goodreads | Goodreads numeric user ID       |
| `CMS_API_TOKEN`      | for cms       | Bearer token for Obsidian plugin |
| `CLOUDFLARE_ACCOUNT_ID`| for cms     | Cloudflare account ID           |
| `CLOUDFLARE_API_TOKEN` | for cms     | Cloudflare API token            |
| `D1_DATABASE_ID`     | for cms       | D1 database ID                  |
| `R2_ENDPOINT`        | for cms       | R2 S3 endpoint                  |
| `R2_ACCESS_KEY_ID`   | for cms       | R2 S3 access key                |
| `R2_SECRET_ACCESS_KEY`| for cms      | R2 S3 secret key                |
| `R2_BUCKET`          | for cms       | R2 bucket name                  |

## Architecture

```
cmd/server/main.go      — entrypoint, server setup, graceful shutdown
internal/
  config/config.go       — typed configuration from env vars
  model/types.go         — shared data types
  cache/cache.go         — in-memory TTL cache
  fetcher/               — upstream API clients (one per service)
  handler/               — HTTP handlers (thin, delegates to fetcher + cache)
  middleware/             — request ID, structured logging, recovery, CORS
openapi.yaml             — OpenAPI 3.0 spec
```

Each service is fetched independently and cached for 5 minutes. The `/api/now` endpoint fans out to all services concurrently.

## Deployment

### Render

1. Create a new Web Service
2. Build command: `cd api && go build -o server ./cmd/server`
3. Start command: `./server`
4. Set environment variables from `.env.example`

### Docker

```dockerfile
FROM golang:1.26-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:3.20
COPY --from=build /app/server /server
EXPOSE 8080
CMD ["/server"]
```

## Development

```bash
make fmt      # format code
make vet      # static analysis
make build    # compile binary
make run      # run locally
```
