# Testing Studio

A web-based tool for testing local cloud services - PubSub, Kafka, REST APIs, GCS, Spanner, and trace journey tracking.

## UI Screenshots

### Homepage
![Homepage](screenshots/ui-homepage.png)

### Settings / Config Editor
![Settings](screenshots/ui-settings.png)

### REST Client (Postman-style)
![REST Client](screenshots/ui-rest-client.png)

## Features

- **Google PubSub** - Pull and view CloudEvents from subscriptions
- **Kafka / EventMesh** - Consume and publish Avro messages
- **REST Client** - Send HTTP requests with collections (Postman-style), TLS certs, and JSON syntax highlighting
- **GCS Browser** - Browse buckets, preview files, and download
- **Spanner Explorer** - Query databases and browse tables
- **Trace Journey Viewer** - Track requests across containers with trace IDs

## Prerequisites

**Requires TMS Suncorp devstack to be running**

```bash
make devstack.start
```

## How to Run

```bash
# Build and run
go build -o testing-studio cmd/server/main.go
./testing-studio

# Or run directly
go run cmd/server/main.go
```

Open **http://localhost:8888** in your browser.

---

Made for TMS Suncorp local development
