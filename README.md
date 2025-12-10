# 📡 Testing Studio (CloudEvents Explorer)

A professional, modular web tool for exploring and testing Kafka and Google Cloud PubSub messages with Avro schema support, message publishing, and comprehensive flow diagrams.

![Version](https://img.shields.io/badge/version-2.0.0-blue)
![Go](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)
![Architecture](https://img.shields.io/badge/architecture-modular-green)

## ✨ Features

### 🎨 Redpanda-Inspired Dark Mode UI
- Beautiful GitHub dark theme color scheme
- Clean, professional interface
- Responsive design that works on any screen

### 📬 Message Management
- **Pull messages** from PubSub subscriptions
- **Persistent storage** - messages stay in memory across pulls
- **Collapsible cards** - each message starts collapsed for easy scanning
- **Expand to view** - click any message to see full details

### 🌈 Syntax-Highlighted JSON
- Color-coded JSON viewer
- Keys, strings, numbers, booleans, and null values all highlighted
- Pretty-printed with proper indentation
- Easy to read and understand complex payloads

### 💾 Configuration Management
- Save multiple PubSub configurations
- Quick dropdown switching between configs
- Stored in JSON for easy editing
- Pre-configured for TMS local development

### 📊 Real-Time Stats
- Total message count
- Last updated timestamp
- Live updates as you pull messages

## 🚀 Quick Start

### Docker (Recommended - connects to devstack)

```bash
cd ~/scratches/cloudevents-explorer

# Start the application (builds if needed)
make start

# View logs
make logs

# Stop the application
make stop
```

Open http://localhost:8888 in your browser.

**Note**: When running in Docker, use the "Docker" configurations which connect to `dep_redpanda` and `dep_pubsub` services in the devstack network.

### Local Development (without Docker)

```bash
cd ~/scratches/cloudevents-explorer

# Option 1: Use the quick start script
./start.sh

# Option 2: Run directly
go run cmd/server/main.go

# Option 3: Build and run
go build -o testing-studio cmd/server/main.go
./testing-studio
```

**Note**: When running locally, use the "Local" configurations which connect to `localhost:9092` and `localhost:8086`.

## 🏗️ Architecture

This project has been refactored from a monolithic `main.go` into a clean, maintainable modular architecture:

```
testing-studio/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── config/                 # Configuration management
│   ├── handlers/               # HTTP request handlers
│   ├── kafka/                  # Kafka operations (pull, publish, Avro)
│   ├── pubsub/                 # PubSub operations
│   ├── templates/              # HTML templates and UI
│   └── types/                  # Shared data structures
├── go.mod
└── README.md
```

For detailed architecture documentation, see [ARCHITECTURE.md](ARCHITECTURE.md).

## 🎯 Why This Tool?

### The Problem
Working with PubSub emulator on macOS with corporate proxies causes HTTP/2 errors:

```bash
# This fails on macOS:
curl localhost:8086/v1/projects/.../pull
# Error: RST_STREAM closed stream. HTTP/2 error code: PROTOCOL_ERROR
```

**Root cause:** Corporate proxy (`localhost:3128`) intercepts localhost connections and adds headers that break the PubSub emulator.

### The Solution
CloudEvents Explorer uses the native Google PubSub Go SDK, completely bypassing HTTP/proxy layers. No more errors! 🎉

## 📖 Usage Guide

### 1. Connection Settings

Select from saved configurations or create new ones:

- **Configuration Name**: Friendly name (e.g., "TMS Local")
- **Emulator Host**: `localhost:8086` (for local emulator)
- **Project ID**: `tms-suncorp-local`
- **Subscription ID**: `cloudevents.subscription`
- **Max Messages**: How many to pull (1-100)

### 2. Pull Messages

Click **▶ Pull Messages** to fetch new messages. They'll be added to the top of the list and persisted in memory.

### 3. View Message Details

Each message shows:
- **Type**: CloudEvent type (e.g., `anzx.migration.tms.v1alpha1.migration.phase.completed`)
- **Subject**: Migration subject (e.g., `migrations/4015648042`)
- **Source**: Event source
- **Schema**: Proto schema reference
- **Timestamp**: When the message was published
- **JSON Payload**: Syntax-highlighted, expandable data

### 4. Expand/Collapse Messages

- Messages start **collapsed** for easy scanning
- **Click any message** to expand and see full details
- **Click again** to collapse

### 5. Clear All

Click **🗑 Clear All** to remove all messages from memory.

## 🎨 UI Features

### Redpanda-Style Design
- **Dark theme** with GitHub color palette
- **Top navigation bar** with tabs
- **Panel-based layout** for organized content
- **Syntax highlighting** for JSON (blue keys, cyan strings, etc.)
- **Smooth animations** for expand/collapse

### Message Cards
- Collapsed by default for quick scanning
- Shows essential info (type, subject, time) at a glance
- Expand to see full metadata + JSON payload
- Color-coded type badges

### Status Notifications
- Success toasts (green)
- Error toasts (red)
- Auto-dismiss after 3 seconds

## 📂 Configuration File

Saved in `configs.json`:

```json
{
  "configs": [
    {
      "name": "TMS Local",
      "emulatorHost": "localhost:8086",
      "projectId": "tms-suncorp-local",
      "subscriptionId": "cloudevents.subscription"
    }
  ]
}
```

## 🔧 Use Cases

### Debugging Workflows
1. Trigger a Temporal workflow (e.g., WriteProfileDiaryNote)
2. Pull messages to see published CloudEvents
3. Expand to view phase started/completed events
4. Inspect JSON payloads for customer data

### Testing Event Publishing
1. Send an event via EventMesh
2. Pull messages to verify it was published
3. Check the event structure and data
4. Confirm schema and type are correct

### Monitoring Migration Events
1. Configure for production PubSub
2. Pull messages to monitor migrations
3. Track phase transitions
4. Debug issues with customer data

## 🏗️ Architecture

### Backend (Go)
- **HTTP server** on port 8888
- **PubSub SDK** for native gRPC connection
- **In-memory storage** for message persistence
- **JSON API** for frontend communication

### Frontend (Vanilla JS)
- **No dependencies** - pure HTML/CSS/JS
- **Syntax highlighter** implemented from scratch
- **Collapsible UI** with expand/collapse state
- **Real-time updates** via fetch API

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | Serve HTML UI |
| `/api/configs` | GET | Get saved configurations |
| `/api/configs` | POST | Save a configuration |
| `/api/pull` | POST | Pull messages from PubSub |
| `/api/messages` | GET | Get stored messages |

## 🎯 Example: Viewing CapDiary Events

After triggering a Unica event → WriteProfileDiaryNote workflow:

1. Open http://localhost:8888
2. Select "TMS Local" config
3. Click "Pull Messages"
4. See 2 events:
   - **MigrationPhaseStarted** (STATUS_IN_PROGRESS)
   - **MigrationPhaseCompleted** (STATUS_COMPLETED)
5. Click to expand and view JSON data
6. Inspect customer ID, group ID, phase, status

## 🐛 Troubleshooting

### Docker: Kafka/PubSub not pulling messages

**Problem**: Running in Docker but Kafka or PubSub pulls return zero messages.

**Solution**: Make sure you're using the **"Docker"** configurations in the UI:
- **Kafka**: Select "Unica Events (Docker)" - connects to `dep_redpanda:9092`
- **PubSub**: Select "TMS PubSub (Docker)" - connects to `dep_pubsub:8086`

The Docker container runs on the `devstack_devstack_network` and uses Docker service names instead of `localhost`.

### Local: Kafka/PubSub not pulling messages

**Problem**: Running locally (`go run`) but can't connect to services.

**Solution**: Use the **"Local"** configurations:
- **Kafka**: Select "Unica Events (Local)" - connects to `localhost:19092`
- **PubSub**: Select "TMS PubSub (Local)" - connects to `localhost:8086`

### Port 8888 already in use

```bash
# Kill any process using port 8888
lsof -ti:8888 | xargs kill -9
```

### Can't connect to devstack services

- Check devstack is running: `docker ps | grep dep_redpanda`
- Ensure Testing Studio container is on devstack network: `docker inspect testing-studio | grep devstack_devstack_network`
- Verify DNS resolution: `docker exec testing-studio getent hosts dep_redpanda`

### Container keeps restarting

Check logs for errors:
```bash
make logs
# or
docker logs testing-studio
```

## 🚀 Advanced Usage

### Multiple Environments
Save configs for different environments:
- **TMS Local** - localhost:8086
- **Dev PubSub** - dev-pubsub-host:8086
- **Staging** - staging-pubsub-host:8086

Switch between them with the dropdown!

### Filtering Messages
Currently shows all messages. To filter:
1. Pull all messages
2. Use browser search (Cmd+F) to find specific types/subjects
3. Expand matching messages

### Exporting Messages
Messages are stored client-side. To export:
1. Open browser console
2. Run: `copy(JSON.stringify(messagesData, null, 2))`
3. Paste into a file

## 📝 Development

### Build
```bash
go build -o cloudevents-explorer main.go
./cloudevents-explorer
```

### Dependencies
- `cloud.google.com/go/pubsub` - PubSub client
- `google.golang.org/api/option` - API options

### Hot Reload
Use `air` or `fresh` for auto-reload during development.

## 🎨 Color Palette

GitHub Dark Theme:
- Background: `#0d1117`
- Panel: `#161b22`
- Border: `#30363d`
- Text: `#c9d1d9`
- Accent: `#58a6ff` (blue)
- Success: `#238636` (green)
- Error: `#da3633` (red)

## 📜 License

Internal tool for TMS Suncorp development.

## 🙏 Acknowledgments

- UI inspired by [Redpanda Console](https://redpanda.com/)
- Color scheme from [GitHub Dark Theme](https://github.com/)
- Built to solve real proxy pain on macOS 🎯

---

**Made with ❤️ for debugging CloudEvents without proxy headaches**