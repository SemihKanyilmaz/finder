# Finder Search API

Fetches content from two providers (JSON + XML), scores each item, saves to PostgreSQL, and returns ranked search results.

---

## Stack

- **Go** — fast, low memory usage
- **Echo** — simple web framework
- **PostgreSQL** — stores content, full-text search with `tsvector`
- **Redis** — caching for provider responses and search results; also **rate limiting**

---

## How it works

1. Request hits `GET /search`
2. Service fetches from both providers in concurrent
3. Each item gets a score (based on views, likes, reactions, freshness)
4. Items are saved to PostgreSQL (upsert)
5. DB returns filtered + sorted + paginated results

---

## Scoring

```
score = (base * type_coefficient) + freshness + engagement

base:
  video:   views/1000 + likes/100
  article: reading_time + reactions/50

type_coefficient:
  video:   1.5
  article: 1.0

freshness:
  ≤ 1 week:   +5
  ≤ 1 month:  +3
  ≤ 3 months: +1
  older:       0

engagement:
  video:   (likes / views) * 10
  article: (reactions / reading_time) * 5
```

---

## Setup

### With Docker (recommended)

```bash
docker compose up --build
```

App starts at `http://localhost:8080`. Migration runs automatically on first boot.

### Without Docker

**Requirements:** Go 1.26+, PostgreSQL, Redis

**1. Set environment variables** (copy `.env` and edit):

```env
SERVER_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/finder?sslmode=disable
REDIS_URL=redis://localhost:6379/0
PROVIDER1_NAME=provider-1
PROVIDER1_BASE_URL=https://raw.githubusercontent.com
PROVIDER2_NAME=provider-2
PROVIDER2_BASE_URL=https://raw.githubusercontent.com
RATE_LIMIT_PER_SEC=10
CACHE_TTL=5m
PROVIDER_CACHE_TTL=1m
CB_TIMEOUT=30s
CB_THRESHOLD=3
```

**2. Run migration:**

```bash
psql $DATABASE_URL -f migrations/001_create_contents.up.sql
```

**3. Start:**

```bash
go run main.go
```

---

## API

Full docs at `http://localhost:8080/swagger/index.html`

### GET /search

| Param | Description | Default |
|---|---|---|
| `q` | keyword | — |
| `type` | `video` or `article` | — |
| `sortBy` | `score` or `freshness` | `score` |
| `page` | page number | `1` |
| `pageSize` | results per page | `10` |

### GET /health

Returns `200 OK`.

---

## Tests

```bash
go test ./...
```
