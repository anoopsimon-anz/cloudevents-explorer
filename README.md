# 📡 Testing Studio

A professional web tool for testing and debugging cloud services locally including Google Cloud PubSub, Kafka/EventMesh, REST APIs, Google Cloud Storage, and Google Cloud Spanner.

![Version](https://img.shields.io/badge/version-3.0.0-blue)
![Go](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)

## ✨ Features

### 📦 GCS Browser
- Browse Google Cloud Storage buckets and objects
- Preview file contents directly in browser
- Download files with one click
- Folder navigation with breadcrumbs
- Works with fake-gcs-server for local development

### 📬 Google PubSub Explorer
- Pull messages from PubSub subscriptions
- View CloudEvents with syntax-highlighted JSON
- Collapsible message cards for easy scanning
- Works with PubSub emulator (bypasses proxy issues on macOS)

### 🌊 Kafka / EventMesh Browser
- Consume Avro messages from Kafka topics
- Automatic schema registry integration
- Publish test events to topics
- Message filtering and search

### 🔌 REST Client
- Send HTTP requests with custom headers
- Support for TLS client certificates
- Request/response viewer
- Perfect for testing internal APIs

### 📊 Flow Diagram Tool
- Visualize message flows
- Base64 encoder/decoder

### 🗄️ Spanner Explorer
- Connect to local Spanner emulator
- Browse tables and view schemas
- Execute SQL queries (SELECT and DML)
- Save multiple connection profiles
- Spanner Studio-style interface

### 🔍 Status Indicators
- Docker status monitoring
- GCloud authentication status with last login time
- Real-time updates

## 🚀 Quick Start

```bash
cd ~/scratches/cloudevents-explorer

# Run directly
go run cmd/server/main.go

# Or build and run
go build -o testing-studio cmd/server/main.go
./testing-studio
```

Open **http://localhost:8888** in your browser.

## 🏗️ Architecture

Clean, modular architecture for maintainability:

```
testing-studio/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── handlers/               # HTTP request handlers
│   │   ├── gcs.go             # GCS Browser handlers
│   │   ├── pubsub.go          # PubSub handlers
│   │   ├── kafka.go           # Kafka handlers
│   │   ├── spanner.go         # Spanner handlers
│   │   ├── rest_client.go     # REST client handlers
│   │   ├── docker_status.go   # Docker status checker
│   │   └── gcloud_status.go   # GCloud auth checker
│   ├── kafka/                  # Kafka operations
│   ├── pubsub/                 # PubSub operations
│   ├── spanner/                # Spanner operations
│   ├── templates/              # HTML templates
│   └── types/                  # Shared types
├── go.mod
└── README.md
```

## 📖 Usage Guide

### GCS Browser
1. Navigate to **http://localhost:8888/gcs**
2. Select a bucket from the sidebar
3. Browse folders and files
4. Click **Preview** to view file contents
5. Click **Download** to save files locally

**Local Development**: Configure fake-gcs-server on `localhost:4443`

### PubSub Explorer
1. Navigate to **http://localhost:8888/pubsub**
2. Configure connection settings:
   - **Emulator Host**: `localhost:8086`
   - **Project ID**: `tms-suncorp-local`
   - **Subscription ID**: `cloudevents.subscription`
3. Click **Pull Messages**
4. Expand messages to view details

### Kafka Browser
1. Navigate to **http://localhost:8888/kafka**
2. Configure Kafka settings:
   - **Bootstrap Servers**: `localhost:19092`
   - **Schema Registry**: `http://localhost:8081`
   - **Topic**: Select from dropdown
3. Click **Pull Messages** to consume
4. Use **Publish** tab to send test events

### REST Client
1. Navigate to **http://localhost:8888/rest-client**
2. Enter request details:
   - Method (GET, POST, PUT, DELETE)
   - URL
   - Headers (one per line, format: `Key: Value`)
   - Request body (for POST/PUT)
3. Optionally configure TLS client certificate
4. Click **Send Request**

### Spanner Explorer
1. Navigate to **http://localhost:8888/spanner**
2. Configure connection settings:
   - **Emulator Host**: `localhost:9010`
   - **Project ID**: `tms-suncorp-local`
   - **Instance ID**: `tms-suncorp-local`
   - **Database ID**: `tms-suncorp-db`
3. Click **Connect** to test connection
4. Click **Load Tables** to browse database tables
5. Click a table to select it
6. Enter SQL query in editor or use example queries
7. Click **Run Query** to execute

**Local Development**: Requires Spanner emulator running on `localhost:9010`

## 🎯 Why This Tool?

### The Problem
Working with local cloud service emulators on macOS with corporate proxies causes issues:
- HTTP/2 errors with PubSub emulator
- Kafka connection failures
- TLS certificate complexity for internal APIs
- No easy way to browse GCS buckets locally

### The Solution
Testing Studio provides native Go SDK integration for all services, bypassing proxy/HTTP issues:
- **PubSub**: Uses native gRPC SDK
- **Kafka**: Direct Kafka protocol connection
- **GCS**: Direct HTTP API access to fake-gcs-server
- **REST**: Built-in HTTP client with TLS support

## 🔧 Local Development Setup

### Prerequisites
- Go 1.21 or higher
- (Optional) PubSub emulator running on `localhost:8086`
- (Optional) Kafka/Redpanda running on `localhost:19092`
- (Optional) fake-gcs-server running on `localhost:4443`
- (Optional) Spanner emulator running on `localhost:9010`

### Running Services

**PubSub Emulator**:
```bash
gcloud beta emulators pubsub start --project=tms-suncorp-local --host-port=localhost:8086
```

**Kafka (via Redpanda)**:
```bash
# See your TMS devstack setup
make devstack.start
```

**fake-gcs-server**:
```bash
docker run -d --name gcs \
  -p 4443:4443 \
  fsouza/fake-gcs-server -scheme http -port 4443
```

**Spanner Emulator**:
```bash
docker run -d --name spanner-emulator \
  -p 9010:9010 \
  gcr.io/cloud-spanner-emulator/emulator
```

## 📂 Configuration

Configurations are stored in `configs.json`:

```json
{
  "pubsub": [
    {
      "name": "TMS Local",
      "emulatorHost": "localhost:8086",
      "projectId": "tms-suncorp-local",
      "subscriptionId": "cloudevents.subscription"
    }
  ],
  "kafka": [
    {
      "name": "Local Kafka",
      "bootstrapServers": "localhost:19092",
      "schemaRegistry": "http://localhost:8081",
      "topic": "cloudevents"
    }
  ],
  "spanner": [
    {
      "name": "TMS Local",
      "emulatorHost": "localhost:9010",
      "projectId": "tms-suncorp-local",
      "instanceId": "tms-suncorp-local",
      "databaseId": "tms-suncorp-db"
    }
  ]
}
```

## 🎨 Features in Detail

### GCS Browser Features
- **Bucket listing**: See all available buckets
- **Folder navigation**: Navigate through object prefixes
- **File preview**: View text files, JSON, and other formats
- **Download**: One-click file downloads
- **Breadcrumb navigation**: Easy navigation back to parent folders

### Status Indicators
- **Docker**: Shows green when Docker daemon is running
- **GCloud**: Shows authentication status and last login time
- Auto-refreshes every 10-30 seconds

### Syntax Highlighting
- JSON syntax highlighting for PubSub and Kafka messages
- Color-coded keys, strings, numbers, and booleans
- Collapsible/expandable message cards

## 🐛 Troubleshooting

### Port 8888 already in use
```bash
lsof -ti:8888 | xargs kill -9
```

### Can't connect to PubSub emulator
- Verify emulator is running: `curl localhost:8086`
- Check project ID matches configuration
- Ensure subscription exists

### Can't connect to Kafka
- Verify Redpanda/Kafka is running: `nc -zv localhost 19092`
- Check bootstrap servers in configuration
- Verify topic exists

### GCS Browser shows no buckets
- Verify fake-gcs-server is running on `localhost:4443`
- Create a test bucket using the GCS API
- Check server logs for errors

### Can't connect to Spanner emulator
- Verify emulator is running: `docker ps | grep spanner`
- Check emulator host in configuration
- Ensure database exists in emulator
- Check emulator logs: `docker logs spanner-emulator`

### GCloud status shows "Not authenticated"
- Run `gcloud auth login`
- Verify token is not expired
- Check `~/.config/gcloud/` directory exists

## 📝 Development

### Hot Reload
Use `air` for development:
```bash
go install github.com/cosmtrek/air@latest
air
```

### Dependencies
```bash
go mod download
```

### Building
```bash
# Development build
go build -o testing-studio cmd/server/main.go

# Production build
go build -ldflags="-s -w" -o testing-studio cmd/server/main.go
```

## 🎨 UI Design

- Clean, professional Google Cloud-inspired interface
- Material Design color palette
- Responsive layout
- Smooth animations and transitions
- Accessibility-friendly

## 📜 Version History

- **v3.0.0** - Added GCS Browser, removed Docker configuration
- **v2.0.0** - Added Kafka support, Flow Diagram, REST Client
- **v1.0.0** - Initial PubSub Explorer

## 🙏 Acknowledgments

- Built for TMS Suncorp local development
- UI inspired by Google Cloud Console
- Solves real proxy/HTTP pain points on macOS

---

**Made with ❤️ for local cloud development without headaches**
