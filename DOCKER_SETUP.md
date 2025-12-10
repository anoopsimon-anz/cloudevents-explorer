# Docker Setup Summary

## What Changed

Successfully dockerized Testing Studio to work with TMS devstack services.

### Key Changes

1. **Dockerfile** - Multi-stage build following TMS datagen pattern:
   - Uses ANZ's internal Go proxy (`platform-gomodproxy.services.x.gcp.anz`)
   - Debian runtime (not Alpine) for glibc compatibility with CGO-enabled Kafka library
   - Non-root user (appuser:1000)

2. **docker-compose.yml** - Joins devstack network:
   - Connects to `devstack_devstack_network` (external network)
   - Accesses devstack services by Docker service names

3. **configs.json** - Dual configurations:
   - **Docker configs**: Use `dep_redpanda:9092` and `dep_pubsub:8086`
   - **Local configs**: Use `localhost:19092` and `localhost:8086`

## Network Architecture

```
┌─────────────────────────────────────────┐
│  devstack_devstack_network (bridge)     │
│                                         │
│  ┌──────────────┐   ┌────────────────┐ │
│  │ dep_redpanda │   │  dep_pubsub    │ │
│  │  :9092       │   │  :8086         │ │
│  └──────────────┘   └────────────────┘ │
│         ▲                   ▲          │
│         │                   │          │
│  ┌──────┴───────────────────┴───────┐  │
│  │   testing-studio                 │  │
│  │   (container)                    │  │
│  └──────────────────────────────────┘  │
│         │                               │
└─────────┼───────────────────────────────┘
          │
          │ Port mapping: 8888:8888
          │
          ▼
    localhost:8888 (host machine)
```

## Usage

### Running with Docker (Recommended)

```bash
# Start (builds if needed)
make start

# View logs
make logs

# Stop
make stop

# Rebuild
make rebuild
```

**Important**: In the UI, select **"Docker"** configurations:
- Unica Events (Docker) → `dep_redpanda:9092`
- TMS PubSub (Docker) → `dep_pubsub:8086`

### Running Locally (without Docker)

```bash
go run cmd/server/main.go
```

**Important**: In the UI, select **"Local"** configurations:
- Unica Events (Local) → `localhost:19092`
- TMS PubSub (Local) → `localhost:8086`

## Troubleshooting

### Problem: Docker container pulls zero Kafka/PubSub messages

**Cause**: Using "Local" configs which try to connect to `localhost` inside the container.

**Solution**: Switch to "Docker" configurations in the UI dropdown.

### Problem: Local run cannot connect to services

**Cause**: Using "Docker" configs which try to connect to `dep_redpanda` (Docker service name).

**Solution**: Switch to "Local" configurations in the UI dropdown.

### Verify Network Setup

```bash
# Check container is on devstack network
docker inspect testing-studio | grep devstack_devstack_network

# Verify DNS resolution inside container
docker exec testing-studio getent hosts dep_redpanda
docker exec testing-studio getent hosts dep_pubsub

# Expected output:
# 172.20.0.X      dep_redpanda
# 172.20.0.X      dep_pubsub

# Check all services are on the same network
docker network inspect devstack_devstack_network --format '{{range .Containers}}{{.Name}}: {{.IPv4Address}}{{"\n"}}{{end}}' | grep -E "(dep_redpanda|dep_pubsub|testing-studio)"

# Expected output:
# testing-studio: 172.20.0.X/16
# devstack-dep_pubsub-1: 172.20.0.X/16
# dep_redpanda: 172.20.0.X/16
```

### Verify Application is Running

```bash
# Test the API endpoint
curl http://localhost:8888/api/configs

# Should return JSON with both Docker and Local configurations
```

## Files Modified

- `Dockerfile` - Multi-stage build with Debian runtime
- `docker-compose.yml` - Joins devstack network
- `.dockerignore` - Optimizes build context
- `Makefile` - Docker commands (start/stop/build/logs)
- `configs.json` - Dual Docker/Local configurations
- `README.md` - Updated with Docker instructions

## Technical Details

### Why Debian instead of Alpine?

The confluent-kafka-go library requires CGO and links against glibc. Alpine uses musl libc, causing runtime errors:
```
exec /app/testing-studio: no such file or directory
Error relocating /app/testing-studio: getcontext: symbol not found
```

Solution: Use `debian:bookworm-slim` which has glibc and `librdkafka1`.

### Why external network?

The devstack services (`dep_redpanda`, `dep_pubsub`) exist on the `devstack_devstack_network`. To access them, Testing Studio must join the same network:

```yaml
networks:
  devstack_devstack_network:
    external: true  # Don't create, use existing
```
