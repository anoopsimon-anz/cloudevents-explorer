# 📡 CloudEvents Explorer

A Redpanda Console-inspired web tool for exploring Google Cloud PubSub CloudEvents with beautiful dark mode UI and message persistence.

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![Go](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)

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

```bash
cd ~/scratches/cloudevents-explorer
go run main.go
```

Open http://localhost:8888 in your browser.

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

### Port 8888 in use?
Change the port in `main.go`:
```go
port := "9999"  // Use different port
```

### Can't connect to PubSub?
- Check devstack is running: `docker ps | grep pubsub`
- Verify emulator host: should be `localhost:8086`
- Check subscription exists in your project

### Messages not appearing?
- Verify subscription ID is correct
- Check that events are being published
- Look at browser console for errors

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